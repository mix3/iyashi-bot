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

type flickrSearcher struct {
	token string
}

func newFlickrSearcher(token string) repository.FlickrSearcher {
	return &flickrSearcher{
		token: token,
	}
}

var (
	flickrDefaultWords = []string{
		"猫", "ねこ",
		"犬", "いぬ",
		"兎", "うさぎ",
		"鳥", "とり",
		"ハムスター",
		"パンダ",
		"日本酒",
	}
)

func (f *flickrSearcher) search(ctx context.Context, keywords []string, page int) (*flickrSearchResponse, error) {
	const baseUrl = "https://api.flickr.com/services/rest/"

	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	args := keywords
	if len(args) == 0 {
		args = append(args, flickrDefaultWords[rand.Intn(len(flickrDefaultWords))])
	}
	args = append(args, "-hentai", "-porn", "-sexy", "-fuck")

	params := url.Values{}
	params.Set("api_key", f.token)
	params.Set("format", "json")
	params.Set("nojsoncallback", "1")
	params.Set("method", "flickr.photos.search")
	params.Set("text", strings.Join(args, " "))
	params.Set("safe_mode", "1")
	params.Set("media", "photo")
	if 0 < page {
		params.Set("page", strconv.Itoa(page))
	}

	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res *flickrSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}

type flickrSearchResponse struct {
	Photos struct {
		Page    int `json:"page"`
		Pages   int `json:"pages"`
		PerPage int `json:"perpage"`
		Total   int `json:"total"`
		Photo   []struct {
			Id       string `json:"id"`
			Owner    string `json:"owner"`
			Secret   string `json:"secret"`
			Server   string `json:"server"`
			Farm     int    `json:"farm"`
			Title    string `json:"title"`
			Ispublic int    `json:"ispublic"`
			Isfriend int    `json:"isfriend"`
			Isfamily int    `json:"isfamily"`
		} `json:"photo"`
	} `json:"photos"`
}

func (f *flickrSearchResponse) RandomImageURL() string {
	if n := len(f.Photos.Photo); 0 < n {
		photo := f.Photos.Photo[rand.Intn(n)]
		return fmt.Sprintf(
			`https://farm%d.staticflickr.com/%s/%s_%s.jpg`,
			photo.Farm,
			photo.Server,
			photo.Id,
			photo.Secret,
		)
	}
	return ""
}

func (f *flickrSearcher) RandomSearch(ctx context.Context, keywords []string) (repository.FlickrRandomSearchResponse, error) {
	const limitPageNum = 40

	res, err := f.search(ctx, keywords, 0)
	if err != nil {
		return nil, err
	}

	pageRange := res.Photos.Pages
	if limitPageNum < pageRange {
		pageRange = limitPageNum
	}

	var imageURL string
	for i := 0; i < 3; i++ {
		res, err = f.search(ctx, keywords, rand.Intn(pageRange+1))
		if err != nil {
			return nil, err
		}
		imageURL = res.RandomImageURL()
		if imageURL != "" {
			break
		}
	}
	if imageURL == "" {
		return nil, repository.ErrorNotFound
	}

	return &flickrRandomSearchResponse{
		imageURL: imageURL,
	}, nil
}

type flickrRandomSearchResponse struct {
	imageURL string
}

func (f *flickrRandomSearchResponse) ImageURL() string {
	return f.imageURL
}
