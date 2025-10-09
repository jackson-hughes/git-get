package urls

import (
	"net/url"
	"path/filepath"
	"testing"
)

func TestIsScpSyntax(t *testing.T) {
	t.Run("scpUrlTrue", func(t *testing.T) {
		url := "git@github.com:jackson-hughes/git-get.git"
		match := IsScpSyntax(url)

		if match != true {
			t.Errorf("got %v want %v", match, true)
		}
	})

	t.Run("notScpUrlTrue", func(t *testing.T) {
		url := "ssh://git@github.com/jackson-hughes/git-get"
		match := IsScpSyntax(url)

		if match == true {
			t.Errorf("got %v want %v", match, false)
		}
	})
}

func TestConvertScpUrl(t *testing.T) {
	url := "git@github.com:jackson-hughes/git-get.git"
	want := "ssh://git@github.com/jackson-hughes/git-get"

	got, err := ConvertScpURL(url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.String() != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetFilepathFromURL(t *testing.T) {
	root := "/tmp"
	repoURL, err := url.Parse("https://github.com/jackson-hughes/git-get.git")
	if err != nil {
		t.Fatalf("unexpected error parsing url: %v", err)
	}

	got := GetFilepathFromURL(*repoURL, root)
	want := filepath.Join(root, "github.com", "jackson-hughes", "git-get")

	if got != want {
		t.Errorf("GetFilepathFromURL() = %q, want %q", got, want)
	}
}
