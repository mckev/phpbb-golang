package myforum

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"phpbb-golang/internal/logger"
)

func OpenDb(ctx context.Context, tableName string) *sql.DB {
	// Refs:
	//   - https://go.dev/doc/tutorial/database-access
	//   - https://www.phpbb.com/demo/
	dbDSN := fmt.Sprintf("file:./db/%s.db?_foreign_keys=on", tableName)
	db, err := sql.Open("sqlite3", dbDSN)
	if err != nil {
		logger.Fatalf(ctx, "Error while opening Database DSN %s: %s", dbDSN, err)
	}
	return db
}

func PopulateDb(ctx context.Context) error {
	// Ref: https://www.erdcloud.com/d/23zvQbme2zHiLtYmf

	// Schema
	var dbForum, dbTopic, dbPost *sql.DB
	{
		dbForum = OpenDb(ctx, "forums")
		sql := `DROP TABLE IF EXISTS forums`
		_, err := dbForum.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error while dropping table forums: %s", err)
		}
		// On PostgreSQL:  forum_id MEDIUMINT(8) PRIMARY KEY AUTOINCREMENT NOT NULL,
		sql = `CREATE TABLE forums (
			forum_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			parent_id MEDIUMINT(8) NOT NULL DEFAULT '0',
			forum_name VARCHAR(255) NOT NULL DEFAULT '',
			forum_desc TEXT NOT NULL DEFAULT '',
			FOREIGN KEY (parent_id) REFERENCES forums(forum_id)
		)`
		_, err = dbForum.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error while creating table forums: %s", err)
		}
		_, err = dbForum.Exec(`INSERT INTO forums (forum_id, parent_id, forum_name, forum_desc) VALUES (?, ?, ?, ?)`, 0, 0, "Root", "")
		if err != nil {
			return fmt.Errorf("Error while inserting table forums: %s", err)
		}
	}
	{
		dbTopic = OpenDb(ctx, "topics")
		sql := `DROP TABLE IF EXISTS topics`
		_, err := dbTopic.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error while dropping table topics: %s", err)
		}
		// On PostgreSQL:  topic_id INT(10) PRIMARY KEY AUTOINCREMENT NOT NULL,
		sql = `CREATE TABLE topics (
			topic_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			forum_id MEDIUMINT(8) NOT NULL DEFAULT '0',
			topic_title VARCHAR(255) NOT NULL DEFAULT '',
			topic_time INT(11) NOT NULL DEFAULT '0',
			topic_first_post_id INT(10) NOT NULL DEFAULT '0',
			topic_last_post_id INT(10) NOT NULL DEFAULT '0',
			topic_views MEDIUMINT(8) NOT NULL DEFAULT '0',
			FOREIGN KEY (forum_id) REFERENCES forums(forum_id)
		)`
		_, err = dbTopic.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error while creating table topics: %s", err)
		}
	}
	{
		dbPost = OpenDb(ctx, "posts")
		sql := `DROP TABLE IF EXISTS posts`
		_, err := dbPost.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error while dropping table posts: %s", err)
		}
		// On PostgreSQL:  post_id INT(10) PRIMARY KEY AUTOINCREMENT NOT NULL,
		sql = `CREATE TABLE posts (
			post_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			topic_id INT(10) NOT NULL DEFAULT '0',
			forum_id MEDIUMINT(8) NOT NULL DEFAULT '0',
			post_subject VARCHAR(255) NOT NULL DEFAULT '',
			post_text MEDIUMTEXT NOT NULL DEFAULT '',
			post_time INT(11) NOT NULL DEFAULT '0',
			FOREIGN KEY (topic_id) REFERENCES topics(topic_id),
			FOREIGN KEY (forum_id) REFERENCES forums(forum_id)
		)`
		_, err = dbPost.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error while creating table posts: %s", err)
		}
	}

	// Data
	{
		res, err := dbForum.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES (?, ?, ?)`, 0, "Your Money", "")
		if err != nil {
			return fmt.Errorf("Error while inserting table forums: %s", err)
		}
		forumId, _ := res.LastInsertId()
		_, err = dbForum.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES (?, ?, ?)`, forumId, "Financial Planning and Building Portfolios", "Asset allocation, risk, diversification and rebalancing. Pros/cons of hiring a financial advisor. Seeking advice on your portfolio?")
		if err != nil {
			return fmt.Errorf("Error while inserting table forums: %s", err)
		}
		{
			// Subforum example: https://forums.linuxmint.com/
			res, err = dbForum.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES (?, ?, ?)`, forumId, "Retirement, Pensions and Peace of Mind", "Preparing for life after work. RRSPs, RRIFs, TFSAs, annuities and meeting future financial and psychological needs.")
			if err != nil {
				return fmt.Errorf("Error while inserting table forums: %s", err)
			}
			subForumId, _ := res.LastInsertId()
			_, err = dbForum.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES (?, ?, ?)`, subForumId, "Jakarta & Bandung", "Retire on Jakarta and Bandung, Indonesia")
			if err != nil {
				return fmt.Errorf("Error while inserting table forums: %s", err)
			}
			_, err = dbForum.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES (?, ?, ?)`, subForumId, "Pattaya", "Retire on Pattaya, Thailand")
			if err != nil {
				return fmt.Errorf("Error while inserting table forums: %s", err)
			}
			_, err = dbForum.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES (?, ?, ?)`, subForumId, "Kuala Lumpur", "Retire on Kuala Lumpur, Malaysia")
			if err != nil {
				return fmt.Errorf("Error while inserting table forums: %s", err)
			}
		}
		_, err = dbForum.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES (?, ?, ?)`, forumId, "Financial News, Policy and Economics", "Recommended reading, economic debates, predictions and opinions.")
		if err != nil {
			return fmt.Errorf("Error while inserting table forums: %s", err)
		}
	}
	{
		res, err := dbForum.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES (?, ?, ?)`, 0, "Your Life", "")
		if err != nil {
			return fmt.Errorf("Error while inserting table forums: %s", err)
		}
		forumId, _ := res.LastInsertId()
		_, err = dbForum.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES (?, ?, ?)`, forumId, "Community Centre", "Non financial topics: autos; computers; entertainment; gatherings; hobbies; sports and travel.")
		if err != nil {
			return fmt.Errorf("Error while inserting table forums: %s", err)
		}
		_, err = dbForum.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES (?, ?, ?)`, forumId, "Now Hear This!", "Announcements from the Management and assistance with forum software. New to FWF? Please consider introducing yourself")
		if err != nil {
			return fmt.Errorf("Error while inserting table forums: %s", err)
		}
	}
	return nil
}
