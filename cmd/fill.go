package main

import (
	"github.com/shayanh/shopify-challenge-2022/models"
	"gorm.io/gorm"
)

func fillInitialData(db *gorm.DB) error {
	invs := []*models.Inventory{
		{Name: "My Inventory"},
		{Name: "Shopify"},
	}
	invRepo := models.InventoryRepository{
		DB: db,
	}
	for _, inv := range invs {
		err := invRepo.FirstOrCreate(inv)
		if err != nil {
			return err
		}
	}

	items := []*models.Item{
		{Name: "Pencil", InventoryID: invs[0].ID},
		{Name: "Bag", InventoryID: invs[0].ID},
		{Name: "Anti Virus", InventoryID: invs[1].ID},
	}
	itemRepo := models.ItemRepository{
		DB: db,
	}
	for _, item := range items {
		err := itemRepo.FirstOrCreate(item)
		if err != nil {
			return err
		}
	}
	return nil
}
