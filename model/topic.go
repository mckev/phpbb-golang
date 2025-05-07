package model

import (
	"context"
	"fmt"
)

type Topic struct {
	TopicId    int    `json:"topic_id"`
	ForumId    int    `json:"forum_id"`
	TopicTitle string `json:"topic_title"`
	TopicTime  int    `json:"topic_time"`
	ForumName  string `json:"forum_name"`
	ForumDesc  string `json:"forum_desc"`
}

func InitTopics(ctx context.Context) error {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	// On PostgreSQL:  topic_id INT(10) PRIMARY KEY AUTOINCREMENT NOT NULL,
	sql := `CREATE TABLE topics (
		topic_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		forum_id MEDIUMINT(8) NOT NULL DEFAULT '0',
		topic_title VARCHAR(255) NOT NULL DEFAULT '',
		topic_time INT(11) NOT NULL DEFAULT '0',
		topic_first_post_id INT(10) NOT NULL DEFAULT '0',
		topic_last_post_id INT(10) NOT NULL DEFAULT '0',
		topic_views MEDIUMINT(8) NOT NULL DEFAULT '0',
		FOREIGN KEY (forum_id) REFERENCES forums(forum_id)
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating topics table: %s", err)
	}
	return nil
}
