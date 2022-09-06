package internal

import (
	"testing"
)

func TestIsScpSyntax(t *testing.T) {
	t.Run("scpUrlTrue", func(t *testing.T) {
		url := "git@github.com:jhughes01/git-get.git"
		match := IsScpSyntax(url)

		if match != true {
			t.Errorf("got %v want %v", match, true)
		}
	})

	t.Run("notScpUrlTrue", func(t *testing.T) {
		url := "ssh://git@github.com/jhughes01/git-get"
		match := IsScpSyntax(url)

		if match == true {
			t.Errorf("got %v want %v", match, false)
		}
	})
}

func TestConvertScpUrl(t *testing.T) {
	url := "git@github.com:jhughes01/git-get.git"
	want := "ssh://git@github.com/jhughes01/git-get"

	got, err := ConvertScpURL(url)
	if got.String() != want && err == nil {
		t.Errorf("got %v want %v", got, want)
	}
}
