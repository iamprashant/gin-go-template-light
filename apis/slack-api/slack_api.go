package slack_api

import (
	"github.com/gin-gonic/gin"
	"iamprashant.in/apps/slack-app/internal/services"
	"iamprashant.in/apps/slack-app/pkg/commons"
)

type SlackApi interface {
	Subscribe(*gin.Context)
}
type slackApi struct {
	logger  commons.Logger
	service services.SlackEventService
}

func NewSlackApi(logger commons.Logger, service services.SlackEventService) SlackApi {
	return &slackApi{logger: logger, service: service}
}

func (sa *slackApi) Subscribe(c *gin.Context) {
	sa.logger.Infof("got request for subscribe")
}
