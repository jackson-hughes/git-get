package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Specification struct {
	Debug bool
	Dir   string
}

var appConfig Specification

var version = "dev"

func loadConfig() {
	if err := envconfig.Process("git_get", &appConfig); err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}

	if appConfig.Dir == "" {
		log.Fatal().Msg(`GIT_GET_DIR is not set; set the clone root directory, e.g. export GIT_GET_DIR="$HOME/Projects"`)
	}

	if appConfig.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.With().Caller().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Debug().Msgf("version: %v", version)
	log.Debug().Msgf("config: %+v", appConfig)
}
