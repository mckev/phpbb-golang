package model

import (
	"context"
	"fmt"
	"time"
)

type Topic struct {
	TopicId    int    `json:"topic_id"`
	ForumId    int    `json:"forum_id"`
	TopicTitle string `json:"topic_title"`
	TopicTime  int    `json:"topic_time"`
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

func InsertTopic(ctx context.Context, forumId int, topicTitle string) (int, error) {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	now := time.Now().UTC()
	topicTime := now.Unix()
	res, err := db.Exec(`INSERT INTO topics (forum_id, topic_title, topic_time) VALUES ($1, $2, $3)`, forumId, topicTitle, topicTime)
	if err != nil {
		return -1, fmt.Errorf("Error while inserting topic '%s' with parent forum %d into topics table: %s", topicTitle, forumId, err)
	}
	topicId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Error while retrieving last insert id for topic '%s': %s", topicTitle, err)
	}
	return int(topicId), nil
}

func ListTopics(ctx context.Context, forumId int) ([]Topic, error) {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	rows, err := db.Query("SELECT topic_id, forum_id, topic_title, topic_time FROM topics WHERE forum_id = $1 ORDER BY topic_id", forumId)
	if err != nil {
		return nil, fmt.Errorf("Error while querying topics table: %s", err)
	}
	defer rows.Close()
	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(&topic.TopicId, &topic.ForumId, &topic.TopicTitle, &topic.TopicTime); err != nil {
			return nil, fmt.Errorf("Error while scanning rows on topics table: %s", err)
		}
		topics = append(topics, topic)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on rows on topics table: %s", err)
	}
	return topics, nil
}
