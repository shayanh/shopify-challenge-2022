package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shayanh/shopify-challenge-2022/models"
)

type ItemHandler struct {
	itemRepo  *models.ItemRepository
	templates *template.Template
}

func NewItemHandler(itemRepo *models.ItemRepository) *ItemHandler {
	return &ItemHandler{
		itemRepo:  itemRepo,
		templates: template.Must(template.ParseFiles("./templates/list.html")),
	}
}

func (h *ItemHandler) listItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.itemRepo.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, item := range items {
		log.Println(item.Name, item.InventoryID, item.Inventory)
	}

	page := struct {
		Items []*models.Item
	}{
		Items: items,
	}

	err = h.templates.ExecuteTemplate(w, "list.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ItemHandler) Handle(router *mux.Router) {
	router.HandleFunc("", h.listItems).Methods("GET")
}
