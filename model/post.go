package model

import (
	"context"
	"fmt"
)

type Post struct {
	PostId      int    `json:"post_id"`
	TopicId     int    `json:"topic_id"`
	ForumId     int    `json:"forum_id"`
	PostSubject string `json:"post_subject"`
	PostText    string `json:"post_text"`
	PostTime    int    `json:"post_time"`
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
