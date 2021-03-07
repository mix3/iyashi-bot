package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/mix3/iyashi-bot/domain/repository"
)

type Command interface {
	MatchStrings() []string
	Match(str string) bool
	Help() string
	Execute(ctx context.Context, channel, user string, args []string) error
}

type helpCommand struct {
	slackAPI repository.SlackAPI
	commands []Command
}

func newHelpCommand(slackAPI repository.SlackAPI, commands []Command) Command {
	return &helpCommand{
		slackAPI: slackAPI,
		commands: commands,
	}
}

func (h *helpCommand) MatchStrings() []string {
	return []string{"help", "?"}
}

func (h *helpCommand) Match(str string) bool {
	for _, s := range h.MatchStrings() {
		if s == str {
			return true
		}
	}
	return false
}

func (h *helpCommand) Help() string {
	helps := make([]string, 0, len(h.commands))
	for _, c := range h.commands {
		helps = append(helps, fmt.Sprintf("%s: %s", strings.Join(c.MatchStrings(), "|"), c.Help()))
	}
	return fmt.Sprintf("```%s```", strings.Join(helps, "\n"))
}

func (h *helpCommand) Execute(ctx context.Context, channel, user string, args []string) error {
	if 0 < len(args) {
		for _, c := range h.commands {
			if c.Match(args[0]) {
				return h.slackAPI.Reply(ctx, channel, user, c.Help())
			}
		}
	}
	return h.slackAPI.Reply(ctx, channel, user, h.Help())
}

type moeCommand struct {
	slackAPI    repository.SlackAPI
	moeSearcher repository.MoeSearcher
}

func newMoeCommand(repo repository.Repository) Command {
	return &moeCommand{
		slackAPI:    repo.SlackAPI(),
		moeSearcher: repo.MoeSearcher(),
	}
}

func (m *moeCommand) MatchStrings() []string {
	return []string{"もえ"}
}

func (m *moeCommand) Match(str string) bool {
	for _, s := range m.MatchStrings() {
		if str == s {
			return true
		}
	}
	return false
}

func (m *moeCommand) Help() string {
	return "mix3 が溜め込んだ画像を返すよ！"
}

func (m *moeCommand) Execute(ctx context.Context, channel, user string, args []string) error {
	res, err := m.moeSearcher.RandomSearch(ctx)
	if err != nil {
		return err
	}
	return m.slackAPI.Reply(ctx, channel, user, res.ImageURL())
}

type iyashiCommand struct {
	slackAPI       repository.SlackAPI
	flickrSearcher repository.FlickrSearcher
}

func newIyashiCommand(repo repository.Repository) Command {
	return &iyashiCommand{
		slackAPI:       repo.SlackAPI(),
		flickrSearcher: repo.FlickrSearcher(),
	}
}

func (m iyashiCommand) MatchStrings() []string {
	return []string{"癒やし", "癒し"}
}

func (m *iyashiCommand) Match(str string) bool {
	for _, s := range m.MatchStrings() {
		if str == s {
			return true
		}
	}
	return false
}

func (m *iyashiCommand) Help() string {
	return "flicker から画像を返すよ！"
}

func (m *iyashiCommand) Execute(ctx context.Context, channel, user string, args []string) error {
	res, err := m.flickrSearcher.RandomSearch(ctx, args)
	if err != nil {
		if err == repository.ErrorNotFound {
			return m.slackAPI.Reply(ctx, channel, user, "見つかんなかったよ(´・ω・｀)")
		}
		return err
	}
	if err := m.slackAPI.DirectMessage(ctx, user, res.ImageURL()); err != nil {
		return err
	}
	return m.slackAPI.Reply(ctx, channel, user, "╭( ･ㅂ･)ﻭ ̑̑ DMしたよ")
}

type tumblrCommand struct {
	slackAPI       repository.SlackAPI
	tumblrSearcher repository.TumblrSearcher
	tumblrID       string
	matchStrings   []string
	appendTags     []string
	isDM           bool
}

func newTumblrCommand(repo repository.Repository, tumblrID string, matchStrings, appendTags []string, isDM bool) Command {
	return &tumblrCommand{
		slackAPI:       repo.SlackAPI(),
		tumblrSearcher: repo.TumblrSearcher(),
		tumblrID:       tumblrID,
		matchStrings:   matchStrings,
		appendTags:     appendTags,
		isDM:           isDM,
	}
}

func (t *tumblrCommand) MatchStrings() []string {
	return t.matchStrings
}

func (t *tumblrCommand) Match(str string) bool {
	for _, s := range t.MatchStrings() {
		if s == str {
			return true
		}
	}
	return false
}

func (t *tumblrCommand) Help() string {
	return fmt.Sprintf("http://%s.tumblr.com/ から画像をランダムで返すよ！", t.tumblrID)
}

func (t *tumblrCommand) Execute(ctx context.Context, channel, user string, args []string) error {
	res, err := t.tumblrSearcher.RandomSearch(ctx, t.tumblrID, append(args, t.appendTags...))
	if err != nil {
		if err == repository.ErrorNotFound {
			return t.slackAPI.Reply(ctx, channel, user, "見つかんなかったよ(´・ω・｀)")
		}
		return err
	}
	if t.isDM {
		if err := t.slackAPI.DirectMessage(ctx, user, res.PhotoURL()); err != nil {
			return err
		}
		return t.slackAPI.Reply(ctx, channel, user, "╭( ･ㅂ･)ﻭ ̑̑ DMしたよ")
	} else {
		return t.slackAPI.Reply(ctx, channel, user, res.PhotoURL())
	}
}
