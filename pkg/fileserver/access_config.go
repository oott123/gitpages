package fileserver

import "net/http"

type AccessConfig struct {
	// if AllowCORS is true, a Access-Control-Allow-Origin header will be set. see CORSOrigins for details.
	AllowCORS bool
	// when CORSOrigin is set, only if the request contains Origin header which matches one of the list will get ACAO header.
	// if CORSOrigin is empty list, ACAO: * will be set.
	CORSOrigins []string
	// Don't use for now
	ContentType string
	// if AllowDotFiles is true, dot files (files started with `.`) will allowed to access by user.
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
