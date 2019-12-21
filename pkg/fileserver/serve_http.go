package fileserver

import (
	"io"
	"mime"
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

	if !accessConfig.AllowDotFiles {
		pieces := strings.Split(filename, string(os.PathSeparator))
		for _, p := range pieces {
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
		if _, err := os.Stat(indexFilename); os.IsNotExist(err) {
			if accessConfig.AllowListDirectory {
				http.ServeFile(w, r, filename)
			} else {
				http.Error(w, "directory listing is not allowed", 403)
			}
		} else {
			http.ServeFile(w, r, filename)
		}
		return
	} else {
		http.ServeFile(w, r, filename)
		return
	}

	// 404 here
	if accessConfig.NotFoundErrorPage != "" {
		filename = f.Resolve(accessConfig.NotFoundErrorPage)
	} else {
		filename = f.Resolve("/404.html")
	}

	if filename != "" {
		ext := filepath.Ext(filename)
		mType := mime.TypeByExtension(ext)
		if mType == "" {
			mType = "text/plain"
		}
		w.Header().Set("Content-Type", mType)
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
