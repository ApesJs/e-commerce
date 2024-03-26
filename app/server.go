package app

import (
	"e-commerce/app/controllers"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Run() {
	var server = controllers.Server{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("cannot load env:", err)
		return
	}

	server.Initialize()
	server.Run(os.Getenv("APP_PORT"))
}
