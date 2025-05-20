package model

import (
	"context"
	"fmt"
	"time"
)

type Post struct {
	PostId      int    `json:"post_id"`
	TopicId     int    `json:"topic_id"`
	ForumId     int    `json:"forum_id"`
	PostSubject string `json:"post_subject"`
	PostText    string `json:"post_text"`
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
		post_time INT(11) NOT NULL DEFAULT '0',
		FOREIGN KEY (topic_id) REFERENCES topics(topic_id),
		FOREIGN KEY (forum_id) REFERENCES forums(forum_id)
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating posts table: %s", err)
	}
	return nil
}

func InsertPost(ctx context.Context, topicId int, forumId int, postSubject string, postText string) (int, error) {
	db := OpenDb(ctx, "posts")
	defer db.Close()
	now := time.Now().UTC()
	postTime := now.Unix()
	res, err := db.Exec(`INSERT INTO posts (topic_id, forum_id, post_subject, post_text, post_time) VALUES ($1, $2, $3, $4, $5)`, topicId, forumId, postSubject, postText, postTime)
	if err != nil {
		return -1, fmt.Errorf("Error while inserting post '%s' with parent topic %d and parent forum %d into posts table: %s", postSubject, topicId, forumId, err)
	}
	postId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Error while retrieving last insert id for post '%s': %s", postSubject, err)
	}
	return int(postId), nil
}

func ListPosts(ctx context.Context, topicId int) ([]Post, error) {
	db := OpenDb(ctx, "posts")
	defer db.Close()
	rows, err := db.Query("SELECT post_id, topic_id, forum_id, post_subject, post_text, post_time FROM posts WHERE topic_id = $1 ORDER BY post_id", topicId)
	if err != nil {
		return nil, fmt.Errorf("Error while querying posts table: %s", err)
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostId, &post.TopicId, &post.ForumId, &post.PostSubject, &post.PostText, &post.PostTime); err != nil {
			return nil, fmt.Errorf("Error while scanning rows on posts table: %s", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on rows on posts table: %s", err)
	}
	return posts, nil
}
