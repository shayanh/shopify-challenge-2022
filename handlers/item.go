package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/shayanh/shopify-challenge-2022/models"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

// ItemHandler implements web handlers related to Item entity. It uses repository
// objects to fetch data from the data store.
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

func getFormItem(r *http.Request) (models.Item, error) {
	var item models.Item
	item.Name = r.FormValue("itemName")
	item.Description = r.FormValue("itemDescription")

	var errs, err error
	item.Quantity, err = strconv.Atoi(r.FormValue("itemQuantity"))
	if err != nil || item.Quantity < 0 {
		errs = errors.New("invalid quantity")
	}

	invID, err := strconv.Atoi(r.FormValue("itemInventory"))
	if err != nil || invID < 0 {
		errs = multierr.Append(errs, errors.New("invalid inventory"))
	} else {
		item.InventoryID = uint(invID)
	}
	return item, errs
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
	Title       string
	FormAction  string
	Inventories []*models.Inventory
	Item        models.Item
	Error       error
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

func (h *ItemHandler) postCreateItem(w http.ResponseWriter, r *http.Request) {
	item, err := getFormItem(r)
	page := editItemPage{
		Title:      "Create Item",
		FormAction: "/items/create",
		Item:       item,
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
	h.renderEditPage(w, editItemPage{
		Title:      "Create Item",
		FormAction: "/items/create",
	})
}

func getParamItemID(r *http.Request) (uint, error) {
	params := mux.Vars(r)
	strID, ok := params["id"]
	if !ok {
		return 0, errors.New("missing item id")
	}
	intID, err := strconv.Atoi(strID)
	if err != nil || intID < 0 {
		return 0, errors.New("invalid item id")
	}
	return uint(intID), nil
}

func (h *ItemHandler) deleteItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := getParamItemID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.itemRepo.FindByID(itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "item not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = h.itemRepo.DeleteByID(itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/items", http.StatusFound)
}

func (h *ItemHandler) postEditItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := getParamItemID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.itemRepo.FindByID(itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "item not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	item, err := getFormItem(r)
	page := editItemPage{
		Title:      "Edit Item",
		FormAction: fmt.Sprintf("/items/%d/edit", itemID),
		Item:       item,
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

	item.ID = itemID
	err = h.itemRepo.Update(&item)
	if err != nil {
		page.Error = err
		h.renderEditPage(w, page)
		return
	}
	http.Redirect(w, r, "/items", http.StatusFound)
}

func (h *ItemHandler) editItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := getParamItemID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.itemRepo.FindByID(itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "item not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.renderEditPage(w, editItemPage{
		Title:      "Edit Item",
		FormAction: fmt.Sprintf("/items/%d/edit", itemID),
		Item:       *item,
	})
}

func (h *ItemHandler) exportCSV(w http.ResponseWriter, r *http.Request) {
	items, err := h.itemRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	records := [][]string{
		{"id", "name", "inventory", "qty", "created_at", "updated_at", "description"},
	}
	for _, item := range items {
		record := []string{
			strconv.Itoa(int(item.ID)), item.Name, item.Inventory.Name,
			strconv.Itoa(item.Quantity), item.CreatedAt.Format(time.RFC3339),
			item.UpdatedAt.Format(time.RFC3339), item.Description,
		}
		records = append(records, record)
	}

	csvWriter := csv.NewWriter(w)
	for _, record := range records {
		if err := csvWriter.Write(record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Println(err)
	}
}

// Handle registers related handlers into a given Router.
func (h *ItemHandler) Handle(router *mux.Router) {
	router.HandleFunc("/items", h.listItems).Methods(http.MethodGet)
	router.HandleFunc("/items/create", h.createItem).Methods(http.MethodGet)
	router.HandleFunc("/items/create", h.postCreateItem).Methods(http.MethodPost)
	router.HandleFunc("/items/{id:[0-9]+}/delete", h.deleteItem).Methods(http.MethodPost)
	router.HandleFunc("/items/{id:[0-9]+}/edit", h.editItem).Methods(http.MethodGet)
	router.HandleFunc("/items/{id:[0-9]+}/edit", h.postEditItem).Methods(http.MethodPost)
	router.HandleFunc("/items/csv", h.exportCSV).Methods(http.MethodGet)
}
