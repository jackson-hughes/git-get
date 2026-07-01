package urls

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var scpSyntaxRegex = regexp.MustCompile(`^([a-zA-Z0-9._-]+)@([a-zA-Z0-9._-]+):(.*)$`)

// Parse parses a user-supplied repository URL, converting scp syntax urls
// (git@github.com:org/repo.git) to ssh transport form. It returns an error
// if the url cannot be parsed or has no host.
func Parse(rawURL string) (*url.URL, error) {
	var u *url.URL
	var err error
	if isScpSyntax(rawURL) {
		u, err = convertScpURL(rawURL)
	} else {
		u, err = url.Parse(rawURL)
	}
	if err != nil {
		return nil, err
	}
	if u.Host == "" {
		return nil, fmt.Errorf("could not determine host from %q; expected a URL like https://github.com/org/repo.git or git@github.com:org/repo.git", rawURL)
	}
	return u, nil
}

// isScpSyntax takes a string url and returns true if the url is in scp format
func isScpSyntax(url string) bool {
	return scpSyntaxRegex.MatchString(url)
}

// convertScpURL converts scp syntax urls into ssh transport urls
// git@github.com:jackson-hughes/git-get.git -> ssh://git@github.com/jackson-hughes/git-get
func convertScpURL(scpSyntaxUrl string) (*url.URL, error) {
	log.Debug().Msgf("convertScpURL: received input %v", scpSyntaxUrl)
	path := strings.Replace(scpSyntaxUrl, ":", "/", 1)
	path = strings.TrimSuffix(path, ".git")
	convertedURL, err := url.Parse(fmt.Sprintf("ssh://%v", path))
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("convertScpURL: returned %v", convertedURL)
	return convertedURL, nil
}

// GetFilepathFromURL determines the go get style filepath based on the git url.
// It returns an error if the resulting path would escape gitProjectRoot.
func GetFilepathFromURL(u url.URL, gitProjectRoot string) (string, error) {
	gitHost := u.Hostname()
	gitProject := strings.TrimSuffix(u.Path, ".git")
	gitProject = strings.TrimPrefix(gitProject, "/")
	pathParts := []string{gitProjectRoot, gitHost}
	if gitProject != "" {
		pathParts = append(pathParts, strings.Split(gitProject, "/")...)
	}
	target := filepath.Join(pathParts...)

	rel, err := filepath.Rel(gitProjectRoot, target)
	if err != nil || rel == "." || !filepath.IsLocal(rel) {
		return "", fmt.Errorf("target path %s escapes the project root %s", target, gitProjectRoot)
	}
	return target, nil
}
