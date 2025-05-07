package myforum

import (
	"context"
	"database/sql"
	"fmt"

	"phpbb-golang/model"
)

func PopulateDb(ctx context.Context) error {
	// Ref: https://www.erdcloud.com/d/23zvQbme2zHiLtYmf

	// Schema
	var dbForum, dbTopic, dbPost *sql.DB
	{
		dbForum = model.OpenDb(ctx, "forums")
		defer dbForum.Close()
		sql := `DROP TABLE IF EXISTS forums`
		_, err := dbForum.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error while dropping forums table: %s", err)
		}
		err = model.InitForums(ctx)
		if err != nil {
			return fmt.Errorf("Error while initializing forums table: %s", err)
		}
	}
	{
		dbTopic = model.OpenDb(ctx, "topics")
		defer dbTopic.Close()
		sql := `DROP TABLE IF EXISTS topics`
		_, err := dbTopic.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error while dropping topics table: %s", err)
		}
		err = model.InitTopics(ctx)
		if err != nil {
			return fmt.Errorf("Error while initializing topics table: %s", err)
		}
	}
	{
		dbPost = model.OpenDb(ctx, "posts")
		defer dbPost.Close()
		sql := `DROP TABLE IF EXISTS posts`
		_, err := dbPost.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error while dropping posts table: %s", err)
		}
		err = model.InitPosts(ctx)
		if err != nil {
			return fmt.Errorf("Error while initializing posts table: %s", err)
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
