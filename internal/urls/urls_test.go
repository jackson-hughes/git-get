package urls

import (
	"net/url"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		rawURL   string
		wantHost string
		wantErr  bool
	}{
		{"https url", "https://github.com/jackson-hughes/git-get.git", "github.com", false},
		{"scp url", "git@github.com:jackson-hughes/git-get.git", "github.com", false},
		{"ssh transport url", "ssh://git@github.com/jackson-hughes/git-get", "github.com", false},
		{"plain string with no host", "not a url at all", "", true},
		{"bare path with no host", "/some/local/path", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.rawURL)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Parse(%q) expected error, got %v", tt.rawURL, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Host != tt.wantHost {
				t.Errorf("Parse(%q).Host = %q, want %q", tt.rawURL, got.Host, tt.wantHost)
			}
		})
	}
}

func TestIsScpSyntax(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"scp url", "git@github.com:jackson-hughes/git-get.git", true},
		{"scp url with hyphenated user", "deploy-bot@github.com:jackson-hughes/git-get.git", true},
		{"scp url with dotted user", "first.last@github.com:jackson-hughes/git-get.git", true},
		{"ssh transport url", "ssh://git@github.com/jackson-hughes/git-get", false},
		{"https url", "https://github.com/jackson-hughes/git-get.git", false},
		{"plain string", "not a url at all", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isScpSyntax(tt.url); got != tt.want {
				t.Errorf("isScpSyntax(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

func TestConvertScpURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"github scp url", "git@github.com:jackson-hughes/git-get.git", "ssh://git@github.com/jackson-hughes/git-get"},
		{"scp url without .git suffix", "git@github.com:jackson-hughes/git-get", "ssh://git@github.com/jackson-hughes/git-get"},
		{"scp url with hyphenated user", "deploy-bot@my-host.com:org/repo.git", "ssh://deploy-bot@my-host.com/org/repo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertScpURL(tt.url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.want {
				t.Errorf("convertScpURL(%q) = %q, want %q", tt.url, got.String(), tt.want)
			}
		})
	}
}

func TestGetFilepathFromURL(t *testing.T) {
	root := "/tmp/root"

	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name: "https url",
			url:  "https://github.com/jackson-hughes/git-get.git",
			want: filepath.Join(root, "github.com", "jackson-hughes", "git-get"),
		},
		{
			name: "host with port",
			url:  "https://github.com:8443/org/repo.git",
			want: filepath.Join(root, "github.com", "org", "repo"),
		},
		{
			name: "ipv6 host with port",
			url:  "ssh://git@[::1]:2222/org/repo.git",
			want: filepath.Join(root, "::1", "org", "repo"),
		},
		{
			name: "ssh transport url",
			url:  "ssh://git@github.com/jackson-hughes/git-get",
			want: filepath.Join(root, "github.com", "jackson-hughes", "git-get"),
		},
		{
			name:    "path traversal escapes root",
			url:     "https://github.com/../../../../home/user/.ssh",
			wantErr: true,
		},
		{
			name:    "path traversal resolves to root itself",
			url:     "https://github.com/..",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Fatalf("unexpected error parsing url: %v", err)
			}

			got, err := GetFilepathFromURL(*u, root)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("GetFilepathFromURL(%q) expected error, got path %q", tt.url, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("GetFilepathFromURL(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}
