package infra

import (
	"github.com/mix3/iyashi-bot/config"
	"github.com/mix3/iyashi-bot/domain/repository"

	"github.com/slack-go/slack"
)

type store struct {
	slackAPI       repository.SlackAPI
	flickrSearcher repository.FlickrSearcher
	tumblrSearcher repository.TumblrSearcher
	moeSearcher    repository.MoeSearcher
}

func NewRepository(conf config.Config) (repository.Repository, error) {
	api := slack.New(conf.SlackBotToken())
	slackAPI, err := newSlackAPI(api)
	if err != nil {
		return nil, err
	}
	return &store{
		slackAPI:       slackAPI,
		flickrSearcher: newFlickrSearcher(conf.FlickrAPIToken()),
		tumblrSearcher: newTumblrSearcher(conf.TumblrAPIToken()),
		moeSearcher:    newMoeSearcher(conf.MoeURL(), conf.MoeKeys()),
	}, nil
}

func (r *store) SlackAPI() repository.SlackAPI {
	return r.slackAPI
}

func (r *store) FlickrSearcher() repository.FlickrSearcher {
	return r.flickrSearcher
}

func (r *store) TumblrSearcher() repository.TumblrSearcher {
	return r.tumblrSearcher
}

func (r *store) MoeSearcher() repository.MoeSearcher {
	return r.moeSearcher
}
