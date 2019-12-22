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
	if containsDotDot(r.URL.Path) {
		// Too many programs use r.URL.Path to construct the argument to
		// serveFile. Reject the request under the assumption that happened
		// here and ".." may not be wanted.
		// Note that name might not contain "..", for example if code (still
		// incorrectly) used filepath.Join(myDir, r.URL.Path).
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}

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

	if len(varyHeader) > 0 {
		w.Header().Set("Vary", strings.Join(varyHeader, ", "))
	}

	if filename == "" || !allowAccess {
		http.Error(w, "access denied", 403)
		return
	} else if os.IsNotExist(statErr) {
		// 404 fall through
	} else if stat.IsDir() {
		// directory
		indexFilename := path.Join(filename, "index.html")
		if _, err := os.Stat(indexFilename); os.IsNotExist(err) && !accessConfig.AllowListDirectory {
			http.Error(w, "directory listing is not allowed", 403)
		} else {
			http.ServeFile(w, r, filename)
		}
		return
	} else {
		f.setMime(w, filename)
		http.ServeFile(w, r, filename)
		return
	}

	// 404 here
	if accessConfig.NotFoundErrorPage != "" {
		filename = f.Resolve(accessConfig.NotFoundErrorPage)
	} else {
		filename = f.Resolve("/404.html")
	}

	if _, err := os.Stat(filename); filename != "" && err == nil {
		f.setMime(w, filename)
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

func (f *FileServer) setMime(w http.ResponseWriter, filename string) {
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

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }
