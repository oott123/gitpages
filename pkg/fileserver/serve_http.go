package fileserver

import (
	"net/http"
	"os"
)

var _ http.Handler = (*FileServer)(nil)

func (f *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := f.Resolve(r.URL.Path)
	stat, statErr := os.Stat(filename)
	accessConfig := f.AccessConfig(r)
	f.log.Debugf("parsed access config: %#v", accessConfig)

	if filename == "" {
		// TODO: 403
	} else if os.IsNotExist(statErr) {
		// TODO: 404
	} else if stat.IsDir() {
		// TODO: directory
	} else {
		http.ServeFile(w, r, filename)
	}
}
