package controllers

import (
	"e-commerce/app/models"
	"github.com/google/uuid"
	render2 "github.com/unrolled/render"
	"net/http"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	render := render2.New(render2.Options{
		Layout:     "layoutLogin",
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

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	render := render2.New(render2.Options{
		Layout:     "layoutLogin",
		Extensions: []string{".html"},
	})

	_ = render.HTML(w, http.StatusOK, "register", map[string]interface{}{
		"error": GetFlash(w, r, "error"),
	})
}

func (server *Server) DoRegister(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if firstName == "" || lastName == "" || email == "" || password == "" {
		SetFlash(w, r, "error", "Semua form harus di isi")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	userModel := models.User{}
	existUser, _ := userModel.FindByEmail(server.DB, email)

	if existUser != nil {
		SetFlash(w, r, "error", "email sudah di gunakan")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	hashedPassword, err := MakePassword(password)
	params := &models.User{
		ID:        uuid.New().String(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hashedPassword,
	}

	user, err := userModel.CreateUser(server.DB, params)
	if err != nil {
		SetFlash(w, r, "error", "registrasi gagal")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
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
