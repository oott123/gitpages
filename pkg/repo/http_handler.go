package repo

import "net/http"

func (r *Repo) HttpHandler() *http.Handler {
	return r.httpHandler
}
