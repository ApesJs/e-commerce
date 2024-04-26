package controllers

import (
	"e-commerce/app/models"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type PageLink struct {
	Page          int32
	Url           string
	IsCurrentPage bool
}

type PaginationLinks struct {
	CurrentPage string
	NextPage    string
	PrevPage    string
	TotalRows   int32
	TotalPages  int32
	Links       []PageLink
}

type PaginationParams struct {
	Path        string
	TotalRows   int32
	PerPage     int32
	CurrentPage int32
}

// nanti coba kembangkan agar hash key mengambil dari file .env
var store = sessions.NewCookieStore([]byte("apes-session-key"))
var sessionShoppingCart = "shopping-cart-session"

func (server *Server) Initialize() {
	server.initializeDB()
	server.initializeRoutes()
}

func (server *Server) initializeDB() {
	var err error

	server.DB, err = gorm.Open(postgres.Open(os.Getenv("DB_SOURCE")))
	if err != nil {
		log.Fatal("cannot connect to database:", err)
		return
	}

	allModels := []interface{}{
		&models.Address{},
		&models.Cart{},
		&models.CartItem{},
		&models.Category{},
		&models.OrderCustomer{},
		&models.OrderItem{},
		&models.Order{},
		&models.Payment{},
		&models.ProductImage{},
		&models.Product{},
		&models.Section{},
		&models.Shipment{},
		&models.User{},
	}

	err = server.DB.AutoMigrate(allModels...)
	if err != nil {
		log.Fatal("cannot migrate database models:", err)
		return
	}

	//for i := 1; i <= 10; i++ {
	//fUser := fakers.UserFaker(server.DB)
	//err = server.DB.Create(&fUser).Error
	//if err != nil {
	//	fmt.Println("Error inserting user:", err)
	//	return
	//}

	//	fProduct := fakers.ProductFaker(server.DB)
	//	err = server.DB.Create(&fProduct).Error
	//	if err != nil {
	//		fmt.Println("Error inserting user:", err)
	//		return
	//	}
	//}

}

func (server *Server) Run(addr string) {
	err := http.ListenAndServe(addr, server.Router)
	if err != nil {
		log.Fatal("cannot run server:", err)
		return
	}
}

func GetPaginationLinks(params PaginationParams) (PaginationLinks, error) {
	var links []PageLink
	totalPages := int32(math.Ceil(float64(params.TotalRows) / float64(params.PerPage)))

	for i := 1; int32(i) <= totalPages; i++ {
		links = append(links, PageLink{
			Page:          int32(i),
			Url:           fmt.Sprintf("%s/%s?page=%s", os.Getenv("APP_URL"), params.Path, fmt.Sprintf(strconv.Itoa(i))),
			IsCurrentPage: int32(i) == params.CurrentPage,
		})
	}

	var nextPage int32
	var prevPage int32

	prevPage = 1
	nextPage = totalPages

	if params.CurrentPage > 2 {
		prevPage = params.CurrentPage - 1
	}

	if params.CurrentPage < totalPages {
		nextPage = params.CurrentPage + 1
	}

	return PaginationLinks{
		CurrentPage: fmt.Sprintf("%s/%s?page=%s", os.Getenv("APP_URL"), params.Path, fmt.Sprintf(string(params.CurrentPage))),
		NextPage:    fmt.Sprintf("%s/%s?page=%s", os.Getenv("APP_URL"), params.Path, fmt.Sprintf(string(nextPage))),
		PrevPage:    fmt.Sprintf("%s/%s?page=%s", os.Getenv("APP_URL"), params.Path, fmt.Sprintf(string(prevPage))),
		TotalRows:   params.TotalRows,
		TotalPages:  totalPages,
		Links:       links,
	}, nil
}
