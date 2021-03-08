package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/mix3/iyashi-bot/config"
	"github.com/mix3/iyashi-bot/usecase"

	"github.com/mattn/go-shellwords"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var (
	re = regexp.MustCompile(`^<@.+?>(.+)`)
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	signingSecret string
	usecase       usecase.Usecase
}

func NewHandler(conf config.Config, u usecase.Usecase) Handler {
	return &handler{
		signingSecret: conf.SlackSigningSecret(),
		usecase:       u,
	}
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, h.signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		if num, ok := r.Header["X-Slack-Retry-Num"]; ok {
			log.Printf("[WARN] X-Slack-Retry-Num:%s X-Slack-Retry-Reason:%s", num, r.Header["X-Slack-Retry-Reason"])
			return
		}
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			log.Printf("[INFO] channel=%s user=%s text=%s", ev.Channel, ev.User, ev.Text)
			text := strings.ReplaceAll(ev.Text, "\u00A0", " ") // コピペするとスペースが non-breaking space になるっぽいので変換
			text = re.ReplaceAllString(text, "$1")             // 自分宛の文言 @<XXXXXX> 削る
			args, err := shellwords.Parse(text)
			if err != nil {
				log.Printf("[ERROR] %s", err)
				return
			}
			log.Printf("[INFO] Run args=%v", args)
			h.usecase.Run(r.Context(), ev.Channel, ev.User, args)
		}
	}
}
