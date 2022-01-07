package models

import "gorm.io/gorm"

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

type Item struct {
	gorm.Model
	Name        string
	Description string `sql:"type:text"`
	InventoryID uint
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

func (rep *ItemRepository) FindAll() ([]*Item, error) {
	var items []*Item
	err := rep.DB.Preload("Inventory").Find(&items).Error
	return items, err
}
