package main

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatal().Err(err).Msgf("failed to init configuration")
	}
	initLogger(config)
	bot := NewBot(config)

	if err := bot.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run bot")
	}
}

type logger struct{}

func (l logger) Output(_ int, msg string) error {
	log.Info().Msg(msg)
	return nil
}

func initLogger(c Config) {
	if !c.isProdLike() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
}
