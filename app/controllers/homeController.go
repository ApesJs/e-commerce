package controllers

import (
	render2 "github.com/unrolled/render"
	"net/http"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	render := render2.New(render2.Options{
		Layout:     "layout",
		Extensions: []string{".html"},
	})

	_ = render.HTML(w, http.StatusOK, "home", map[string]interface{}{
		"title": "Home Title",
		"body":  "Home Body",
	})
}
