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

	logger.Infof(ctx, "Topics of 'Now Hear This!':")
	topics, err := model.ListTopics(ctx, 10)
	if err != nil {
		logger.Errorf(ctx, "Error while listing topics: %s", err)
	}
	for _, topic := range topics {
		logger.Infof(ctx, "%s", helper.JsonDumps(topic))
	}
	logger.Infof(ctx, "")
}
