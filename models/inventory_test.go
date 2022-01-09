package models

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = Migrate(db)
	if err != nil {
		return nil, err
	}
	return db, err
}

func tearDownDB() error {
	return os.Remove("test.db")
}

func TestInventoryRepository(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tearDownDB(); err != nil {
			log.Fatal(err)
		}
	}()

	inv := Inventory{Name: "test"}
	invRepo := &InventoryRepository{DB: db}
	createdInv, err := invRepo.Create(inv)
	assert.Nil(t, err)
	assert.NotEqual(t, createdInv.ID, 0)
	assert.Equal(t, createdInv.Name, inv.Name)

	foundInv, err := invRepo.FindByID(createdInv.ID)
	assert.Nil(t, err)
	assert.Equal(t, foundInv.ID, createdInv.ID)
	assert.Equal(t, foundInv.Name, createdInv.Name)
}

func TestItemRepository(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tearDownDB(); err != nil {
			log.Fatal(err)
		}
	}()

	invRepo := &InventoryRepository{DB: db}
	inv, err := invRepo.Create(Inventory{Name: "test"})
	assert.Nil(t, err)

	itemRepo := &ItemRepository{DB: db}
	item := Item{Name: "t1", InventoryID: inv.ID, Quantity: 8,
		Description: "d1"}

	createdItem, err := itemRepo.Create(item)
	assert.Nil(t, err)
	assert.NotEqual(t, createdItem.ID, 0)
	assert.Equal(t, createdItem.Name, item.Name)

	foundItem, err := itemRepo.FindByID(createdItem.ID)
	assert.Nil(t, err)
	assert.Equal(t, foundItem.ID, createdItem.ID)
	assert.Equal(t, foundItem.Name, createdItem.Name)

	item.Quantity += 1
	updatedItem, err := itemRepo.Update(item)
	assert.Equal(t, updatedItem.Quantity, item.Quantity)
}
