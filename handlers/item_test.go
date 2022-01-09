package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/shayanh/shopify-challenge-2022/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func fillInitialData(db *gorm.DB) ([]models.Inventory, []models.Item, error) {
	var invs = []models.Inventory{
		{Name: "School"},
		{Name: "Software"},
		{Name: "Phones"},
	}
	invRepo := models.InventoryRepository{
		DB: db,
	}
	for i := range invs {
		updatedInv, err := invRepo.FirstOrCreate(invs[i])
		if err != nil {
			return nil, nil, err
		}
		invs[i] = updatedInv
	}

	var items = []models.Item{
		{Name: "Pencil", InventoryID: invs[0].ID, Quantity: 8,
			Description: "Black writing pencil for school days."},
		{Name: "Backpack", InventoryID: invs[0].ID, Quantity: 11,
			Description: "Medium sized school backpack."},
		{Name: "Anti Virus", InventoryID: invs[1].ID, Quantity: 3,
			Description: "Strong protection for your machine."},
		{Name: "iPhone 13", InventoryID: invs[2].ID, Quantity: 9,
			Description: "Smartphone by Apple company."},
	}
	itemRepo := models.ItemRepository{
		DB: db,
	}
	for i := range items {
		updatedItem, err := itemRepo.FirstOrCreate(items[i])
		if err != nil {
			return nil, nil, err
		}
		items[i] = updatedItem
	}
	return invs, items, nil
}

type mockedRenderer struct {
	mock.Mock
}

func (m *mockedRenderer) Render(w http.ResponseWriter, name string, data interface{}) {
	m.Called(w, name, data)
}

type ItemHandlerTestSuite struct {
	suite.Suite

	db       *gorm.DB
	invRepo  *models.InventoryRepository
	itemRepo *models.ItemRepository

	h        *ItemHandler
	renderer *mockedRenderer

	initInvs  []models.Inventory
	initItems []models.Item
}

func (s *ItemHandlerTestSuite) SetupTest() {
	var err error
	s.db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = models.Migrate(s.db)
	if err != nil {
		log.Fatal(err)
	}

	s.initInvs, s.initItems, err = fillInitialData(s.db)
	if err != nil {
		log.Fatal(err)
	}

	s.invRepo = &models.InventoryRepository{DB: s.db}
	s.itemRepo = &models.ItemRepository{DB: s.db}
	s.renderer = &mockedRenderer{}
	s.renderer.On("Render", mock.Anything, mock.Anything, mock.Anything)
	s.h = NewItemHandler(s.itemRepo, s.invRepo, s.renderer)
}

func (s *ItemHandlerTestSuite) TearDownTest() {
	if err := os.Remove("test.db"); err != nil {
		log.Fatal(err)
	}
}

func (s *ItemHandlerTestSuite) TestListItems() {
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	w := httptest.NewRecorder()

	s.h.ListItems(w, req)

	resp := w.Result()
	s.Equal(resp.StatusCode, http.StatusOK, "status must be ok")
	s.renderer.AssertNumberOfCalls(s.T(), "Render", 1)
	calledListItemsPage := s.renderer.Calls[0].Arguments[2].(listItemsPage)
	s.Require().Equal(len(calledListItemsPage.Items), len(s.initItems))
	for i := range s.initItems {
		s.Equal(s.initItems[i].Name, calledListItemsPage.Items[i].Name)
		s.Equal(s.initItems[i].InventoryID, calledListItemsPage.Items[i].InventoryID)
		s.NotNil(calledListItemsPage.Items[i].Inventory)
	}
}

func makeItemPostForm(item models.Item) url.Values {
	form := url.Values{}
	form.Add("itemName", item.Name)
	form.Add("itemDescription", item.Description)
	form.Add("itemQuantity", strconv.Itoa(item.Quantity))
	form.Add("itemInventory", strconv.Itoa(int(item.InventoryID)))
	return form
}

func (s *ItemHandlerTestSuite) TestPostCreateItem_Successful() {
	item := models.Item{
		Name:        "test",
		Description: "test test",
		Quantity:    10,
		InventoryID: s.initInvs[2].ID,
	}

	data := makeItemPostForm(item)
	req := httptest.NewRequest(http.MethodPost, "/items/create", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	s.h.PostCreateItem(w, req)

	resp := w.Result()
	s.Equal(resp.StatusCode, http.StatusFound)
	s.renderer.AssertNumberOfCalls(s.T(), "Render", 0)

	items, err := s.itemRepo.FindAll()
	s.Require().Nil(err)
	s.Equal(len(items), len(s.initItems)+1)
}

func (s *ItemHandlerTestSuite) TestPostCreateItem_NoName() {
	item := models.Item{
		Description: "test test",
		Quantity:    10,
		InventoryID: s.initInvs[2].ID,
	}

	data := makeItemPostForm(item)
	req := httptest.NewRequest(http.MethodPost, "/items/create", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	s.h.PostCreateItem(w, req)

	resp := w.Result()
	s.Equal(resp.StatusCode, http.StatusOK, "status must be ok")
	s.renderer.AssertNumberOfCalls(s.T(), "Render", 1)
	calledEditItemPage := s.renderer.Calls[0].Arguments[2].(editItemPage)
	s.NotNil(calledEditItemPage.Error)
}

func (s *ItemHandlerTestSuite) TestPostEditItem_Successful() {
	item := s.initItems[1]
	item.Quantity += 1

	data := makeItemPostForm(item)
	target := fmt.Sprintf("/items/%d/edit", item.ID)
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(item.ID))})
	w := httptest.NewRecorder()

	s.h.PostEditItem(w, req)

	resp := w.Result()
	s.Equal(resp.StatusCode, http.StatusFound)
	s.renderer.AssertNumberOfCalls(s.T(), "Render", 0)

	retItem, err := s.itemRepo.FindByID(item.ID)
	s.Nil(err)
	s.Equal(item.Quantity, retItem.Quantity)
}

func (s *ItemHandlerTestSuite) TestPostEditItem_NotFound() {
	id := 100
	target := fmt.Sprintf("/items/%d/edit", id)
	req := httptest.NewRequest(http.MethodPost, target, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(id)})
	w := httptest.NewRecorder()

	s.h.PostEditItem(w, req)

	resp := w.Result()
	s.Equal(resp.StatusCode, http.StatusNotFound)
}

func TestItemHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ItemHandlerTestSuite))
}
