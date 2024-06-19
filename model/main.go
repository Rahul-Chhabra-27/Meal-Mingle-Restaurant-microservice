package model

import "gorm.io/gorm"

type Restaurant struct {
	gorm.Model 
	Name string `gorm:"unique"`
	Phone string `gorm:"unique"`
	Availability string `gorm:"default:open"`
	Rating float32 `gorm:"default:0"`
	RestaurantOwnerMail string
	ImageUrl string
	OperationDays string
	OperationHours string
}
type Address struct {
	gorm.Model
	RestaurantId uint `gorm:"foreignKey:RestaurantId;references:ID;uniqueIndex:idx_restaurant_address"` // foreign key referencing the primary key of the Restaurant table
	StreetName string
	Pincode string
	City string
	Country string
}
// {
// 	"restaurantId": "1",
// 	"restaurantName": "The Pizza Kings",
// 	"restaurantAddress": {
// 		"streetNumber": "12/92",
// 		"streetName": "Babar Pur, Shahdara, Geeta Colony",
// 		"city": "Delhi",
// 		"country": "India"
// 	},
// 	"restaurantRating": 3.6,
// 	"restaurantDiscount": 12,
// 	"restaurantAvailability": true,
// 	"restaurantImage": "restaurant_1.png",
// 	"restaurantOperationDays": "Mon-Fri",
// 	"restaurantOperationHours": "10:00AM-11:00PM",
// 	"restaurantItems": [
// 		{
// 			"restaurantItemId": "1_1",
// 			"restaurantItemName": "Margherita Pizza",
// 			"restaurantItemPrice": 390,
// 			"restaurantItemCategory": "Pizza",
// 			"restaurantItemImageUrl": "restaurant_1_menu_1.png",
// 			"restaurantItemCuisineType": "North Indian",
// 			"restaurantItemVeg": true
// 		}
// 	]
// },
type RestaurantItem struct {
    gorm.Model
    ItemName string `gorm:"type:varchar(255);uniqueIndex:idx_restaurant_items"`
    ItemPrice int64
    ImageUrl string
    RestaurantId uint  `gorm:"foreignKey:RestaurantId;references:ID;uniqueIndex:idx_restaurant_items"` // foreign key referencing the primary key of the Restaurant table
	Category string
	CuisineType string
	Veg bool
}
