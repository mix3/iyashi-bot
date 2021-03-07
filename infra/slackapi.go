package infra

import (
	"context"
	"fmt"

	"github.com/mix3/iyashi-bot/domain/repository"
	"github.com/slack-go/slack"
)

type slackAPI struct {
	api    *slack.Client
	userID string
}

func newSlackAPI(api *slack.Client) (repository.SlackAPI, error) {
	res, err := api.AuthTest()
	if err != nil {
		return nil, err
	}
	return &slackAPI{
		api:    api,
		userID: res.UserID,
	}, nil
}

func (s *slackAPI) PostMessage(ctx context.Context, channel, text string) error {
	_, _, err := s.api.PostMessageContext(ctx, channel, slack.MsgOptionText(text, true))
	return err
}

func (s *slackAPI) DirectMessage(ctx context.Context, user, text string) error {
	_, _, err := s.api.PostMessageContext(ctx, user, slack.MsgOptionText(text, true))
	return err
}

func (s *slackAPI) Reply(ctx context.Context, channel, user, text string) error {
	_, _, err := s.api.PostMessageContext(ctx, channel, slack.MsgOptionText(fmt.Sprintf("<@%s> %s", user, text), false))
	return err
}

func (s *slackAPI) UserID() string {
	return s.userID
}
