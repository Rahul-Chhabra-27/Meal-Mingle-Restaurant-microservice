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

func GoDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
func ConnectDB(dsn string) (*gorm.DB, *gorm.DB, error) {
	restaurantDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	restaurantDB.AutoMigrate(&model.Restaurant{})
	restaurantitemDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	restaurantitemDB.AutoMigrate(&model.RestaurantItem{})
	return restaurantDB, restaurantitemDB, nil
}

func ValidateRestaurantFields(RestaurantName string, RestaurantCity string, RestaurantAddress string, RestaurantPhone string, RestaurantAvailability string, RestaurantImageUrl string) bool {

	if RestaurantImageUrl == "" ||  RestaurantName == "" || RestaurantCity == "" || RestaurantAddress == "" || RestaurantPhone == "" || RestaurantAvailability == "" {
		return false
	}
	return true
}

func ValidateRestaurantItemFields(RestaurantItemName string, RestaurantItemImageUrl string) bool {
	fmt.Println(RestaurantItemName, RestaurantItemImageUrl)
	if RestaurantItemName == "" || RestaurantItemImageUrl == "" {
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