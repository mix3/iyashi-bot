package infra

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/mix3/iyashi-bot/domain/repository"
)

type moeSearcher struct {
	moeURL string
	keys   []string
}

func newMoeSearcher(moeURL string, keys []string) repository.MoeSearcher {
	return &moeSearcher{
		moeURL: moeURL,
		keys:   keys,
	}
}

func (m *moeSearcher) RandomSearch(ctx context.Context) (repository.MoeRandomSearchResponse, error) {
	if len(m.keys) == 0 {
		return nil, repository.ErrorNotFound
	}
	return &moeRandomSearchResponse{
		imageURL: fmt.Sprintf("%s/%s", m.moeURL, m.keys[rand.Intn(len(m.keys))]),
	}, nil
}

type moeRandomSearchResponse struct {
	imageURL string
}

func (m *moeRandomSearchResponse) ImageURL() string {
	return m.imageURL
}
