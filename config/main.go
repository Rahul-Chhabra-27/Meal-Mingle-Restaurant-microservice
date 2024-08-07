package config

import (
	"fmt"
	"log"
	"os"
	"restaurant-micro/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DatabaseDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
}
func GoDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
func ConnectDB() (*gorm.DB, *gorm.DB, *gorm.DB, error) {
	restaurantDB, err := gorm.Open(mysql.Open(DatabaseDsn()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	restaurantDB.AutoMigrate(&model.Restaurant{})
	restaurantitemDB, err := gorm.Open(mysql.Open(DatabaseDsn()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	restaurantitemDB.AutoMigrate(&model.RestaurantItem{})

	restaurantAddress, err := gorm.Open(mysql.Open(DatabaseDsn()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	restaurantAddress.AutoMigrate(&model.Address{})

	return restaurantDB, restaurantitemDB, restaurantAddress, nil
}

func ValidateRestaurantFields(RestaurantName string, RestaurantAddress model.Address,
	RestaurantPhone string, RestaurantAvailability string,
	RestaurantImageUrl string, RestaurantOperationHours string,
	RestaurantOperationDays string,
	RestauarntRating float32, restaurantMinimumOrderAmount int64,
	restaurantDiscountPercentage float32,

) bool {
	if RestauarntRating <= 0 || RestaurantImageUrl == "" || RestaurantName == "" || RestaurantAddress.City == "" ||
		RestaurantAddress.Country == "" || RestaurantAddress.Pincode == "" ||
		RestaurantAddress.StreetName == "" || RestaurantPhone == "" ||
		RestaurantAvailability == "" || RestaurantOperationHours == "" ||
		RestaurantOperationDays == "" || restaurantMinimumOrderAmount < 0 ||
		restaurantDiscountPercentage < 0 {
		return false
	}
	return true
}

func ValidateRestaurantItemFields(RestaurantItemName string, RestaurantItemImageUrl string,
	RestaurantItemPrice int64, RestaurantItemCategory string,
	RestaurantItemCuisine string, RestaurantId string) bool {
	if RestaurantItemName == "" || RestaurantItemImageUrl == "" || RestaurantItemPrice <= 0 ||
		RestaurantItemCategory == "" || RestaurantItemCuisine == "" || RestaurantId == "" {
		return false
	}
	return true
}

func ValidateRestaurantPhone(restaurantPhone string) bool {
	if len(restaurantPhone) != 10 {
		return false
	}

	for _, char := range restaurantPhone {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func ValidateAddressFields(City string, Country string, Pincode string, StreetName string) bool {
	if City == "" || Country == "" || Pincode == "" || StreetName == "" {
		return false
	}
	return true
}
