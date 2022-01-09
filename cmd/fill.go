package main

import (
	"log"

	"github.com/shayanh/shopify-challenge-2022/models"
	"gorm.io/gorm"
)

func fillInitialData(db *gorm.DB) error {
	invs := []models.Inventory{
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
			return err
		}
		invs[i] = updatedInv
		log.Printf("Created inventory - ID=%d, Name=%s\n", updatedInv.ID, updatedInv.Name)
	}

	items := []models.Item{
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
	for _, item := range items {
		updatedItem, err := itemRepo.FirstOrCreate(item)
		if err != nil {
			return err
		}
		log.Printf("Created item - ID=%d, Name=%s\n", updatedItem.ID, updatedItem.Name)
	}
	return nil
}
