package controllers

import (
	render2 "github.com/unrolled/render"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render := render2.New(render2.Options{
		Layout: "layout",
	})

	_ = render.HTML(w, http.StatusOK, "home", map[string]interface{}{
		"title": "Home Title",
		"body":  "Home Body",
	})
}
