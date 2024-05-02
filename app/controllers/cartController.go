package controllers

import (
	"e-commerce/app/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	render2 "github.com/unrolled/render"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

func GetShoppingCartID(w http.ResponseWriter, r *http.Request) string {
	session, _ := store.Get(r, sessionShoppingCart)
	if session.Values["cart-id"] == nil {
		session.Values["cart-id"] = uuid.New().String()
		_ = session.Save(r, w)
	}

	return fmt.Sprintf("%v", session.Values["cart-id"])
}

func GetShoppingCart(db *gorm.DB, cartID string) (*models.Cart, error) {
	var cart models.Cart

	existCart, err := cart.GetCart(db, cartID)
	if err != nil {
		existCart, _ = cart.CreateCart(db, cartID)
	}

	_, _ = existCart.CalculateCart(db, cartID)

	updateCart, _ := cart.GetCart(db, cartID)

	totalWeight := 0
	productModel := models.Product{}
	for _, cartItem := range updateCart.CartItems {
		product, _ := productModel.FindByID(db, cartItem.ProductID)
		productWeight, _ := product.Weight.Float64()
		ceilWeight := math.Ceil(productWeight)
		itemWeight := cartItem.Qty * int(ceilWeight)
		totalWeight += itemWeight
	}

	updateCart.TotalWeight = totalWeight

	return updateCart, nil
}

func (server *Server) GetCart(w http.ResponseWriter, r *http.Request) {
	render := render2.New(render2.Options{
		Layout:     "layout",
		Extensions: []string{".html"},
	})

	var cart *models.Cart

	cartID := GetShoppingCartID(w, r)
	cart, _ = GetShoppingCart(server.DB, cartID)
	provinces, err := server.GetProvinces()
	if err != nil {
		log.Fatal(err)
	}

	_ = render.HTML(w, http.StatusOK, "cart", map[string]interface{}{
		"cart":      cart,
		"provinces": provinces,
	})
}

func (server *Server) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	productID := r.FormValue("product_id")
	qty, _ := strconv.Atoi(r.FormValue("qty"))

	productModel := models.Product{}
	product, err := productModel.FindByID(server.DB, productID)
	if err != nil {
		http.Redirect(w, r, "/products/"+product.Slug, http.StatusSeeOther)
	}

	if qty > product.Stock {
		http.Redirect(w, r, "/products/"+product.Slug, http.StatusSeeOther)
	}

	var cart *models.Cart

	cartID := GetShoppingCartID(w, r)
	cart, _ = GetShoppingCart(server.DB, cartID)
	_, err = cart.AddItem(server.DB, models.CartItem{
		ProductID: productID,
		Qty:       qty,
	})
	if err != nil {
		http.Redirect(w, r, "/products/"+product.Slug, http.StatusSeeOther)
	}

	http.Redirect(w, r, "/carts", http.StatusSeeOther)
}

func (server *Server) UpdateCart(w http.ResponseWriter, r *http.Request) {
	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(server.DB, cartID)

	for _, item := range cart.CartItems {
		qty, _ := strconv.Atoi(r.FormValue(item.ID))

		_, err := cart.UpdateItemQty(server.DB, item.ID, qty)
		if err != nil {
			http.Redirect(w, r, "/carts", http.StatusSeeOther)
		}
	}

	http.Redirect(w, r, "/carts", http.StatusSeeOther)
}

func (server *Server) RemoveItemByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		http.Redirect(w, r, "/carts", http.StatusSeeOther)
	}

	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(server.DB, cartID)

	err := cart.RemoveItemByID(server.DB, vars["id"])
	if err != nil {
		http.Redirect(w, r, "/carts", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/carts", http.StatusSeeOther)
}

func (server *Server) GetCitiesByProvince(w http.ResponseWriter, r *http.Request) {
	provinceID := r.URL.Query().Get("province_id")

	cities, err := server.GetCitiesByProvinceID(provinceID)
	if err != nil {
		log.Fatal(err)
	}

	res := Result{
		Code:    200,
		Data:    cities,
		Message: "Success",
	}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)
}

func (server *Server) CalculateShipping(w http.ResponseWriter, r *http.Request) {
	origin := os.Getenv("API_ONGKIR_ORIGIN")
	destination := r.FormValue("city_id")
	courier := r.FormValue("courier")

	if destination == "" {
		http.Error(w, "invalid destination", http.StatusInternalServerError)
	}

	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(server.DB, cartID)

	shippingFeeOption, err := server.CalculateShippingFee(models.ShippingFeeParams{
		Origin:      origin,
		Destination: destination,
		Weight:      cart.TotalWeight,
		Courier:     courier,
	})

	if err != nil {
		http.Error(w, "calculatin failed", http.StatusInternalServerError)
	}

	fmt.Println(shippingFeeOption, "TAHHHH")
}
