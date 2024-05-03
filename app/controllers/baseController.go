package controllers

import (
	"e-commerce/app/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
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

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// nanti coba kembangkan agar hash key mengambil dari file .env
var store = sessions.NewCookieStore([]byte("apes-session-key"))
var sessionShoppingCart = "shopping-cart-session"
var sessionFlash = "flash-session"

func (server *Server) Initialize() {
	server.initializeDB()
	server.initializeRoutes()
}

func (server *Server) initializeDB() {
	var err error

	server.DB, err = gorm.Open(postgres.Open(os.Getenv("DB_SOURCE")))
	//server.DB, err = gorm.Open(postgres.Open(os.Getenv("DB_SOURCE2")), &gorm.Config{})
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

func (server *Server) GetProvinces() ([]models.Province, error) {
	response, err := http.Get(os.Getenv("API_ONGKIR_BASE_URL") + "province?key=" + os.Getenv("API_ONGKIR_KEY"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	provinceResponse := models.ProvinceResponse{}

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, readErr
	}

	jsonErr := json.Unmarshal(body, &provinceResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return provinceResponse.ProvinceData.Results, nil
}

func (server *Server) GetCitiesByProvinceID(provinceID string) ([]models.City, error) {
	response, err := http.Get(os.Getenv("API_ONGKIR_BASE_URL") + "city?key=" + os.Getenv("API_ONGKIR_KEY") + "&province=" + provinceID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	cityResponse := models.CityResponse{}

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, readErr
	}

	jsonErr := json.Unmarshal(body, &cityResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return cityResponse.CityData.Results, nil
}

func (server *Server) CalculateShippingFee(shippingParams models.ShippingFeeParams) ([]models.ShippingFeeOption, error) {
	if shippingParams.Origin == "" || shippingParams.Destination == "" || shippingParams.Weight <= 0 || shippingParams.Courier == "" {
		return nil, errors.New("invalid params")
	}

	params := url.Values{}
	params.Add("key", os.Getenv("API_ONGKIR_KEY"))
	params.Add("origin", shippingParams.Origin)
	params.Add("destination", shippingParams.Destination)
	params.Add("weight", strconv.Itoa(shippingParams.Weight))
	params.Add("courier", shippingParams.Courier)

	response, err := http.PostForm(os.Getenv("API_ONGKIR_BASE_URL")+"cost", params)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	ongkirResponse := models.OngkirResponse{}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(body, &ongkirResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	var shippingFeeOptions []models.ShippingFeeOption
	for _, result := range ongkirResponse.OngkirData.Results {
		for _, cost := range result.Costs {
			shippingFeeOptions = append(shippingFeeOptions, models.ShippingFeeOption{
				Service: cost.Service,
				Fee:     cost.Cost[0].Value,
			})
		}
	}

	return shippingFeeOptions, nil
}

func SetFlash(w http.ResponseWriter, r *http.Request, name, value string) {
	session, err := store.Get(r, sessionFlash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.AddFlash(value, name)
	_ = session.Save(r, w)
}

func GetFlash(w http.ResponseWriter, r *http.Request, name string) []string {
	session, err := store.Get(r, sessionFlash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	fm := session.Flashes(name)
	if len(fm) < 0 {
		return nil
	}

	_ = session.Save(r, w)

	var flashes []string
	for _, fl := range fm {
		flashes = append(flashes, fl.(string))
	}

	return flashes
}
