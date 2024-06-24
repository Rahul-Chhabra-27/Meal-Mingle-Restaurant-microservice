package main

import (
	"context"
	"fmt"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"
)

func (*RestaurantService) AddRestaurantItem(ctx context.Context, request *restaurantpb.AddRestaurantItemRequest) (*restaurantpb.AddRestaurantItemResponse, error) {
	userEmail, emailCtxError := ctx.Value("userEmail").(string)
	userRole, roleCtxError := ctx.Value("userRole").(string)
	
	if !emailCtxError || !roleCtxError {
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.AddRestaurantItemResponse{ 
			Data: nil,
			Message: "Failed to get user mail from context",
			Error: "Internal Server Error", 
			StatusCode: int64(500),
		}, nil
	}
	
	if userRole != model.AdminRole {
		return &restaurantpb.AddRestaurantItemResponse{
			Data:       nil,
			Message:    "You do not have permission to perform this action. Only admin can add a restaurant item",
			StatusCode: 403,
			Error:      "Forbidden",
		}, nil
	}

	if  request.RestaurantItem == nil || 
		!config.ValidateRestaurantItemFields(request.RestaurantItem.RestaurantItemName,
		request.RestaurantItem.RestaurantItemImageUrl, request.RestaurantItem.RestaurantItemPrice,
		request.RestaurantItem.RestaurantItemCategory,
		request.RestaurantItem.RestaurantItemCuisineType, 
		request.RestaurantItem.RestaurantName) {
		return &restaurantpb.AddRestaurantItemResponse{
			Data: 	 nil,
			Message:    "Invalid restaurant item data provided.",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}
	var restaurantItem model.RestaurantItem
	restaurantItem.ItemPrice = request.RestaurantItem.RestaurantItemPrice
	restaurantItem.ImageUrl = request.RestaurantItem.RestaurantItemImageUrl
	restaurantItem.ItemName = request.RestaurantItem.RestaurantItemName
	restaurantItem.Category = request.RestaurantItem.RestaurantItemCategory
	restaurantItem.CuisineType = request.RestaurantItem.RestaurantItemCuisineType
	restaurantItem.Veg = request.RestaurantItem.RestaurantItemVeg

	// fetch restaurant from restaurantDB
	var restaurant model.Restaurant
	primaryKey := restaurantDBConnector.Where("name = ?", request.RestaurantItem.RestaurantName).First(&restaurant)
	// check if the restaurant is exist or not.
	if primaryKey.Error != nil || restaurant.RestaurantOwnerMail != userEmail {
		return &restaurantpb.AddRestaurantItemResponse{
			Data:       nil,
			Message:    "You are not authorized to modify this restaurant's data Or Restaurant does not exist",
			StatusCode: 403,
			Error:      "Forbidden",
		}, nil
	}
	restaurantItem.RestaurantId = restaurant.ID
	// Create a new restaurant item in the database and return the primary key if successful or an error if it fails
	result := restaurantItemDBConnector.Create(&restaurantItem)
	// Check if there is an error while creating the restaurant item
	if result.Error != nil {
		return &restaurantpb.AddRestaurantItemResponse{
			Message:    "A food item with similar details might already exist on this restaurant's menu.",
			StatusCode: 409,
			Error:      "Food item creation failed",
		}, nil
	}
	// Return a success message if the restaurant item is created successfully
	restaurantItemResponse := request.RestaurantItem
	restaurantItemResponse.RestaurantItemId = strconv.FormatUint(uint64(restaurantItem.ID), 10)
	return &restaurantpb.AddRestaurantItemResponse{
		Data: &restaurantpb.AddRestaurantItemResponseData{
			RestaurantItem: restaurantItemResponse,
		},
		Message:    "Restaurant item added successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
