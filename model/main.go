package model

import "gorm.io/gorm"

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type Restaurant struct {
	gorm.Model
	Name                         string  `gorm:"unique"`
	Phone                        string  `gorm:"unique"`
	Availability                 string  `gorm:"default:open"`
	Rating                       float32 `gorm:"default:0"`
	RestaurantOwnerMail          string
	ImageUrl                     string
	OperationDays                string
	OperationHours               string
	RestaurantMinimumOrderAmount int64
	RestaurantDiscountPercentage float32
}
type Address struct {
	gorm.Model
	RestaurantId uint `gorm:"foreignKey:RestaurantId;references:ID;uniqueIndex:idx_restaurant_address"` // foreign key referencing the primary key of the Restaurant table
	StreetName   string
	Pincode      string
	City         string
	Country      string
}
type RestaurantItem struct {
	gorm.Model
	ItemName     string `gorm:"type:varchar(255);uniqueIndex:idx_restaurant_items"`
	ItemPrice    int64
	ImageUrl     string
	RestaurantId uint `gorm:"foreignKey:RestaurantId;references:ID;uniqueIndex:idx_restaurant_items"` // foreign key referencing the primary key of the Restaurant table
	Category     string
	CuisineType  string
	Veg          bool
}
