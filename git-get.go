package main

import (
	"fmt"
	"git-get/internal/urls"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"os/exec"
)

func displayHelp() {
	helpMessage := `git-get is a simple utility that clones git repositories into a go get style directory structure.

git-get accepts one argument, the git project url to clone, for example:
	git-get https://github.com/jackson-hughes/git-get.git`
	fmt.Println(helpMessage)
}

func clone(gitProjectURL url.URL, gitProjectsDir string) error {
	cmd := exec.Command("git", "clone", gitProjectURL.String(), gitProjectsDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		displayHelp()
		return
	}
	gitProjectURL := os.Args[1]

	var gitURL url.URL

	if ok := urls.IsScpSyntax(gitProjectURL); ok {
		URL, err := urls.ConvertScpURL(gitProjectURL)
		if err != nil {
			log.Fatal().Err(err)
		}
		gitURL = *URL
	} else {
		URL, err := url.Parse(gitProjectURL)
		if err != nil {
			log.Fatal().Err(err)
		}
		gitURL = *URL
	}

	projectFilepath := urls.GetFilepathFromURL(gitURL, appConfig.Dir)
	log.Debug().Msgf("input url: %v", gitURL.String())
	log.Debug().Msgf("path to write repo to: %v", projectFilepath)

	if err := clone(gitURL, projectFilepath); err != nil {
		log.Err(err)
	}
}
