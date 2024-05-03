package controllers

import (
	"e-commerce/app/models"
	"github.com/gorilla/mux"
	render2 "github.com/unrolled/render"
	"log"
	"net/http"
	"strconv"
)

func (server *Server) Products(w http.ResponseWriter, r *http.Request) {
	render := render2.New(render2.Options{
		Layout:     "layout",
		Extensions: []string{".html"},
	})

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page = 1
	}
	perPage := 9

	productModel := models.Product{}
	products, totalRows, err := productModel.GetProducts(server.DB, perPage, page)
	if err != nil {
		log.Fatal("error:", err)
		return
	}

	pagination, _ := GetPaginationLinks(PaginationParams{
		Path:        "products",
		TotalRows:   int32(totalRows),
		PerPage:     int32(perPage),
		CurrentPage: int32(page),
	})

	_ = render.HTML(w, http.StatusOK, "products", map[string]interface{}{
		"products":   products,
		"pagination": pagination,
	})
}

func (server *Server) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	render := render2.New(render2.Options{
		Layout:     "layout",
		Extensions: []string{".html"},
	})

	vars := mux.Vars(r)
	if vars["slug"] == "" {
		return
	}

	productModel := models.Product{}
	product, err := productModel.FindBySlug(server.DB, vars["slug"])
	if err != nil {
		return
	}

	_ = render.HTML(w, http.StatusOK, "product", map[string]interface{}{
		"product": product,
		"success": GetFlash(w, r, "success"),
		"error":   GetFlash(w, r, "error"),
	})
}
