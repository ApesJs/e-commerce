package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (server *Server) initializeRoutes() {
	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/", server.Home).Methods("GET")
	server.Router.HandleFunc("/products", server.Products).Methods("GET")
	server.Router.HandleFunc("/product/{slug}", server.GetProductBySlug).Methods("GET")

	server.Router.HandleFunc("/carts", server.GetCart).Methods("GET")
	server.Router.HandleFunc("/carts", server.AddItemToCart).Methods("POST")

	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/public/", http.FileServer(staticFileDirectory))
	server.Router.PathPrefix("/public/").Handler(staticFileHandler).Methods("GET")
}
