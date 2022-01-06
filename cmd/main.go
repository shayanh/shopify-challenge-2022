package main

import (
	"log"

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
	log.Println(err)
}
