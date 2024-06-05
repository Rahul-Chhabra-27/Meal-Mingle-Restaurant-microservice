package main

import (
	"context"
	"fmt"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"

	"google.golang.org/grpc/codes"
)

func (*RestaurantService) AddRestaurantItem(ctx context.Context, response *restaurantpb.AddRestaurantItemRequest) (*restaurantpb.AddRestaurantItemResponse, error) {
	// Get the user email from the context
	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.AddRestaurantItemResponse{
			Message:    "",
			StatusCode: int64(codes.Internal),
			Error:      "Internal Server Error",
		}, nil
	}
	var restaurantItem model.RestaurantItem
	if !config.ValidateRestaurantItemFields(response.RestaurantItemName, response.RestaurantItemImageUrl) {
		return &restaurantpb.AddRestaurantItemResponse{
			Message:    "",
			StatusCode: int64(codes.InvalidArgument),
			Error:      "Invalid restaurant item fields",
		}, nil
	}
	restaurantItem.ItemName = response.RestaurantItemName
	restaurantItem.ItemPrice = response.RestaurantItemPrice
	restaurantItem.ImageUrl = response.RestaurantItemImageUrl

	if response.RestaurantName == "" {
		return &restaurantpb.AddRestaurantItemResponse{
			Message:    "",
			StatusCode: 400,
			Error:      "Restaurant name is required",
		}, nil
	}
	// fetch restaurant from restaurantDB
	var restaurant model.Restaurant
	primaryKey := restaurantDBConnector.Where("name = ?", response.RestaurantName).First(&restaurant)
	// check if the restaurant is exist or nor
	if primaryKey.Error != nil || restaurant.RestaurantOwnerMail != userEmail{
		return &restaurantpb.AddRestaurantItemResponse{
			Message:    "",
			StatusCode: int64(codes.NotFound),
			Error:      "Restaurant Does not exist OR you are not the owner of this restaurant",
		}, nil
	}
	restaurantItem.RestaurantId = restaurant.ID
	// Create a new restaurant item in the database and return the primary key if successful or an error if it fails
	result := restaurantItemDBConnector.Create(&restaurantItem)
	// Check if there is an error while creating the restaurant item
	if result.Error != nil {
		fmt.Println(result.Error)
		return &restaurantpb.AddRestaurantItemResponse{
			Message:    "",
			StatusCode: int64(codes.Internal),
			Error:      "Food Item is already exist",
		}, nil
	}
	// Return a success message if the restaurant item is created successfully
	return &restaurantpb.AddRestaurantItemResponse{
		Message:    "Restaurant item added successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
