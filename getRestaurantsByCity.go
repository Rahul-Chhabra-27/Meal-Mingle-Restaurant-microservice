package main

// import (
// 	"context"
// 	"fmt"
// 	"restaurant-micro/model"
// 	restaurantpb "restaurant-micro/proto/restaurant"
// )

// func (*RestaurantService) GetRestaurantsByCity(ctx context.Context, request *restaurantpb.GetRestaurantsByCityRequest) (*restaurantpb.GetRestaurantsByCityResponse, error) {
// 	// get all the restaurants by city.
// 	var restaurants []model.Restaurant
// 	primaryKey := restaurantDBConnector.Where("city = ?", request.City).Find(&restaurants).Error
// 	if primaryKey != nil {
// 		fmt.Println("Failed to get restaurants by city")
// 		return &restaurantpb.GetRestaurantsByCityResponse{
// 			Data: &restaurantpb.RestaurantData{
// 				TotalRestaurants: 0,
// 				Restaurants:      nil,
// 			},
// 			Message:    "",
// 			StatusCode: int64(500),
// 			Error:      "Failed to get restaurants",
// 		}, nil
// 	}

// 	var response restaurantpb.GetRestaurantsByCityResponse
// 	response.Data = &restaurantpb.RestaurantData{}
// 	for _, restaurant := range restaurants {
// 		response.Data.Restaurants = append(response.Data.Restaurants, &restaurantpb.Restaurant{
// 			RestaurantName:         restaurant.Name,
// 			RestaurantCity:         restaurant.City,
// 			RestaurantAddress:      restaurant.Address,
// 			RestaurantPhone:        restaurant.Phone,
// 			RestaurantAvailability: restaurant.Availability,
// 			RestaurantRating:       restaurant.Rating,
// 			RestaurantOwnerMail:    restaurant.RestaurantOwnerMail,
// 			RestaurantImageUrl:    restaurant.ImageUrl,
// 		})
// 	}
// 	response.Data.TotalRestaurants = int64(len(restaurants))
// 	response.Message = "Restaurants fetched successfully"
// 	response.StatusCode = 200
// 	response.Error = ""
// 	return &response, nil
// }
