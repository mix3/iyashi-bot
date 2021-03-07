package repository

import "context"

type RepositoryError string

func (r RepositoryError) Error() string {
	return string(r)
}

const (
	ErrorNotFound RepositoryError = "ErrorNotFound"
)

type Repository interface {
	SlackAPI() SlackAPI
	FlickrSearcher() FlickrSearcher
	TumblrSearcher() TumblrSearcher
	MoeSearcher() MoeSearcher
}

type SlackAPI interface {
	DirectMessage(ctx context.Context, user, text string) error
	PostMessage(ctx context.Context, channel, text string) error
	Reply(ctx context.Context, channel, user, text string) error
	UserID() string
}

type FlickrSearcher interface {
	RandomSearch(ctx context.Context, keywords []string) (FlickrRandomSearchResponse, error)
}

type FlickrRandomSearchResponse interface {
	ImageURL() string
}

type TumblrSearcher interface {
	RandomSearch(ctx context.Context, tumblrID string, tags []string) (TumblrRandomSearchResponse, error)
}

type TumblrRandomSearchResponse interface {
	PhotoURL() string
}

type MoeSearcher interface {
	RandomSearch(ctx context.Context) (MoeRandomSearchResponse, error)
}

type MoeRandomSearchResponse interface {
	ImageURL() string
}
