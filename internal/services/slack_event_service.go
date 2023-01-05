package services

import (
	"iamprashant.in/apps/slack-app/pkg/commons"
	"iamprashant.in/apps/slack-app/pkg/connectors"
)

type SlackEventService interface{}

type slackEventService struct {
	logger commons.Logger
	psql   connectors.PostgresConnector
}

func NewSlackEventService(logger commons.Logger, psql connectors.PostgresConnector) SlackEventService {
	return &slackEventService{logger: logger, psql: psql}
}
