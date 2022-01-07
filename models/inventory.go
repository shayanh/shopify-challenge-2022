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

func (rep *InventoryRepository) Create(inventory *Inventory) error {
	return rep.DB.Create(inventory).Error
}

func (rep *InventoryRepository) FirstOrCreate(inventory *Inventory) error {
	return rep.DB.FirstOrCreate(inventory, *inventory).Error
}

func (rep *InventoryRepository) FindByID(id uint) (*Inventory, error) {
	var inventory *Inventory
	err := rep.DB.First(&inventory, id).Error
	return inventory, err
}

func (rep *InventoryRepository) FindAll() ([]*Inventory, error) {
	var inventories []*Inventory
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

func (rep *ItemRepository) Create(item *Item) error {
	return rep.DB.Create(item).Error
}

func (rep *ItemRepository) FirstOrCreate(item *Item) error {
	return rep.DB.FirstOrCreate(item, *item).Error
}

func (rep *ItemRepository) Update(item *Item) error {
	return rep.DB.Model(item).Updates(item).Error
}

func (rep *ItemRepository) DeleteByID(id uint) error {
	return rep.DB.Delete(&Item{}, id).Error
}

func (rep *ItemRepository) FindByID(id uint) (*Item, error) {
	var item *Item
	err := rep.DB.First(&item, id).Error
	return item, err
}

func (rep *ItemRepository) FindAll() ([]*Item, error) {
	var items []*Item
	err := rep.DB.Preload("Inventory").Find(&items).Error
	return items, err
}
