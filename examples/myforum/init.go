package myforum

import (
	"context"

	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

func InitMyforum(ctx context.Context) {
	err := PopulateDb(ctx)
	if err != nil {
		logger.Fatalf(ctx, "Error while populating database: %s", err)
	}
}

func DebugMyforum(ctx context.Context) {
	logger.Infof(ctx, "Forums:")
	forums, err := model.ListForums(ctx)
	if err != nil {
		logger.Errorf(ctx, "Error while listing forums: %s", err)
	}
	for _, forum := range forums {
		logger.Infof(ctx, "%s", helper.JsonDumps(forum))
	}
	logger.Infof(ctx, "")

	logger.Infof(ctx, "Hierarchy of Root Forum:")
	forumChildNodes := model.ComputeForumChildNodes(ctx, forums, model.ROOT_FORUM_ID, 0)
	for _, forumNode := range forumChildNodes {
		logger.Infof(ctx, "  - %t %s", forumNode.IsLeaf, helper.JsonDumps(forumNode.Forum))
		if len(forumNode.ForumChildNodes) > 0 {
			for _, forumNode := range forumNode.ForumChildNodes {
				logger.Infof(ctx, "      - %t %s", forumNode.IsLeaf, helper.JsonDumps(forumNode.Forum))
				if len(forumNode.ForumChildNodes) > 0 {
					for _, forumNode := range forumNode.ForumChildNodes {
						logger.Infof(ctx, "          - %t %s", forumNode.IsLeaf, helper.JsonDumps(forumNode.Forum))
					}
				}
			}
		}
	}
	logger.Infof(ctx, "")

	FORUM_ID := 10
	logger.Infof(ctx, "Navigation trails of 'Now Hear This!' forum:")
	forumNavTrails, err := model.ComputeForumNavTrails(ctx, FORUM_ID)
	if err != nil {
		logger.Errorf(ctx, "Error while computing Forum Nav Trails for forum id %d: %s", FORUM_ID, err)
	}
	for _, forum := range forumNavTrails {
		logger.Infof(ctx, "  - %s", helper.JsonDumps(forum))
	}
	logger.Infof(ctx, "")

	TOPIC_ID := 2
	logger.Infof(ctx, "Users of 'We're now powered by phpBB 3.3':")
	users, err := model.ListUsers(ctx, TOPIC_ID)
	if err != nil {
		logger.Errorf(ctx, "Error while listing users: %s", err)
	}
	for _, user := range users {
		logger.Infof(ctx, "%s", helper.JsonDumps(user))
	}
	logger.Infof(ctx, "")

	logger.Infof(ctx, "Topics of 'Now Hear This!':")
	topics, err := model.ListTopics(ctx, FORUM_ID, 0)
	if err != nil {
		logger.Errorf(ctx, "Error while listing topics: %s", err)
	}
	for _, topic := range topics {
		logger.Infof(ctx, "%s", helper.JsonDumps(topic))
	}
	logger.Infof(ctx, "")

	logger.Infof(ctx, "Posts of 'We're now powered by phpBB 3.3':")
	posts, err := model.ListPosts(ctx, TOPIC_ID, 0)
	if err != nil {
		logger.Errorf(ctx, "Error while listing posts: %s", err)
	}
	for _, post := range posts {
		logger.Infof(ctx, "%s", helper.JsonDumps(post))
	}
	logger.Infof(ctx, "")
}
