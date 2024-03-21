package app

import (
	"e-commerce/app/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

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

func Run() {
	var server = Server{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("cannot load env:", err)
		return
	}

	server.Initialize()
	server.Run(os.Getenv("APP_PORT"))
}
