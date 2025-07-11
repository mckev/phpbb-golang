package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	INVALID_POST_ID    = -1
	MAX_POSTS_PER_PAGE = 25
)

type Post struct {
	PostId       int    `json:"post_id"`
	TopicId      int    `json:"topic_id"`
	ForumId      int    `json:"forum_id"`
	PostSubject  string `json:"post_subject"`
	PostText     string `json:"post_text"`
	PostUserId   int    `json:"post_user_id"`
	PostUserName string `json:"post_user_name"`
	PostTime     int64  `json:"post_time"`
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
		post_user_id INT(10) NOT NULL DEFAULT '0',
		post_user_name VARCHAR(255) NOT NULL DEFAULT '',
		post_time INT(11) NOT NULL DEFAULT '0',
		FOREIGN KEY (topic_id) REFERENCES topics(topic_id),
		FOREIGN KEY (forum_id) REFERENCES forums(forum_id),
		FOREIGN KEY (post_user_id) REFERENCES users(user_id)
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating posts table: %s", err)
	}
	return nil
}

func InsertPost(ctx context.Context, topicId int, forumId int, postSubject string, postText string, postUserId int, postUserName string) (int, error) {
	db := OpenDb(ctx, "posts")
	defer db.Close()
	now := time.Now().UTC()
	postTime := now.Unix()
	result, err := db.Exec("INSERT INTO posts (topic_id, forum_id, post_subject, post_text, post_user_id, post_user_name, post_time) VALUES ($1, $2, $3, $4, $5, $6, $7)", topicId, forumId, postSubject, postText, postUserId, postUserName, postTime)
	if err != nil {
		return INVALID_POST_ID, fmt.Errorf("Error while inserting post subject '%s' with topic id %d and forum id %d into posts table: %s", postSubject, topicId, forumId, err)
	}
	postId, err := result.LastInsertId()
	if err != nil {
		return INVALID_POST_ID, fmt.Errorf("Error while retrieving last insert id for post subject '%s': %s", postSubject, err)
	}
	return int(postId), nil
}

func ListPosts(ctx context.Context, topicId int, startItem int) ([]Post, error) {
	db := OpenDb(ctx, "posts")
	defer db.Close()
	rows, err := db.Query("SELECT post_id, topic_id, forum_id, post_subject, post_text, post_user_id, post_user_name, post_time FROM posts WHERE topic_id = $1 ORDER BY post_id LIMIT $2 OFFSET $3", topicId, MAX_POSTS_PER_PAGE, startItem)
	if err != nil {
		return nil, fmt.Errorf("Error while querying posts table for topic id %d: %s", topicId, err)
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostId, &post.TopicId, &post.ForumId, &post.PostSubject, &post.PostText, &post.PostUserId, &post.PostUserName, &post.PostTime); err != nil {
			return nil, fmt.Errorf("Error while scanning rows on posts table for topic id %d: %s", topicId, err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on rows on posts table for topic id %d: %s", topicId, err)
	}
	return posts, nil
}

func GetPost(ctx context.Context, postId int) (Post, error) {
	db := OpenDb(ctx, "posts")
	defer db.Close()
	var post Post
	err := db.
		QueryRow("SELECT post_id, topic_id, forum_id, post_subject, post_text, post_user_id, post_user_name, post_time FROM posts WHERE post_id = $1", postId).
		Scan(&post.PostId, &post.TopicId, &post.ForumId, &post.PostSubject, &post.PostText, &post.PostUserId, &post.PostUserName, &post.PostTime)
	if err != nil {
		if err == sql.ErrNoRows {
			// No result found
			return Post{}, nil
		}
		return Post{}, fmt.Errorf("Error while scanning row on posts table for post id %d: %s", postId, err)
	}
	return post, nil
}

func CountPostCurItem(ctx context.Context, topicId int, postId int) (int, error) {
	// Given a post id, return how many items are there before it
	// 0 <= curItem < totalItems
	db := OpenDb(ctx, "posts")
	defer db.Close()
	var curItem int
	err := db.
		QueryRow("SELECT COUNT(*) AS cur_item FROM posts WHERE topic_id = $1 AND post_id < $2", topicId, postId).
		Scan(&curItem)
	if err != nil {
		return 0, fmt.Errorf("Error while counting current item on topic id %d for post id %d: %s", topicId, postId, err)
	}
	return curItem, nil
}
