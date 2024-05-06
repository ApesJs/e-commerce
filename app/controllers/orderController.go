package controllers

import (
	render2 "github.com/unrolled/render"
	"net/http"
)

func (server *Server) Checkout(w http.ResponseWriter, r *http.Request) {
	render := render2.New(render2.Options{
		Layout:     "layout",
		Extensions: []string{".html"},
	})

	if !IsLoggedIn(r) {
		SetFlash(w, r, "error", "Anda Perlu Login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	_ = render.HTML(w, http.StatusOK, "checkout", map[string]interface{}{
		"user": server.CurrentUser(w, r),
	})
}
