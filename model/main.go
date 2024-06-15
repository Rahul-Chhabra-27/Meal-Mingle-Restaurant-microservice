package model

import "gorm.io/gorm"

type Restaurant struct {
	gorm.Model 
	Name string `gorm:"unique"`
	Address string 	
	Phone string `gorm:"unique"`
	Availability string `gorm:"default:open"`
	Rating float32 `gorm:"default:0"`
	RestaurantOwnerMail string
	City string
	ImageUrl string
}
type RestaurantItem struct {
    gorm.Model
    ItemName string `gorm:"type:varchar(255);uniqueIndex:idx_restaurant_items"`
    ItemPrice int64
    ImageUrl string
	// restuarant-item category...
    RestaurantId uint  `gorm:"foreignKey:RestaurantId;references:ID;uniqueIndex:idx_restaurant_items"` // foreign key referencing the primary key of the Restaurant table
}
