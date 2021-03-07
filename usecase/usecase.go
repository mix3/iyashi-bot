package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/mix3/iyashi-bot/domain/repository"
)

type Usecase interface {
	Run(ctx context.Context, channel, user string, args []string)
}

type usecase struct {
	repo     repository.Repository
	commands []Command
}

func NewUsecase(repo repository.Repository) Usecase {
	cmds := []Command{
		newMoeCommand(repo),
		newIyashiCommand(repo),
		newTumblrCommand(repo, "grass-tree-garden", []string{"しばき"}, []string{}, false),
		newTumblrCommand(repo, "honobonoarc", []string{"萌え"}, []string{}, true),
		newTumblrCommand(repo, "ganbaruzoi", []string{"ぞい"}, []string{}, false),
		newTumblrCommand(repo, "tawawa-of-monday", []string{"たわわ"}, []string{"safe"}, false),
	}
	helpcmd := newHelpCommand(repo.SlackAPI(), cmds)
	return &usecase{
		repo:     repo,
		commands: append(cmds, helpcmd),
	}
}

func (u *usecase) Run(ctx context.Context, channel, user string, args []string) {
	defer func() {
		if err := recover(); err != nil {
			u.err(ctx, channel, user, fmt.Errorf("panic: %w", err))
		}
	}()
	if err := u.run(ctx, channel, user, args); err != nil {
		u.err(ctx, channel, user, err)
	}
}

func (u *usecase) run(ctx context.Context, channel, user string, args []string) error {
	for _, c := range u.commands {
		if c.Match(args[0]) {
			return c.Execute(ctx, channel, user, args[1:])
		}
	}
	return nil
}

func (u *usecase) err(ctx context.Context, channel, user string, err error) {
	log.Printf("[WARN] channel=%s user=%s err:%s", channel, user, err)
	u.repo.SlackAPI().Reply(ctx, channel, user, fmt.Sprintf("エラっちゃった(´・ω・｀) err:%s", err))
}
