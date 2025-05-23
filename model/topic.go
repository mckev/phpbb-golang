package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	INVALID_TOPIC       = -1
	MAX_TOPICS_PER_PAGE = 25
)

type Topic struct {
	TopicId       int    `json:"topic_id"`
	ForumId       int    `json:"forum_id"`
	TopicTitle    string `json:"topic_title"`
	TopicTime     int    `json:"topic_time"`
	TopicNumPosts int    `json:"topic_num_posts"`
	TopicNumViews int    `json:"topic_num_views"`
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
		topic_num_posts MEDIUMINT(8) NOT NULL DEFAULT '0',
		topic_num_views MEDIUMINT(8) NOT NULL DEFAULT '0',
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
		return INVALID_TOPIC, fmt.Errorf("Error while inserting topic title '%s' with forum id %d into topics table: %s", topicTitle, forumId, err)
	}
	topicId, err := res.LastInsertId()
	if err != nil {
		return INVALID_TOPIC, fmt.Errorf("Error while retrieving last insert id for topic title '%s': %s", topicTitle, err)
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

func ListTopics(ctx context.Context, forumId int) ([]Topic, error) {
	db := OpenDb(ctx, "topics")
	defer db.Close()
	rows, err := db.Query("SELECT topic_id, forum_id, topic_title, topic_time, topic_num_posts, topic_num_views FROM topics WHERE forum_id = $1 ORDER BY topic_id", forumId)
	if err != nil {
		return nil, fmt.Errorf("Error while querying topics table for forum id %d: %s", forumId, err)
	}
	defer rows.Close()
	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(&topic.TopicId, &topic.ForumId, &topic.TopicTitle, &topic.TopicTime, &topic.TopicNumPosts, &topic.TopicNumViews); err != nil {
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
	row := db.QueryRow("SELECT topic_id, forum_id, topic_title, topic_time, topic_num_posts, topic_num_views FROM topics WHERE topic_id = $1", topicId)
	var topic Topic
	if err := row.Scan(&topic.TopicId, &topic.ForumId, &topic.TopicTitle, &topic.TopicTime, &topic.TopicNumPosts, &topic.TopicNumViews); err != nil {
		if err == sql.ErrNoRows {
			// No result found
			return Topic{}, nil
		}
		return Topic{}, fmt.Errorf("Error while scanning row on topics table for topic id %d: %s", topicId, err)
	}
	return topic, nil
}
