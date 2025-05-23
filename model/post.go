package model

import (
	"context"
	"fmt"
	"time"
)

const (
	INVALID_POST       = -1
	MAX_POSTS_PER_PAGE = 25
)

type Post struct {
	PostId      int    `json:"post_id"`
	TopicId     int    `json:"topic_id"`
	ForumId     int    `json:"forum_id"`
	PostSubject string `json:"post_subject"`
	PostText    string `json:"post_text"`
	UserId      int    `json:"user_id"`
	PostTime    int64  `json:"post_time"`
}

func InitPosts(ctx context.Context) error {
	db := OpenDb(ctx, "posts")
	defer db.Close()
	// On PostgreSQL:  post_id INT(10) PRIMARY KEY AUTOINCREMENT NOT NULL,
	sql := `CREATE TABLE posts (
		post_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		topic_id INT(10) NOT NULL DEFAULT '0',
		forum_id MEDIUMINT(8) NOT NULL DEFAULT '0',
		post_subject VARCHAR(255) NOT NULL DEFAULT '',
		post_text MEDIUMTEXT NOT NULL DEFAULT '',
		user_id INT(10) NOT NULL DEFAULT '0',
		post_time INT(11) NOT NULL DEFAULT '0',
		FOREIGN KEY (topic_id) REFERENCES topics(topic_id),
		FOREIGN KEY (forum_id) REFERENCES forums(forum_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating posts table: %s", err)
	}
	return nil
}

func InsertPost(ctx context.Context, topicId int, forumId int, postSubject string, postText string, userId int) (int, error) {
	db := OpenDb(ctx, "posts")
	defer db.Close()
	now := time.Now().UTC()
	postTime := now.Unix()
	res, err := db.Exec(`INSERT INTO posts (topic_id, forum_id, post_subject, post_text, post_time, user_id) VALUES ($1, $2, $3, $4, $5, $6)`, topicId, forumId, postSubject, postText, postTime, userId)
	if err != nil {
		return INVALID_POST, fmt.Errorf("Error while inserting post subject '%s' with topic id %d and forum id %d into posts table: %s", postSubject, topicId, forumId, err)
	}
	postId, err := res.LastInsertId()
	if err != nil {
		return INVALID_POST, fmt.Errorf("Error while retrieving last insert id for post subject '%s': %s", postSubject, err)
	}
	return int(postId), nil
}

func ListPosts(ctx context.Context, topicId int) ([]Post, error) {
	db := OpenDb(ctx, "posts")
	defer db.Close()
	rows, err := db.Query("SELECT post_id, topic_id, forum_id, post_subject, post_text, post_time, user_id FROM posts WHERE topic_id = $1 ORDER BY post_id", topicId)
	if err != nil {
		return nil, fmt.Errorf("Error while querying posts table for topic id %d: %s", topicId, err)
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostId, &post.TopicId, &post.ForumId, &post.PostSubject, &post.PostText, &post.PostTime, &post.UserId); err != nil {
			return nil, fmt.Errorf("Error while scanning rows on posts table for topic id %d: %s", topicId, err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on rows on posts table for topic id %d: %s", topicId, err)
	}
	return posts, nil
}
