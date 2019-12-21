package config

import (
	"github.com/minio/minio/pkg/wildcard"
	"strings"
)

type Server struct {
	// Host is the server matched http Host header; '*' is allowed for wildcard match (e.g. ex*mple.example.org)
	Host string
	// Remote is the git repo to serve; e.g. https://github.com/ghost/example.git
	Remote string
	// WebHookSecret is the secret used by web hook; your web hook endpoint will be your host /_gitpages/update/{your_secret}
	WebHookSecret string
	// Branch is the deploy target branch from the git repo
	Branch string
	// Dir is the serve root
	Dir string
	// AllowSymlink will allow symlinks point to files out of the git repo
	AllowSymlink bool
}

func (s *Server) MatchHost(host string) bool {
	return wildcard.MatchSimple(strings.ToLower(s.Host), strings.ToLower(host))
}
