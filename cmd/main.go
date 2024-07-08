package main

import (
	"github.com/joho/godotenv"
	"log"
	_ "tooki/docs"
	"tooki/pkg/handler"
	"tooki/pkg/repository"
)

// @title						Tooki??
// @version					1.0
// @license.name				Apache 2.0
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := repository.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	h := handler.NewHandler(db)
	router := h.InitRoutes()

	err = router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
