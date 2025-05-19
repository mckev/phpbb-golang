package myforum

import (
	"context"
	"fmt"

	"phpbb-golang/model"
)

func PopulateDb(ctx context.Context) error {
	// Schema
	{
		err := model.DropDb(ctx, "posts")
		if err != nil {
			return fmt.Errorf("Error while dropping posts table: %s", err)
		}
		err = model.DropDb(ctx, "topics")
		if err != nil {
			return fmt.Errorf("Error while dropping topics table: %s", err)
		}
		err = model.DropDb(ctx, "forums")
		if err != nil {
			return fmt.Errorf("Error while dropping forums table: %s", err)
		}
	}
	{
		err := model.InitForums(ctx)
		if err != nil {
			return fmt.Errorf("Error while initializing forums table: %s", err)
		}
		err = model.InitTopics(ctx)
		if err != nil {
			return fmt.Errorf("Error while initializing topics table: %s", err)
		}
		err = model.InitPosts(ctx)
		if err != nil {
			return fmt.Errorf("Error while initializing posts table: %s", err)
		}
	}

	// Forums
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
		forumBBId, err := model.InsertForum(ctx, forumBId, "Now Hear This!", "Announcements from the Management and assistance with forum software. New to FWF? Please consider introducing yourself")
		if err != nil {
			return err
		}
		{
			// Topics: Now Hear This!
			topicBB1Id, err := model.InsertTopic(ctx, forumBBId, "We're now powered by phpBB 3.3")
			if err != nil {
				return err
			}
			_, err = model.InsertPost(ctx, topicBB1Id, forumBBId, "We're now powered by phpBB 3.3", `We are pleased to announce that the board has been upgraded to the [url=https://www.phpbb.com/about/launch/]phpBB 3.3[/url] Feature Release.

There have only been minor changes to the user interface and features that a keen eyed observer might see. We expect that many of you probably won't be able to notice any difference.

For the most part this upgrade was about getting the underlying components and frameworks to a more modern base, which should result in some performance improvement.

For those keeping track, the fix for [url=https://www.financialwisdomforum.org/forum/viewtopic.php?p=650049#p650049]this bug[/url] hasn't been included in this phpBB feature release. It is due in the next, currently unscheduled, bugfix release. A reminder that the workaround is to delete the PM when reading it, not from the list of PMs.

Please use this topic if you encounter any problems.`)
			if err != nil {
				return err
			}
			_, err = model.InsertPost(ctx, topicBB1Id, forumBBId, "Re: We're now powered by phpBB 3.3", `[blockquote user_name="Peculiar_Investor" user_id="636" post_id="659301" time="1586687280"]Has anyone else even noticed we upgraded and have you found anything else that might have changed?[/blockquote]

I would't know anything had changed if not for your posts/updates.
As always, thanks for the work you and others do to keep FWF such an excellent site and resource.`)
			if err != nil {
				return err
			}
			_, err = model.InsertPost(ctx, topicBB1Id, forumBBId, "Re: We're now powered by phpBB 3.3", `Haven't noticed any differences.`)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
