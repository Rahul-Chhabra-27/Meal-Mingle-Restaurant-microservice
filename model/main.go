package model

import "gorm.io/gorm"

type Restaurant struct {
	gorm.Model 
	Name string `gorm:"unique"`
	Address string 	
	Phone string `gorm:"unique"`
	Availability string `gorm:"default:open"`
	Rating float32 `gorm:"default:0"`
}
type RestaurantItem struct {
	gorm.Model
	ItemName string
	ItemPrice float32
	ImageUrl string
	RestaurantId uint `gorm:"foreignKey:ID"`
}
