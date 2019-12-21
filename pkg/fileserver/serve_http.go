package fileserver

import (
	"net/http"
	"os"
	"path"
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

	if !f.accessConfig.AllowDotFiles {
		baseName := path.Base(filename)
		if strings.Index(baseName, ".") == 0 {
			// is dotfile
			allowAccess = false
		}
	}

	if f.accessConfig.AllowCORS {
		if len(f.accessConfig.CORSOrigins) == 0 {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			varyHeader = append(varyHeader, "Origin")
			for _, o := range f.accessConfig.CORSOrigins {
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
		// TODO: 403
	} else if os.IsNotExist(statErr) {
		// TODO: 404
	} else if stat.IsDir() {
		// directory
		indexFilename := path.Join(filename, "index.html")
		if _, err := os.Stat(indexFilename); os.IsNotExist(err) {
			// TODO: directory listing
		} else {
			http.ServeFile(w, r, filename)
		}
	} else {
		http.ServeFile(w, r, filename)
	}
}
