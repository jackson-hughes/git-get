package internal

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/url"
	"regexp"
	"strings"
)

// IsScpSyntax takes a string url and returns true if the url is in scp format
func IsScpSyntax(url string) bool {
	scpSyntax := regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)
	match := scpSyntax.FindStringSubmatch(url)
	if match != nil {
		return true
	} else {
		return false
	}
}

// ConvertScpURL converts scp syntax urls into ssh transport urls
// git@github.com:jhughes01/git-get.git -> ssh://git@github.com/jhughes01/git-get
func ConvertScpURL(scpSyntaxUrl string) (*url.URL, error) {
	log.Debug().Msgf("ConvertScpURL: received input %v: ", scpSyntaxUrl)
	convertedUrl, err := url.Parse(fmt.Sprintf("ssh://%v", strings.Replace(
		strings.Replace(scpSyntaxUrl, ":", "/", 1),
		".git", "", 1)))
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("ConvertScpURL: returned %v: ", convertedUrl)
	return convertedUrl, nil
}

// GetFilepathFromURL determines the go get style filepath based on the git url
func GetFilepathFromURL(url url.URL, gitProjectRoot string) string {
	gitHost := url.Host
	gitProject := strings.Replace(url.Path, ".git", "", 1)
	path := fmt.Sprint(gitProjectRoot, "/", gitHost, gitProject)
	return path
}
