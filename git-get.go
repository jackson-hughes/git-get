package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/jackson-hughes/git-get/internal/urls"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func displayHelp() {
	helpMessage := `git-get is a simple utility that clones git repositories into a go get style directory structure.

Usage:
	git-get <repository-url>

Example:
	git-get https://github.com/jackson-hughes/git-get.git

Flags:
	-h, --help	display this help message
	--version	print the version and exit

Configuration:
	GIT_GET_DIR	(required) root directory to clone repositories into
	GIT_GET_DEBUG	set to true to enable debug logging`
	fmt.Fprintln(os.Stderr, helpMessage)
}

func clone(rawURL string, targetDir string) error {
	cmd := exec.Command("git", "clone", rawURL, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flag.Usage = displayHelp
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	if flag.NArg() != 1 {
		displayHelp()
		os.Exit(2)
	}

	loadConfig()

	if err := validateGitExists(exec.LookPath); err != nil {
		log.Fatal().Err(err).Send()
	}

	rawURL := flag.Arg(0)

	// The parsed URL is used only to derive the target directory; the
	// original input is what gets passed to git clone.
	gitURL, err := urls.Parse(rawURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid repository URL")
	}

	projectFilepath, err := urls.GetFilepathFromURL(*gitURL, appConfig.Dir)
	if err != nil {
		log.Fatal().Err(err).Msg("Error determining target directory")
	}
	log.Debug().Msgf("input url: %v", rawURL)
	log.Debug().Msgf("path to write repo to: %v", projectFilepath)

	exists, err := checkDirectoryExists(projectFilepath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error checking target directory")
	}
	if exists {
		log.Fatal().Msgf("Target directory already exists: %s", projectFilepath)
	}

	if err := clone(rawURL, projectFilepath); err != nil {
		log.Fatal().Err(err).Msg("Error cloning repository")
	}
}
