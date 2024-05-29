package config

import (
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
		return nil, nil, err
	}
	restaurantitemDB.AutoMigrate(&model.RestaurantItem{})
	return restaurantDB, restaurantitemDB, nil
}
