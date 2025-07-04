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
		err = model.DropDb(ctx, "sessions")
		if err != nil {
			return fmt.Errorf("Error while dropping sessions table: %s", err)
		}
		err = model.DropDb(ctx, "users")
		if err != nil {
			return fmt.Errorf("Error while dropping users table: %s", err)
		}
	}
	{
		err := model.InitUsers(ctx)
		if err != nil {
			return fmt.Errorf("Error while initializing users table: %s", err)
		}
		err = model.InitSessions(ctx)
		if err != nil {
			return fmt.Errorf("Error while initializing sessions table: %s", err)
		}
		err = model.InitForums(ctx)
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

	// Users
	user2Name := "The Management"
	user2Id, err := model.InsertUser(ctx, user2Name, "Password1", "user2@example.com", "Sic transit gloria mundi. Tuesday is usually worse. - Robert A. Heinlein, Starman Jones")
	if err != nil {
		return err
	}
	err = model.SetUserType(ctx, user2Id, model.USER_TYPE_FOUNDER)
	if err != nil {
		return err
	}
	user3Name := "PeculiarInvestor"
	user3Id, err := model.InsertUser(ctx, user3Name, "Password1", "user3@example.com", `[url=https://www.finiki.org/wiki/Main_Page][img]https://www.financialwisdomforum.org/forum/images/icons/icon_wiki.svg[/img]finiki, the Canadian financial wiki[/url] New editors wanted and welcomed, please help collaborate and improve the wiki.

"Normal people... believe that if it ain't broke, don't fix it. Engineers believe that if it ain't broke, it doesn't have enough features yet." - Scott Adams`)
	if err != nil {
		return err
	}
	err = model.SetUserType(ctx, user3Id, model.USER_TYPE_FOUNDER)
	if err != nil {
		return err
	}
	user4Name := "OnlyMyOpinion"
	user4Id, err := model.InsertUser(ctx, user4Name, "Password1", "", "")
	if err != nil {
		return err
	}
	user5Name := "DenisD"
	user5Id, err := model.InsertUser(ctx, user5Name, "Password4", "user5@example.com", "")
	if err != nil {
		return err
	}
	userXssName := `"]an escape[/blockquote]<script>alert('Test XSS User name')</script>`
	userXssId, err := model.InsertUser(ctx, userXssName, "Password1", "user_xss@example.com", "<script>alert('Test XSS User Signature')</script>")
	if err != nil {
		return err
	}

	// Forums
	{
		forumAId, err := model.InsertForum(ctx, model.ROOT_FORUM_ID, "Your Money", "", model.ADMIN_USER_ID)
		if err != nil {
			return err
		}
		_, err = model.InsertForum(ctx, forumAId, "Financial Planning and Building Portfolios", "Asset allocation, risk, diversification and rebalancing. Pros/cons of hiring a financial advisor. Seeking advice on your portfolio?", user2Id)
		if err != nil {
			return err
		}
		forumABId, err := model.InsertForum(ctx, forumAId, "Retirement, <script>alert('Test XSS Forum name')</script> Pensions and Peace of Mind", "Preparing for life after work. RRSPs, RRIFs, TFSAs, annuities and meeting future financial and psychological needs.", user2Id)
		if err != nil {
			return err
		}
		{
			// Subforum example: https://forums.linuxmint.com/
			_, err := model.InsertForum(ctx, forumABId, "Jakarta & Bandung", "Retire on Jakarta and Bandung, Indonesia", user2Id)
			if err != nil {
				return err
			}
			_, err = model.InsertForum(ctx, forumABId, "Pattaya", "Retire on Pattaya, Thailand", user2Id)
			if err != nil {
				return err
			}
			_, err = model.InsertForum(ctx, forumABId, "Kuala Lumpur", "Retire on Kuala Lumpur, Malaysia", user2Id)
			if err != nil {
				return err
			}
		}
		forumACId, err := model.InsertForum(ctx, forumAId, "Financial News, Policy and Economics", "Recommended reading, economic debates, predictions and opinions.", model.ADMIN_USER_ID)
		if err != nil {
			return err
		}
		_, err = model.InsertTopic(ctx, forumACId, "A finanical news", user2Id, user2Name)
		if err != nil {
			return err
		}
		err = model.IncreaseNumTopicsForForum(ctx, forumACId)
		if err != nil {
			return err
		}
	}
	{
		forumBId, err := model.InsertForum(ctx, model.ROOT_FORUM_ID, "Your Life", "", model.ADMIN_USER_ID)
		if err != nil {
			return err
		}
		_, err = model.InsertForum(ctx, forumBId, "Community Centre", "Non financial topics: autos; computers; entertainment; gatherings; hobbies; sports and travel.", user2Id)
		if err != nil {
			return err
		}
		forumBBId, err := model.InsertForum(ctx, forumBId, "Now Hear This!", "Announcements from the Management and assistance with forum software. New to FWF? Please consider introducing yourself", user2Id)
		if err != nil {
			return err
		}
		{
			// Topic of "Now Hear This!" forum : Introduce Yourself
			topicBB1Id, err := model.InsertTopic(ctx, forumBBId, "Introduce Yourself", user2Id, user2Name)
			if err != nil {
				return err
			}
			err = model.IncreaseNumTopicsForForum(ctx, forumBBId)
			if err != nil {
				return err
			}
			postBB1AId, err := model.InsertPost(ctx, topicBB1Id, forumBBId, "Introduce Yourself", `Greetings. The purpose of this thread is to allow new posters to introduce themselves if they wish, giving as much - or as little - background as they want.

Posting on this thread is entirely voluntary - but, if you do wish to post, thank you and welcome to the Financial Wisdom Forum (FWF)!

-- The Management`, user2Id)
			if err != nil {
				return err
			}
			err = model.UpdateFirstPostOfTopic(ctx, topicBB1Id, postBB1AId)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfTopic(ctx, topicBB1Id, postBB1AId, user2Id, user2Name)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfForum(ctx, forumBBId, postBB1AId, "Introduce Yourself", user2Id, user2Name)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForUser(ctx, user2Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForTopic(ctx, topicBB1Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForForum(ctx, forumBBId)
			if err != nil {
				return err
			}
			postBB1BId, err := model.InsertPost(ctx, topicBB1Id, forumBBId, "Re: Introduce <script>alert('Test XSS Post Subject')</script> Yourself", `Hello, <script>alert('Test XSS Post Text')</script> there!
[blockquote=<script>alert('Test BB Attack')</script> user_name="User<script>alert('Test BB Attack')</script>" <script> user_id="123<script>alert('Test BB Attack')</script>" post_id="<script>alert('Test BB Attack')</script>456" time="<script>alert('Test BB Attack')</script>" <script>="<script>"]a <script>alert('Test BB Attack')</script> test[/blockquote]`, userXssId)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfTopic(ctx, topicBB1Id, postBB1BId, userXssId, userXssName)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfForum(ctx, forumBBId, postBB1BId, "Re: Introduce Yourself", userXssId, userXssName)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForUser(ctx, userXssId)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForTopic(ctx, topicBB1Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForForum(ctx, forumBBId)
			if err != nil {
				return err
			}
			postBB1CId, err := model.InsertPost(ctx, topicBB1Id, forumBBId, "Re: Introduce Yourself", `[blockquote user_name=spicy86 post_id=782359 time=1735650047 user_id=17457]
I'm just wondering when this forum was started.
[/blockquote]
Read all about it in [url=https://www.financialwisdomforum.org/history-of-fwf/]History of FWF - Financial Wisdom Forum[/url]
[blockquote]February 18, 2005, FWF goes live using phpBB v2.0.11[/blockquote]
so closing in on 20 years of providing a place "Where Canadian Investors Meet for Financial Education and Empowerment". Most important is FWF is an independent, non-commercial site that is solely run by volunteers.`, user3Id)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfTopic(ctx, topicBB1Id, postBB1CId, user3Id, user3Name)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfForum(ctx, forumBBId, postBB1CId, "Re: Introduce Yourself", user3Id, user3Name)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForUser(ctx, user3Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForTopic(ctx, topicBB1Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForForum(ctx, forumBBId)
			if err != nil {
				return err
			}

			// Topic of "Now Hear This!" forum : We're now powered by phpBB 3.3
			topicBB2Id, err := model.InsertTopic(ctx, forumBBId, "We're now powered by phpBB 3.3", user3Id, user3Name)
			if err != nil {
				return err
			}
			err = model.IncreaseNumTopicsForForum(ctx, forumBBId)
			if err != nil {
				return err
			}
			postBB2AId, err := model.InsertPost(ctx, topicBB2Id, forumBBId, "We're now powered by phpBB 3.3", `We are pleased to announce that the board has been upgraded to the [url=https://www.phpbb.com/about/launch/]phpBB 3.3[/url] Feature Release.

There have only been minor changes to the user interface and features that a keen eyed observer might see. We expect that many of you probably won't be able to notice any difference.

For the most part this upgrade was about getting the underlying components and frameworks to a more modern base, which should result in some performance improvement.

For those keeping track, the fix for [url=https://www.financialwisdomforum.org/forum/viewtopic.php?p=650049#p650049]this bug[/url] hasn't been included in this phpBB feature release. It is due in the next, currently unscheduled, bugfix release. A reminder that the workaround is to delete the PM when reading it, not from the list of PMs.

Please use this topic if you encounter any problems.`, user3Id)
			if err != nil {
				return err
			}
			err = model.UpdateFirstPostOfTopic(ctx, topicBB2Id, postBB2AId)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfTopic(ctx, topicBB2Id, postBB2AId, user3Id, user3Name)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfForum(ctx, forumBBId, postBB2AId, "We're now powered by phpBB 3.3", user3Id, user3Name)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForUser(ctx, user3Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForUser(ctx, user3Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForTopic(ctx, topicBB2Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForForum(ctx, forumBBId)
			if err != nil {
				return err
			}
			postBB2BId, err := model.InsertPost(ctx, topicBB2Id, forumBBId, "Re: We're now powered by phpBB 3.3", `[blockquote user_name="Peculiar_Investor" user_id="636" post_id="2" time="1586687280"]Has anyone else even noticed we upgraded and have you found anything else that might have changed?[/blockquote]
I would't know anything had changed if not for your posts/updates.
As always, thanks for the work you and others do to keep FWF such an excellent site and resource.`, user4Id)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfTopic(ctx, topicBB2Id, postBB2BId, user4Id, user4Name)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfForum(ctx, forumBBId, postBB2BId, "Re: We're now powered by phpBB 3.3", user4Id, user4Name)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForUser(ctx, user4Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForTopic(ctx, topicBB2Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForForum(ctx, forumBBId)
			if err != nil {
				return err
			}
			postBB2CId, err := model.InsertPost(ctx, topicBB2Id, forumBBId, "Re: We're now powered by phpBB 3.3", "Haven't noticed any differences.", user5Id)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfTopic(ctx, topicBB2Id, postBB2CId, user5Id, user5Name)
			if err != nil {
				return err
			}
			err = model.UpdateLastPostOfForum(ctx, forumBBId, postBB2CId, "Re: We're now powered by phpBB 3.3", user5Id, user5Name)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForUser(ctx, user5Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForTopic(ctx, topicBB2Id)
			if err != nil {
				return err
			}
			err = model.IncreaseNumPostsForForum(ctx, forumBBId)
			if err != nil {
				return err
			}

			// Spam posts of "Now Hear This!" forum : We're now powered by phpBB 3.3
			// for i := 4; i <= 250; i++ {
			// 	postBB2DId, err := model.InsertPost(ctx, topicBB2Id, forumBBId, "Re: We're now powered by phpBB 3.3", fmt.Sprintf("Spam Post %d", i), user4Id)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	err = model.UpdateLastPostOfTopic(ctx, topicBB2Id, postBB2DId, user4Id, user4Name)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	err = model.UpdateLastPostOfForum(ctx, forumBBId, postBB2DId, "Re: We're now powered by phpBB 3.3", user4Id, user4Name)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	err = model.IncreaseNumPostsForUser(ctx, user4Id)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	err = model.IncreaseNumPostsForTopic(ctx, topicBB2Id)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	err = model.IncreaseNumPostsForForum(ctx, forumBBId)
			// 	if err != nil {
			// 		return err
			// 	}
			// }

			// Spam topic of "Now Hear This!" forum
			// for i := 3; i <= 250; i++ {
			// 	_, err := model.InsertTopic(ctx, forumBBId, fmt.Sprintf("Spam Topic %d", i), user2Id, user2Name)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	err = model.IncreaseNumTopicsForForum(ctx, forumBBId)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
		}
	}
	return nil
}
