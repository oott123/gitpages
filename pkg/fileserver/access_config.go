package fileserver

import "net/http"

type AccessConfig struct {
	AllowCORS          bool
	CORSOrigins        []string
	ContentType        string
	AllowDotFiles      bool
	AllowListDirectory bool
	HotlinkProtection  bool
	HotlinkOrigins     []string
	NotFoundErrorPage  string
	ForbiddenErrorPage string
}

func (f *FileServer) AccessConfig(r *http.Request) *AccessConfig {
	// TODO: add scripting or sth for further control
	a := f.accessConfig
	return &a
}
