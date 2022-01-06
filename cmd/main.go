package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shayanh/shopify-challenge-2022/handlers"
	"github.com/shayanh/shopify-challenge-2022/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = models.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	err = fillInitialData(db)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.StrictSlash(true)

	itemRepo := &models.ItemRepository{
		DB: db,
	}
	itemHandler := handlers.NewItemHandler(itemRepo)
	itemHandler.Handle(router.PathPrefix("/items").Subrouter())

	err = http.ListenAndServe("localhost:8000", router)
	if err != nil {
		log.Fatal(err)
	}
}
