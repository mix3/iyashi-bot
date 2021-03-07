package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mix3/iyashi-bot/domain/repository"
)

const (
	tumblrPageLimit = 20
)

type tumblrSearcher struct {
	token string
}

func newTumblrSearcher(token string) repository.TumblrSearcher {
	return &tumblrSearcher{
		token: token,
	}
}

func (t *tumblrSearcher) RandomSearch(ctx context.Context, tumblrID string, tags []string) (repository.TumblrRandomSearchResponse, error) {
	res, err := t.search(ctx, tumblrID, tags, 0)
	if err != nil {
		return nil, err
	}

	var photoURL string
	for i := 0; i < 3; i++ {
		n := res.Response.TotalPosts - tumblrPageLimit + 1
		if n < 0 {
			break
		}
		offset := rand.Intn(n)
		res, err = t.search(ctx, tumblrID, tags, offset)
		if err != nil {
			return nil, err
		}
		photoURL = res.RandomPhotoURL()
		if photoURL != "" {
			break
		}
	}
	if photoURL == "" {
		return nil, repository.ErrorNotFound
	}

	return &tumblrRandomSearchResponse{
		photoURL: photoURL,
	}, nil
}

func (t *tumblrSearcher) search(ctx context.Context, tumblrID string, tags []string, offset int) (*TumblrSearchResponse, error) {
	baseURL := fmt.Sprintf("http://api.tumblr.com/v2/blog/%s.tumblr.com/posts/photo", tumblrID)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("api_key", t.token)
	params.Set("limit", strconv.Itoa(tumblrPageLimit))
	if 0 < offset {
		params.Set("offset", strconv.Itoa(offset))
	}
	if 0 < len(tags) {
		params.Set("tag", strings.Join(tags, "+"))
	}

	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res *TumblrSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}

type TumblrSearchResponse struct {
	Response struct {
		Posts []struct {
			Photos []struct {
				OriginalSize struct {
					Url string `json:"url"`
				} `json:"original_size"`
			} `json:"photos"`
		} `json:"posts"`
		TotalPosts int `json:"total_posts"`
	} `json:"response"`
}

func (t *TumblrSearchResponse) RandomPhotoURL() string {
	urls := make([]string, 0, tumblrPageLimit)
	for _, post := range t.Response.Posts {
		for _, photo := range post.Photos {
			urls = append(urls, photo.OriginalSize.Url)
		}
	}
	if n := len(urls); 0 < n {
		return urls[rand.Intn(n)]
	}
	return ""
}

type tumblrRandomSearchResponse struct {
	photoURL string
}

func (t *tumblrRandomSearchResponse) PhotoURL() string {
	return t.photoURL
}
