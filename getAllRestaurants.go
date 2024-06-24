package main

import (
	"context"
	"fmt"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"
)

func (*RestaurantService) GetAllRestaurants(ctx context.Context, request *restaurantpb.GetAllRestaurantsRequest) (*restaurantpb.GetAllRestaurantsResponse, error) {

	// get the user mail from the context
	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.GetAllRestaurantsResponse{
			Data:       &restaurantpb.GetAllRestaurantsResponseData{
				TotalRestaurants: 0,
				Restaurants:      nil,
			},
			Message:    "Failed to get user email from context",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}

	var restaurants []model.Restaurant
	err := restaurantDBConnector.Where("restaurant_owner_mail = ?", userEmail).Find(&restaurants).Error
	
	if err != nil {
		fmt.Println("[ GetAllRestaurants ] Failed to get restaurants from the database", err)
		return &restaurantpb.GetAllRestaurantsResponse{
			Data:       &restaurantpb.GetAllRestaurantsResponseData{
				TotalRestaurants: 0,
				Restaurants:      nil,
			},
			Message:    "Failed to get the restaurants to the database. Please try again later.",
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
				Data:       &restaurantpb.GetAllRestaurantsResponseData{
					TotalRestaurants: 0,
					Restaurants:      nil,
				},
				Message:    "Failed to get the restaurant addresses to the database. Please try again later.",
				StatusCode: 500,
				Error:      "Internal Server Error",
			}, nil
		}
		restaurantsResponse = append(restaurantsResponse, &restaurantpb.Restaurant{
			RestaurantId:             strconv.FormatUint(uint64(restaurant.ID), 10),
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
		Data: &restaurantpb.GetAllRestaurantsResponseData{
			TotalRestaurants: int64(totalRestaurants),
			Restaurants:      restaurantsResponse,
		},
		Message:    "Restaurants fetched successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
