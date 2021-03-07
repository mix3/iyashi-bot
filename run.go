package iyashibot

import (
	crand "crypto/rand"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"os"

	"github.com/fujiwara/ridge"
	"github.com/hashicorp/logutils"
	"github.com/mix3/iyashi-bot/config"
	"github.com/mix3/iyashi-bot/handler"
	"github.com/mix3/iyashi-bot/infra"
	"github.com/mix3/iyashi-bot/usecase"
)

func init() {
	log.SetOutput(&logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	})

	seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	rand.Seed(seed.Int64())
}

func Run(opts ...config.Option) error {
	conf, err := config.NewConfig(opts...)
	if err != nil {
		return err
	}
	if err := conf.Valid(); err != nil {
		return err
	}

	repo, err := infra.NewRepository(conf)
	if err != nil {
		return err
	}

	uc := usecase.NewUsecase(repo)
	h := handler.NewHandler(conf, uc)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Index)
	log.Println("[INFO] Server listening")
	ridge.Run(":8080", "/", mux)
	return nil
}
