package fileserver

import "net/http"

type AccessConfig struct {
	// if AllowCORS is true, a Access-Control-Allow-Origin header will be set. see CORSOrigins for details.
	AllowCORS bool
	// when CORSOrigin is set, only if the request contains Origin header which matches one of the list will get ACAO header.
	// if CORSOrigin is empty list, ACAO: * will be set.
	CORSOrigins []string
	// if AllowDotFiles is true, dot files (files started with `.`) will allowed to access by user.
	AllowDotFiles bool
	// if AllowListDirectory is true, access directory which don't contains index files will result a simple list
	AllowListDirectory bool
	// if NotFoundErrorPage is defined, 404 will be sent to this file (relative to site root); else will be sent to 404.html
	NotFoundErrorPage string
	// Don't use for now
	ContentType string
	// Don't use for now
	HotlinkProtection bool
	// Don't use for now
	HotlinkOrigins []string
	// Don't use for now
	ForbiddenErrorPage string
}

func (f *FileServer) AccessConfig(r *http.Request) *AccessConfig {
	// TODO: add scripting or sth for further control
	a := f.accessConfig
	return &a
}
