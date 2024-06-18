package main

import (
	"context"
	"fmt"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
)

func (*RestaurantService) GetAllRestaurants(context.Context, *restaurantpb.GetAllRestaurantsRequest) (*restaurantpb.GetAllRestaurantsResponse, error) {
	var restaurants []model.Restaurant
	err := restaurantDBConnector.Find(&restaurants)
	if err.Error != nil {
		fmt.Println("[ GetAllRestaurants ] Failed to get restaurants from the database", err.Error)
		return &restaurantpb.GetAllRestaurantsResponse{
			Data: &restaurantpb.GetAllRestaurantsData{
				TotalRestaurants: 0,
				Restaurants:      nil,
			},
			Message:    "",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}

	restaurantsResponse := []*restaurantpb.Restaurant{}
	totalRestaurants := 0
	for _, restaurant := range restaurants {
		totalRestaurants++
		// fetch all address of the restaurant from the database.
		var restaurantAddress model.Address
		restaurantAddressErr := restaurantAddressDBConnector.Where("restaurant_id = ?", restaurant.ID).First(&restaurantAddress).Error
		if restaurantAddressErr != nil {
			fmt.Println("[ GetAllRestaurants ] Failed to get restaurant address from the database", restaurantAddressErr)
			return &restaurantpb.GetAllRestaurantsResponse{
				Data: &restaurantpb.GetAllRestaurantsData{
					TotalRestaurants: 0,
					Restaurants:      nil,
				},
				Message:    "",
				StatusCode: 500,
				Error:      "Internal Server Error",
			}, nil
		}
		restaurantsResponse = append(restaurantsResponse, &restaurantpb.Restaurant{
			RestaurantName:           restaurant.Name,
			RestaurantAvailability:   restaurant.Availability,
			RestaurantRating:         restaurant.Rating,
			RestaurantImageUrl:       restaurant.ImageUrl,
			RestaurantPhoneNumber:    restaurant.Phone,
			RestaurantOperationDays:  restaurant.OperationDays,
			RestaurantOperationHours: restaurant.OperationHours,
			RestaurantOwnerMail:      restaurant.RestaurantOwnerMail,
			RestaurantAddress: &restaurantpb.Address{
				City:       restaurantAddress.City,
				Country:    restaurantAddress.Country,
				Pincode:    restaurantAddress.Pincode,
				StreetName: restaurantAddress.StreetName,
			},
		})
	}
	return &restaurantpb.GetAllRestaurantsResponse{
		Data: &restaurantpb.GetAllRestaurantsData{
			TotalRestaurants: int64(totalRestaurants),
			Restaurants:      restaurantsResponse,
		},
		Message:    "Restaurants fetched successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
