package controllers

import (
	"e-commerce/app/models"
	render2 "github.com/unrolled/render"
	"log"
	"net/http"
)

func (server *Server) Products(w http.ResponseWriter, r *http.Request) {
	render := render2.New(render2.Options{
		Layout: "layout",
	})

	productModel := models.Product{}
	products, err := productModel.GetProducts(server.DB)
	if err != nil {
		log.Fatal("error:", err)
		return
	}

	_ = render.HTML(w, http.StatusOK, "products", map[string]interface{}{
		"products": products,
	})
}
