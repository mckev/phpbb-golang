package myforum

import (
	"context"
	"fmt"

	"phpbb-golang/model"
)

func PopulateDb(ctx context.Context) error {
	// Schema
	{
		err := model.DropDb(ctx, "forums")
		if err != nil {
			return fmt.Errorf("Error while dropping forums table: %s", err)
		}
		err = model.InitForums(ctx)
		if err != nil {
			return fmt.Errorf("Error while initializing forums table: %s", err)
		}
	}
	{
		err := model.DropDb(ctx, "topics")
		if err != nil {
			return fmt.Errorf("Error while dropping topics table: %s", err)
		}
		err = model.InitTopics(ctx)
		if err != nil {
			return fmt.Errorf("Error while initializing topics table: %s", err)
		}
	}
	{
		err := model.DropDb(ctx, "posts")
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
		forumAId, err := model.InsertForum(ctx, 0, "Your Money", "")
		if err != nil {
			return err
		}
		_, err = model.InsertForum(ctx, forumAId, "Financial Planning and Building Portfolios", "Asset allocation, risk, diversification and rebalancing. Pros/cons of hiring a financial advisor. Seeking advice on your portfolio?")
		if err != nil {
			return err
		}
		forumABId, err := model.InsertForum(ctx, forumAId, "Retirement, <script>alert('Test XSS')</script> Pensions and Peace of Mind", "Preparing for life after work. RRSPs, RRIFs, TFSAs, annuities and meeting future financial and psychological needs.")
		if err != nil {
			return err
		}
		{
			// Subforum example: https://forums.linuxmint.com/
			_, err := model.InsertForum(ctx, forumABId, "Jakarta & Bandung", "Retire on Jakarta and Bandung, Indonesia")
			if err != nil {
				return err
			}
			_, err = model.InsertForum(ctx, forumABId, "Pattaya", "Retire on Pattaya, Thailand")
			if err != nil {
				return err
			}
			_, err = model.InsertForum(ctx, forumABId, "Kuala Lumpur", "Retire on Kuala Lumpur, Malaysia")
			if err != nil {
				return err
			}
		}
		_, err = model.InsertForum(ctx, forumAId, "Financial News, Policy and Economics", "Recommended reading, economic debates, predictions and opinions.")
		if err != nil {
			return err
		}
	}
	{
		forumBId, err := model.InsertForum(ctx, 0, "Your Life", "")
		if err != nil {
			return err
		}
		_, err = model.InsertForum(ctx, forumBId, "Community Centre", "Non financial topics: autos; computers; entertainment; gatherings; hobbies; sports and travel.")
		if err != nil {
			return err
		}
		_, err = model.InsertForum(ctx, forumBId, "Now Hear This!", "Announcements from the Management and assistance with forum software. New to FWF? Please consider introducing yourself")
		if err != nil {
			return err
		}
	}
	return nil
}
