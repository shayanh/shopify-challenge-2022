package models

import "gorm.io/gorm"

// Inventory denotes an inventory.
type Inventory struct {
	gorm.Model
	Name string `gorm:"not null;unique"`
}

type InventoryRepository struct {
	DB *gorm.DB
}

func (rep *InventoryRepository) Create(inventory Inventory) (Inventory, error) {
	err := rep.DB.Create(&inventory).Error
	return inventory, err
}

func (rep *InventoryRepository) FirstOrCreate(inventory Inventory) (Inventory, error) {
	var res Inventory
	err := rep.DB.FirstOrCreate(&res, inventory).Error
	return res, err
}

func (rep *InventoryRepository) FindByID(id uint) (Inventory, error) {
	var inventory Inventory
	err := rep.DB.First(&inventory, id).Error
	return inventory, err
}

func (rep *InventoryRepository) FindAll() ([]Inventory, error) {
	var inventories []Inventory
	err := rep.DB.Find(&inventories).Error
	return inventories, err
}

// Item is an inventory item.
type Item struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `sql:"type:text"`
	Quantity    int    `gorm:"default:0"`
	InventoryID uint   `gorm:"not null"`
	Inventory   Inventory
}

type ItemRepository struct {
	DB *gorm.DB
}

func (rep *ItemRepository) Create(item Item) (Item, error) {
	err := rep.DB.Create(&item).Error
	return item, err
}

func (rep *ItemRepository) FirstOrCreate(item Item) (Item, error) {
	var res Item
	err := rep.DB.FirstOrCreate(&res, item).Error
	return res, err
}

func (rep *ItemRepository) Update(item Item) (Item, error) {
	err := rep.DB.Save(&item).Error
	return item, err
}

func (rep *ItemRepository) DeleteByID(id uint) error {
	return rep.DB.Delete(&Item{}, id).Error
}

func (rep *ItemRepository) FindByID(id uint) (Item, error) {
	var item Item
	err := rep.DB.First(&item, id).Error
	return item, err
}

func (rep *ItemRepository) FindAll() ([]Item, error) {
	var items []Item
	err := rep.DB.Find(&items).Error
	return items, err
}
