package config

import (
	"fmt"
)

type Option func(*config) error

func SlackBotToken(v string) Option {
	return func(c *config) error {
		if v == "" {
			return fmt.Errorf("SlackBotToken required")
		}
		c.slackBotToken = v
		return nil
	}
}

func SlackSigningSecret(v string) Option {
	return func(c *config) error {
		if v == "" {
			return fmt.Errorf("SlackSigningSecret required")
		}
		c.slackSigningSecret = v
		return nil
	}
}

func FlickrAPIToken(v string) Option {
	return func(c *config) error {
		if v == "" {
			return fmt.Errorf("FlickrAPIToken required")
		}
		c.flickrAPIToken = v
		return nil
	}
}

func TumblrAPIToken(v string) Option {
	return func(c *config) error {
		if v == "" {
			return fmt.Errorf("TumblrAPIToken required")
		}
		c.tumblrAPIToken = v
		return nil
	}
}

func MoeURL(v string) Option {
	return func(c *config) error {
		if v == "" {
			return fmt.Errorf("MoeURL required")
		}
		c.moeURL = v
		return nil
	}
}

func MoeKeys(v []string) Option {
	return func(c *config) error {
		c.moeKeys = v
		return nil
	}
}

type Config interface {
	SlackBotToken() string
	SlackSigningSecret() string
	FlickrAPIToken() string
	TumblrAPIToken() string
	MoeURL() string
	MoeKeys() []string
	Valid() error
}

type config struct {
	slackBotToken      string
	slackSigningSecret string
	flickrAPIToken     string
	tumblrAPIToken     string
	moeURL             string
	moeKeys            []string
}

func (c *config) SlackBotToken() string {
	return c.slackBotToken
}

func (c *config) SlackSigningSecret() string {
	return c.slackSigningSecret
}

func (c *config) FlickrAPIToken() string {
	return c.flickrAPIToken
}

func (c *config) TumblrAPIToken() string {
	return c.tumblrAPIToken
}

func (c *config) MoeURL() string {
	return c.moeURL
}

func (c *config) MoeKeys() []string {
	return c.moeKeys
}

func (c *config) Valid() error {
	if c.slackBotToken == "" {
		return fmt.Errorf("SlackBotToken required")
	}
	if c.slackSigningSecret == "" {
		return fmt.Errorf("SlackSigningSecret required")
	}
	if c.flickrAPIToken == "" {
		return fmt.Errorf("FlickrAPIToken required")
	}
	if c.tumblrAPIToken == "" {
		return fmt.Errorf("TumblrAPIToken required")
	}
	if c.moeURL == "" {
		return fmt.Errorf("MoeURL required")
	}
	return nil
}

func NewConfig(opts ...Option) (Config, error) {
	c := &config{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}
