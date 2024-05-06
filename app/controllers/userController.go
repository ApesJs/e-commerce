package controllers

import (
	"e-commerce/app/models"
	render2 "github.com/unrolled/render"
	"net/http"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	render := render2.New(render2.Options{
		Layout:     "layout",
		Extensions: []string{".html"},
	})

	_ = render.HTML(w, http.StatusOK, "login", map[string]interface{}{
		"error": GetFlash(w, r, "error"),
	})
}

func (server *Server) DoLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	userModel := models.User{}
	//ubah nanti dimana pengecekan email dan password dalam 1 method saya di dalam model
	user, err := userModel.FindByEmail(server.DB, email)
	if err != nil {
		SetFlash(w, r, "error", "Email atau Password tidak sesuai")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !ComparePassword(password, user.Password) {
		SetFlash(w, r, "error", "Email atau Password tidak sesuai")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, sessionUser)
	session.Values["id"] = user.ID
	_ = session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionUser)
	session.Values["id"] = "kosong"
	_ = session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
