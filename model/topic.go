package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	INVALID_TOPIC_ID    = -1
	MAX_TOPICS_PER_PAGE = 25
)

type Topic struct {
	TopicId     int    `json:"topic_id"`
	ForumId     int    `json:"forum_id"`
	TopicTitle  string `json:"topic_title"`
	TopicUserId int    `json:"topic_user_id"`
	TopicTime   int64  `json:"topic_time"`
	// Derived properties to speed up
	TopicNumPosts          int    `json:"topic_num_posts"`
	TopicNumViews          int    `json:"topic_num_views"`
	TopicFirstPostId       int    `json:"topic_first_post_id"`
	TopicFirstPostUserName string `json:"topic_first_post_user_name"`
	TopicLastPostId        int    `json:"topic_last_post_id"`
	TopicLastPostUserId    int    `json:"topic_last_post_user_id"`
	TopicLastPostUserName  string `json:"topic_last_post_user_name"`
	TopicLastPostTime      int64  `json:"topic_last_post_time"`
}

func InitTopics(ctx context.Context) error {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	// On PostgreSQL:  topic_id INT(10) PRIMARY KEY AUTOINCREMENT NOT NULL,
	sql := `CREATE TABLE topics (
		topic_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		forum_id MEDIUMINT(8) NOT NULL DEFAULT '0',
		topic_title VARCHAR(255) NOT NULL DEFAULT '',
		topic_user_id INT(10) NOT NULL DEFAULT '0',
		topic_time INT(11) NOT NULL DEFAULT '0',
		topic_num_posts MEDIUMINT(8) NOT NULL DEFAULT '0',
		topic_num_views MEDIUMINT(8) NOT NULL DEFAULT '0',
		topic_first_post_id INT(10) NOT NULL DEFAULT '0',
		topic_first_post_user_name VARCHAR(255) NOT NULL DEFAULT '',
		topic_last_post_id INT(10) NOT NULL DEFAULT '0',
		topic_last_post_time INT(11) NOT NULL DEFAULT '0',
		topic_last_post_user_id INT(10) NOT NULL DEFAULT '0',
		topic_last_post_user_name VARCHAR(255) NOT NULL DEFAULT '',
		FOREIGN KEY (forum_id) REFERENCES forums(forum_id),
		FOREIGN KEY (topic_user_id) REFERENCES users(user_id),
		FOREIGN KEY (topic_last_post_user_id) REFERENCES users(user_id)
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating topics table: %s", err)
	}
	return nil
}

func InsertTopic(ctx context.Context, forumId int, topicTitle string, topicUserId int, topicUserName string) (int, error) {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	now := time.Now().UTC()
	topicTime := now.Unix()
	res, err := db.Exec(`INSERT INTO topics (forum_id, topic_title, topic_time, topic_user_id, topic_first_post_user_name, topic_last_post_user_id, topic_last_post_user_name) VALUES ($1, $2, $3, $4, $5, $6, $7)`, forumId, topicTitle, topicTime, topicUserId, topicUserName, topicUserId, topicUserName)
	if err != nil {
		return INVALID_TOPIC_ID, fmt.Errorf("Error while inserting topic title '%s' with forum id %d into topics table: %s", topicTitle, forumId, err)
	}
	topicId, err := res.LastInsertId()
	if err != nil {
		return INVALID_TOPIC_ID, fmt.Errorf("Error while retrieving last insert id for topic title '%s': %s", topicTitle, err)
	}
	return int(topicId), nil
}

func IncreaseNumPostsForTopic(ctx context.Context, topicId int) error {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	_, err := db.Exec(`UPDATE topics SET topic_num_posts = topic_num_posts + 1 WHERE topic_id = $1`, topicId)
	if err != nil {
		return fmt.Errorf("Error while increasing num posts for topic id %d: %s", topicId, err)
	}
	return nil
}

func UpdateFirstPostOfTopic(ctx context.Context, topicId int, topicFirstPostId int) error {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	result, err := db.Exec(`UPDATE topics SET topic_first_post_id = $1 WHERE topic_id = $2`, topicFirstPostId, topicId)
	if err != nil {
		return fmt.Errorf("Error while updating the first post for topic id %d: %s", topicId, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error while retrieving rows affected while updating the first post for topic id %d: %s", topicId, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No rows were updated while updating the first post for topic id %d", topicId)
	}
	return nil
}

func UpdateLastPostOfTopic(ctx context.Context, topicId int, topicLastPostId int, topicLastPostUserId int, topicLastPostUserName string) error {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	now := time.Now().UTC()
	topicLastPostTime := now.Unix()
	result, err := db.Exec(`UPDATE topics SET topic_last_post_id = $1, topic_last_post_user_id = $2, topic_last_post_user_name = $3, topic_last_post_time = $4 WHERE topic_id = $5`, topicLastPostId, topicLastPostUserId, topicLastPostUserName, topicLastPostTime, topicId)
	if err != nil {
		return fmt.Errorf("Error while updating the last post for topic id %d: %s", topicId, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error while retrieving rows affected while updating the last post for topic id %d: %s", topicId, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No rows were updated while updating the last post for topic id %d", topicId)
	}
	return nil
}

func ListTopics(ctx context.Context, forumId int, startItem int) ([]Topic, error) {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	rows, err := db.Query("SELECT topic_id, forum_id, topic_title, topic_user_id, topic_time, topic_num_posts, topic_num_views, topic_first_post_id, topic_first_post_user_name, topic_last_post_id, topic_last_post_user_id, topic_last_post_user_name, topic_last_post_time FROM topics WHERE forum_id = $1 ORDER BY topic_id LIMIT $2 OFFSET $3", forumId, MAX_TOPICS_PER_PAGE, startItem)
	if err != nil {
		return nil, fmt.Errorf("Error while querying topics table for forum id %d: %s", forumId, err)
	}
	defer rows.Close()
	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(&topic.TopicId, &topic.ForumId, &topic.TopicTitle, &topic.TopicUserId, &topic.TopicTime, &topic.TopicNumPosts, &topic.TopicNumViews, &topic.TopicFirstPostId, &topic.TopicFirstPostUserName, &topic.TopicLastPostId, &topic.TopicLastPostUserId, &topic.TopicLastPostUserName, &topic.TopicLastPostTime); err != nil {
			return nil, fmt.Errorf("Error while scanning rows on topics table for forum id %d: %s", forumId, err)
		}
		topics = append(topics, topic)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on rows on topics table for forum id %d: %s", forumId, err)
	}
	return topics, nil
}

func GetTopic(ctx context.Context, topicId int) (Topic, error) {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	row := db.QueryRow("SELECT topic_id, forum_id, topic_title, topic_user_id, topic_time, topic_num_posts, topic_num_views, topic_first_post_id, topic_first_post_user_name, topic_last_post_id, topic_last_post_user_id, topic_last_post_user_name, topic_last_post_time FROM topics WHERE topic_id = $1", topicId)
	var topic Topic
	if err := row.Scan(&topic.TopicId, &topic.ForumId, &topic.TopicTitle, &topic.TopicUserId, &topic.TopicTime, &topic.TopicNumPosts, &topic.TopicNumViews, &topic.TopicFirstPostId, &topic.TopicFirstPostUserName, &topic.TopicLastPostId, &topic.TopicLastPostUserId, &topic.TopicLastPostUserName, &topic.TopicLastPostTime); err != nil {
		if err == sql.ErrNoRows {
			// No result found
			return Topic{}, nil
		}
		return Topic{}, fmt.Errorf("Error while scanning row on topics table for topic id %d: %s", topicId, err)
	}
	return topic, nil
}
