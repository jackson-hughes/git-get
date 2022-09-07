package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type Specification struct {
	Debug bool
	Dir   string `default:"/tmp"`
}

var appConfig Specification

var vcsCommit string
var vcsURL string

func init() {
	if err := envconfig.Process("git_get", &appConfig); err != nil {
		log.Fatal().Msgf("error loading config:\n %v", err)
	}

	if appConfig.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.With().Caller().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Debug().Msgf("vcs url: %v\tvcs commit: %v", vcsURL, vcsCommit)
	log.Debug().Msgf("config: %+v", appConfig)
}
