package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shayanh/shopify-challenge-2022/models"
)

type ItemHandler struct {
	itemRepo  *models.ItemRepository
	invRepo   *models.InventoryRepository
	templates *template.Template
}

func NewItemHandler(itemRepo *models.ItemRepository, invRepo *models.InventoryRepository) *ItemHandler {
	templatesDir := "./templates"
	templateNames := []string{
		"list.html",
		"edit.html",
	}
	var templateFileNames []string
	for _, tn := range templateNames {
		templateFileNames = append(templateFileNames, filepath.Join(templatesDir, tn))
	}

	return &ItemHandler{
		itemRepo:  itemRepo,
		invRepo:   invRepo,
		templates: template.Must(template.ParseFiles(templateFileNames...)),
	}
}

func (h *ItemHandler) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ItemHandler) listItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.itemRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := struct {
		Items []*models.Item
	}{
		Items: items,
	}
	h.renderTemplate(w, "list.html", page)
}

func getItem(r *http.Request) (models.Item, error) {
	var item models.Item
	item.Name = r.FormValue("itemName")
	item.Description = r.FormValue("itemDescription")
	invID, err := strconv.Atoi(r.FormValue("itemInventory"))
	if err != nil {
		return item, errors.New("invalid inventory")
	}
	item.InventoryID = uint(invID)
	return item, nil
}

func (h *ItemHandler) validateItem(item *models.Item) error {
	if item.Name == "" {
		return errors.New("item name cannot be empty")
	}
	_, err := h.invRepo.FindByID(item.InventoryID)
	if err != nil {
		return errors.New("invalid inventory")
	}
	return nil
}

type editItemPage struct {
	Inventories []*models.Inventory
	Error       error
	Name        string
	Description string
	InventoryID uint
}

func (h *ItemHandler) renderEditPage(w http.ResponseWriter, page editItemPage) {
	inventories, err := h.invRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page.Inventories = inventories
	h.renderTemplate(w, "edit.html", page)
}

func (h *ItemHandler) postItem(w http.ResponseWriter, r *http.Request) {
	item, err := getItem(r)
	page := editItemPage{
		Name:        item.Name,
		Description: item.Description,
		InventoryID: item.InventoryID,
	}
	if err != nil {
		page.Error = err
		h.renderEditPage(w, page)
		return
	}
	err = h.validateItem(&item)
	if err != nil {
		page.Error = err
		h.renderEditPage(w, page)
		return
	}

	err = h.itemRepo.Create(&item)
	if err != nil {
		page.Error = err
		h.renderEditPage(w, page)
		return
	}
	http.Redirect(w, r, "/items", http.StatusFound)
}

func (h *ItemHandler) createItem(w http.ResponseWriter, r *http.Request) {
	h.renderEditPage(w, editItemPage{})
}

func (h *ItemHandler) Handle(router *mux.Router) {
	router.HandleFunc("/items", h.listItems).Methods("GET")
	router.HandleFunc("/items", h.postItem).Methods("POST")
	router.HandleFunc("/items/create", h.createItem).Methods("GET")
}
