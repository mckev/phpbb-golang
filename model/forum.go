package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	ROOT_FORUM_ID              = 0
	INVALID_FORUM_ID           = -1
	MAX_FORUM_NAV_TRAILS_DEPTH = 7
	MAX_FORUM_NODES_DEPTH      = 3
)

type Forum struct {
	ForumId     int    `json:"forum_id"`
	ParentId    int    `json:"parent_id"`
	ForumName   string `json:"forum_name"`
	ForumDesc   string `json:"forum_desc"`
	ForumUserId string `json:"forum_user_id"`
	ForumTime   int64  `json:"forum_time"`
	// Derived properties to speed up
	ForumNumTopics        int    `json:"forum_num_topics"`
	ForumNumPosts         int    `json:"forum_num_posts"`
	ForumLastPostId       int    `json:"forum_last_post_id"`
	ForumLastPostSubject  string `json:"forum_last_post_subject"`
	ForumLastPostTime     int64  `json:"forum_last_post_time"`
	ForumLastPostUserId   int    `json:"forum_last_post_user_id"`
	ForumLastPostUserName string `json:"forum_last_post_user_name"`
}

func InitForums(ctx context.Context) error {
	// Ref: https://www.erdcloud.com/d/23zvQbme2zHiLtYmf
	db := OpenDb(ctx, "forums")
	defer db.Close()
	// On PostgreSQL:  forum_id MEDIUMINT(8) PRIMARY KEY AUTOINCREMENT NOT NULL,
	sql := `CREATE TABLE forums (
		forum_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		parent_id MEDIUMINT(8) NOT NULL DEFAULT '0',
		forum_name VARCHAR(255) NOT NULL DEFAULT '',
		forum_desc TEXT NOT NULL DEFAULT '',
		forum_user_id INT(10) NOT NULL DEFAULT '0',
		forum_time INT(11) NOT NULL DEFAULT '0',
		forum_num_topics MEDIUMINT(8) NOT NULL DEFAULT '0',
		forum_num_posts MEDIUMINT(8) NOT NULL DEFAULT '0',
		forum_last_post_id INT(10) NOT NULL DEFAULT '0',
		forum_last_post_subject VARCHAR(255) NOT NULL DEFAULT '',
		forum_last_post_user_id INT(10) NOT NULL DEFAULT '0',
		forum_last_post_user_name VARCHAR(255) NOT NULL DEFAULT '',
		forum_last_post_time INT(11) NOT NULL DEFAULT '0',
		FOREIGN KEY (parent_id) REFERENCES forums(forum_id),
		FOREIGN KEY (forum_user_id) REFERENCES users(user_id),
		FOREIGN KEY (forum_last_post_user_id) REFERENCES users(user_id)
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating forums table: %s", err)
	}
	_, err = db.Exec("INSERT INTO forums (forum_id, parent_id, forum_name, forum_desc, forum_user_id, forum_last_post_user_id) VALUES ($1, $2, $3, $4, $5, $6)", ROOT_FORUM_ID, ROOT_FORUM_ID, "Root Forum", "", ADMIN_USER_ID, ADMIN_USER_ID)
	if err != nil {
		return fmt.Errorf("Error while inserting Root Forum into forums table: %s", err)
	}
	return nil
}

func InsertForum(ctx context.Context, parentId int, forumName string, forumDesc string, forumUserId int) (int, error) {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	now := time.Now().UTC()
	forumTime := now.Unix()
	res, err := db.Exec("INSERT INTO forums (parent_id, forum_name, forum_desc, forum_user_id, forum_time, forum_last_post_user_id) VALUES ($1, $2, $3, $4, $5, $6)", parentId, forumName, forumDesc, forumUserId, forumTime, forumUserId)
	if err != nil {
		return INVALID_FORUM_ID, fmt.Errorf("Error while inserting forum name '%s' with forum description '%s' and parent forum %d into forums table: %s", forumName, forumDesc, parentId, err)
	}
	forumId, err := res.LastInsertId()
	if err != nil {
		return INVALID_FORUM_ID, fmt.Errorf("Error while retrieving last insert id for forum name '%s': %s", forumName, err)
	}
	return int(forumId), nil
}

func IncreaseNumTopicsForForum(ctx context.Context, forumId int) error {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	_, err := db.Exec("UPDATE forums SET forum_num_topics = forum_num_topics + 1 WHERE forum_id = $1", forumId)
	if err != nil {
		return fmt.Errorf("Error while increasing num topics for forum id %d: %s", forumId, err)
	}
	return nil
}

func IncreaseNumPostsForForum(ctx context.Context, forumId int) error {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	_, err := db.Exec("UPDATE forums SET forum_num_posts = forum_num_posts + 1 WHERE forum_id = $1", forumId)
	if err != nil {
		return fmt.Errorf("Error while increasing num posts for forum id %d: %s", forumId, err)
	}
	return nil
}

func UpdateLastPostOfForum(ctx context.Context, forumId int, forumLastPostId int, forumLastPostSubject string, forumLastPostUserId int, forumLastPostUserName string) error {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	now := time.Now().UTC()
	topicLastPostTime := now.Unix()
	result, err := db.Exec("UPDATE forums SET forum_last_post_id = $1, forum_last_post_subject = $2, forum_last_post_user_id = $3, forum_last_post_user_name = $4, forum_last_post_time = $5 WHERE forum_id = $6", forumLastPostId, forumLastPostSubject, forumLastPostUserId, forumLastPostUserName, topicLastPostTime, forumId)
	if err != nil {
		return fmt.Errorf("Error while updating the last post for forum id %d: %s", forumId, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error while retrieving rows affected while updating the last post for forum id %d: %s", forumId, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No rows were updated while updating the last post for forum id %d", forumId)
	}
	return nil
}

func ListForums(ctx context.Context) ([]Forum, error) {
	// Ref: https://go.dev/doc/tutorial/database-access
	db := OpenDb(ctx, "forums")
	defer db.Close()
	rows, err := db.Query("SELECT forum_id, parent_id, forum_name, forum_desc, forum_user_id, forum_time, forum_num_topics, forum_num_posts, forum_last_post_id, forum_last_post_subject, forum_last_post_user_id, forum_last_post_user_name, forum_last_post_time FROM forums ORDER BY forum_id")
	if err != nil {
		return nil, fmt.Errorf("Error while querying forums table: %s", err)
	}
	defer rows.Close()
	var forums []Forum
	for rows.Next() {
		var forum Forum
		if err := rows.Scan(&forum.ForumId, &forum.ParentId, &forum.ForumName, &forum.ForumDesc, &forum.ForumUserId, &forum.ForumTime, &forum.ForumNumTopics, &forum.ForumNumPosts, &forum.ForumLastPostId, &forum.ForumLastPostSubject, &forum.ForumLastPostUserId, &forum.ForumLastPostUserName, &forum.ForumLastPostTime); err != nil {
			return nil, fmt.Errorf("Error while scanning rows on forums table: %s", err)
		}
		forums = append(forums, forum)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on rows on forums table: %s", err)
	}
	return forums, nil
}

func GetForum(ctx context.Context, forumId int) (Forum, error) {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	row := db.QueryRow("SELECT forum_id, parent_id, forum_name, forum_desc, forum_user_id, forum_time, forum_num_topics, forum_num_posts, forum_last_post_id, forum_last_post_subject, forum_last_post_user_id, forum_last_post_user_name, forum_last_post_time FROM forums WHERE forum_id = $1", forumId)
	var forum Forum
	if err := row.Scan(&forum.ForumId, &forum.ParentId, &forum.ForumName, &forum.ForumDesc, &forum.ForumUserId, &forum.ForumTime, &forum.ForumNumTopics, &forum.ForumNumPosts, &forum.ForumLastPostId, &forum.ForumLastPostSubject, &forum.ForumLastPostUserId, &forum.ForumLastPostUserName, &forum.ForumLastPostTime); err != nil {
		if err == sql.ErrNoRows {
			// No result found
			return Forum{}, nil
		}
		return Forum{}, fmt.Errorf("Error while scanning row on forums table for forum id %d: %s", forumId, err)
	}
	return forum, nil
}
