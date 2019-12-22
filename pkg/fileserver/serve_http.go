package fileserver

import (
	"fmt"
	gomime "github.com/cubewise-code/go-mime"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var _ http.Handler = (*FileServer)(nil)

func (f *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := f.Resolve(r.URL.Path)
	stat, statErr := os.Stat(filename)
	accessConfig := f.AccessConfig(r)
	f.log.Debugf("parsed access config: %#v", accessConfig)

	allowAccess := true
	varyHeader := make([]string, 0)

	if v := w.Header().Get("Vary"); v != "" {
		varyHeader = append(varyHeader, v)
	}

	if !accessConfig.AllowDotFiles && strings.Contains(filename, ".") {
		for _, p := range strings.FieldsFunc(filename, isSlashRune) {
			if strings.Index(p, ".") == 0 {
				// is dotfile
				f.log.Debugf("%s is dotfile, not allowed", filename)
				allowAccess = false
				break
			}
		}
	}

	if accessConfig.AllowCORS {
		if len(accessConfig.CORSOrigins) == 0 {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			varyHeader = append(varyHeader, "Origin")
			for _, o := range accessConfig.CORSOrigins {
				if o == r.Header.Get("Origin") {
					w.Header().Set("Access-Control-Allow-Origin", o)
					break
				}
			}
		}
	}

	if accessConfig.CrossSiteProtection || accessConfig.ReferrerProtection {
		varyHeader = append(varyHeader, "Referer")
	}

	if accessConfig.CrossSiteProtection {
		varyHeader = append(varyHeader, "Sec-Fetch-Dest", "Sec-Fetch-Site", "Sec-Fetch-Mode")
		if r.Header.Get("Referer") != "" || accessConfig.CrossSiteProtectionOnlyEmptyReferrer {
			fetchDest := r.Header.Get("Sec-Fetch-Dest")
			fetchSite := r.Header.Get("Sec-Fetch-Site")
			fetchMode := r.Header.Get("Sec-Fetch-Mode")
			if fetchSite == "cross-site" && fetchMode == "no-cors" {
				if fetchDest != "empty" && fetchDest != "document" {
					allowAccess = false
				}
			}
		}
	}

	if accessConfig.ReferrerProtection {
		referrer := r.Header.Get("Referer")
		allowAccess = false
		for _, r := range accessConfig.ReferrerAllowed {
			if r.MatchString(referrer) {
				allowAccess = true
			}
		}
	}

	if len(varyHeader) > 0 {
		w.Header().Set("Vary", strings.Join(varyHeader, ", "))
	}

	for k, v := range accessConfig.AddHeaders {
		w.Header().Set(k, v)
	}

	if filename == "" || !allowAccess {
		f.serveForbidden(w, r, "access denied", accessConfig)
		return
	} else if os.IsNotExist(statErr) || stat == nil {
		// 404 fall through
	} else if stat.IsDir() {
		// directory
		indexFilename := path.Join(filename, "index.html")
		if _, err := os.Stat(indexFilename); os.IsNotExist(err) && !accessConfig.AllowListDirectory {
			f.serveForbidden(w, r, "directory listing is not allowed", accessConfig)
		} else {
			f.serveFile(w, r, filename, stat, accessConfig)
		}
		return
	} else {
		// 200 OK for regular files
		f.serveFile(w, r, filename, stat, accessConfig)
		return
	}

	// 404 here
	if accessConfig.NotFoundErrorPage != "" {
		filename = f.Resolve(accessConfig.NotFoundErrorPage)
	} else {
		filename = f.Resolve("/404.html")
	}

	if _, err := os.Stat(filename); filename != "" && err == nil {
		f.setMime(w, filename, accessConfig)
		w.Header().Del("Date")
		w.Header().Del("Last-Modified")
		w.Header().Set("Cache-Control", "no-cache")

		w.WriteHeader(404)
		file, err := os.Open(filename)
		defer file.Close()
		if err != nil {
			f.log.Errorf("error while opening 404 file: %s", err)
			http.Error(w, "404 Not found", 404)
			return
		}
		_, _ = io.Copy(w, file)
		return
	}
}

func (f *FileServer) serveFile(w http.ResponseWriter, r *http.Request, filename string, stat os.FileInfo, accessConfig *AccessConfig) {
	if !stat.IsDir() {
		f.setMime(w, filename, accessConfig)
	}
	http.ServeFile(w, r, filename)
}

func (f *FileServer) serveForbidden(w http.ResponseWriter, r *http.Request, reason string, accessConfig *AccessConfig) {
	w.Header().Del("Date")
	w.Header().Del("Last-Modified")
	w.Header().Set("Cache-Control", "no-cache")

	filename := ""
	if accessConfig.ForbiddenErrorPage != "" {
		filename = f.Resolve(accessConfig.ForbiddenErrorPage)
	}

	if filename != "" {
		f.setMime(w, filename, accessConfig)
		w.WriteHeader(403)

		file, err := os.Open(filename)
		defer file.Close()
		if err != nil {
			f.log.Errorf("error while opening 403 file: %s", err)
			http.Error(w, "403 Forbidden", 403)
			return
		}
		_, _ = io.Copy(w, file)
	} else {
		http.Error(w, reason, 403)
	}
}

func (f *FileServer) setMime(w http.ResponseWriter, filename string, accessConfig *AccessConfig) {
	ext := filepath.Ext(filename)
	mType := gomime.TypeByExtension(ext)
	if mType == "" {
		mType = "text/plain"
	}

	if mType == "application/javascript" || strings.Contains(mType, "text/") {
		mType = fmt.Sprintf("%s; charset=utf-8", mType)
	}

	w.Header().Set("Content-Type", mType)
}
