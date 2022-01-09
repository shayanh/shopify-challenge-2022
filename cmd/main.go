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

type ResponseWriterWrapper struct {
	Status int
	http.ResponseWriter
}

func (rww *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rww.Status = statusCode
	rww.ResponseWriter.WriteHeader(statusCode)
}

func NewResponseWriterWrapper(rww http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{http.StatusOK, rww}
}

func logDecorator(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := NewResponseWriterWrapper(w)
		h.ServeHTTP(wr, r)
		log.Printf("%s %s %d", r.Method, r.URL.Path, wr.Status)
	})
}

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
	invRepo := &models.InventoryRepository{
		DB: db,
	}
	renderer := handlers.NewHTMLRenderer("./templates")

	itemHandler := handlers.NewItemHandler(itemRepo, invRepo, renderer)
	itemHandler.HandleFuncs(router)

	err = http.ListenAndServe("localhost:8000", logDecorator(router))
	if err != nil {
		log.Fatal(err)
	}
}
