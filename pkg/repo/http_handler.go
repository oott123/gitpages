package repo

import "net/http"

func (r *Repo) HttpHandler() http.Handler {
	return r.httpHandler
}

func (r *Repo) SetHttpHandler(h http.Handler) {
	r.httpHandler = h
}
