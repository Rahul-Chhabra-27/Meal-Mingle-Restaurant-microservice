package main

import (
	"context"
	"fmt"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"

	"google.golang.org/grpc/codes"
)

func (*RestaurantService) AddRestaurant(ctx context.Context, response *restaurantpb.AddRestaurantRequest) (*restaurantpb.AddRestaurantResponse, error) {
	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.AddRestaurantResponse{Message: "", Error: "Internal Server Error", StatusCode: int64(500)}, nil
	}
	restaurantName := response.Restaurant.RestaurantName;
	restaurantCity := response.Restaurant.RestaurantCity;
	restaurantAvailability := response.Restaurant.RestaurantAvailability;
	restaurantAddress := response.Restaurant.RestaurantAddress
	restaurantPhone := response.Restaurant.RestaurantPhone
	restaurantRating := response.Restaurant.RestaurantRating
	restaurantImageUrl := response.Restaurant.RestaurantImageUrl

	if !config.ValidateRestaurantFields(restaurantName, restaurantCity,restaurantAddress, restaurantPhone, restaurantAvailability,restaurantImageUrl) {

		return &restaurantpb.AddRestaurantResponse{
			Message:    "",
			StatusCode: 400,
			Error:      "Invalid restaurant fields",
		}, nil
	}
	if !config.ValidateRestaurantPhone(restaurantPhone) {
		return &restaurantpb.AddRestaurantResponse{
			Message:    "",
			StatusCode: 400,
			Error:      "Invalid phone number",
		}, nil
	}
	var restaurant model.Restaurant
	restaurant.Name = restaurantName
	restaurant.City = restaurantCity
	restaurant.Address = restaurantAddress
	restaurant.Phone = restaurantPhone
	restaurant.Availability = restaurantAvailability
	restaurant.Rating = restaurantRating
	restaurant.RestaurantOwnerMail = userEmail
	restaurant.ImageUrl = restaurantImageUrl
	restaurantNotFoundErr := restaurantDBConnector.Where("name = ?", restaurant.Name).First(&restaurant).Error

	if restaurantNotFoundErr == nil {
		return &restaurantpb.AddRestaurantResponse{
			Message:    "Restaurant with the same name already exists",
			StatusCode: int64(409),
			Error:      "",
		}, nil
	}
	primaryKey := restaurantDBConnector.Create(&restaurant)
	if primaryKey.Error != nil {
		return &restaurantpb.AddRestaurantResponse{
			Message:    "Failed to add restaurant",
			StatusCode: 500,
			Error:      primaryKey.Error.Error(),
		}, nil
	}
	return &restaurantpb.AddRestaurantResponse{
		Data:&restaurantpb.RestaurantData {
			Restaurants: []*restaurantpb.Restaurant{
				{
					RestaurantName:restaurantName,
					RestaurantCity:restaurantCity,
					RestaurantAddress:restaurantAddress,
					RestaurantPhone:restaurantPhone,
					RestaurantAvailability:restaurantAvailability,
					RestaurantRating:restaurantRating,
					RestaurantOwnerMail:userEmail,
					RestaurantImageUrl:restaurantImageUrl,
				},
			},
			TotalRestaurants: 1,
		},
		Message:    "Restaurant added successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
