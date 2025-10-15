package urls

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// IsScpSyntax takes a string url and returns true if the url is in scp format
func IsScpSyntax(url string) bool {
	scpSyntax := regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)
	match := scpSyntax.FindStringSubmatch(url)
	if match != nil {
		return true
	}
	return false
}

// ConvertScpURL converts scp syntax urls into ssh transport urls
// git@github.com:jackson-hughes/git-get.git -> ssh://git@github.com/jackson-hughes/git-get
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
	// trim port from filepath
	if strings.Contains(gitHost, ":") {
		gitHost = strings.Split(url.Host, ":")[0]
		log.Debug().Msgf("port found in hostname, %v has been replaced with %v", url.Host, gitHost)
	}
	gitProject := strings.TrimSuffix(url.Path, ".git")
	gitProject = strings.TrimPrefix(gitProject, "/")
	pathParts := []string{gitProjectRoot, gitHost}
	if gitProject != "" {
		pathParts = append(pathParts, strings.Split(gitProject, "/")...)
	}
	return filepath.Join(pathParts...)
}
