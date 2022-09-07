package main

import (
	"fmt"
	"git-get/internal"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"os/exec"
)

func displayHelp() {
	helpMessage := `git-get is a simple utility that clones git repositories into a go get style directory structure.

git-get accepts one argument, the git project url to clone, for example:
	git-get https://github.com/jhughes01/git-get.git`
	fmt.Println(helpMessage)
}

func clone(gitProjectURL url.URL, gitProjectsDir string) {
	cmd := exec.Command("git", "clone", gitProjectURL.String(), gitProjectsDir)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	log.Info().Msg(string(stdoutStderr))
}

func main() {
	if len(os.Args) != 2 {
		displayHelp()
		return
	}
	gitProjectURL := os.Args[1]

	var gitURL url.URL

	if ok := internal.IsScpSyntax(gitProjectURL); ok {
		URL, err := internal.ConvertScpURL(gitProjectURL)
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

	projectFilepath := internal.GetFilepathFromURL(gitURL, appConfig.ProjectsDir)
	log.Debug().Msgf("input url: %v", gitURL.String())
	log.Debug().Msgf("path to write repo to: %v", projectFilepath)

	clone(gitURL, projectFilepath)

}
