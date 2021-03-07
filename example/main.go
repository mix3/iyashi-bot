package main

import (
	"log"
	"os"

	iyashibot "github.com/mix3/iyashi-bot"
	"github.com/mix3/iyashi-bot/config"
)

func main() {
	log.Fatal(iyashibot.Run(
		config.SlackBotToken(os.Getenv("IYASHI_BOT_SLACK_BOT_TOKEN")),
		config.SlackSigningSecret(os.Getenv("IYASHI_BOT_SLACK_SIGNING_SECRET")),
		config.FlickrAPIToken(os.Getenv("IYASHI_BOT_FLICKR_API_TOKEN")),
		config.TumblrAPIToken(os.Getenv("IYASHI_BOT_TUMBLR_API_TOKEN")),
		config.MoeURL("https://example.com"),
		config.MoeKeys([]string{
			"xxx",
			"yyy",
			"zzz",
		}),
	))
}
