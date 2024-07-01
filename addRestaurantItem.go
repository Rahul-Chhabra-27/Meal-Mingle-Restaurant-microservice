package main

import (
	"context"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"

	"go.uber.org/zap"
)

func (*RestaurantService) AddRestaurantItem(ctx context.Context, request *restaurantpb.AddRestaurantItemRequest) (*restaurantpb.AddRestaurantItemResponse, error){
	logger.Info("Received AddRestaurantItem request",
		zap.String("restaurantItemName", request.RestaurantItem.RestaurantItemName))

	userEmail, emailCtxError := ctx.Value("userEmail").(string)
	userRole, roleCtxError := ctx.Value("userRole").(string)

	if !emailCtxError || !roleCtxError {
		logger.Error("Failed to get user email or role from context")
		return &restaurantpb.AddRestaurantItemResponse{
			Data:       nil,
			Message:    "Failed to get user mail from context",
			Error:      "Internal Server Error",
			StatusCode: int64(500),
		}, nil
	}

	logger.Info("Context values retrieved", zap.String("userEmail", userEmail),
		zap.String("userRole", userRole))

	if userRole != model.AdminRole {
		logger.Warn("Permission denied", zap.String("userRole", userRole))
		return &restaurantpb.AddRestaurantItemResponse{
			Data:       nil,
			Message:    "You do not have permission to perform this action. Only admin can add a restaurant item",
			StatusCode: 403,
			Error:      "Forbidden",
		}, nil
	}

	if request.RestaurantItem == nil ||
		!config.ValidateRestaurantItemFields(request.RestaurantItem.RestaurantItemName,
			request.RestaurantItem.RestaurantItemImageUrl, request.RestaurantItem.RestaurantItemPrice,
			request.RestaurantItem.RestaurantItemCategory,
			request.RestaurantItem.RestaurantItemCuisineType,
			request.RestaurantItem.RestaurantName) {
		logger.Warn("Invalid restaurant item data provided")
		return &restaurantpb.AddRestaurantItemResponse{
			Data:       nil,
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
		logger.Warn("Restaurant does not exist or you are not authorized to modify this restaurant's data",
			zap.String("restaurantName", request.RestaurantItem.RestaurantName))
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
		logger.Error("Failed to create restaurant item", zap.String("restaurantItemName", request.RestaurantItem.RestaurantItemName), zap.Error(result.Error))
		return &restaurantpb.AddRestaurantItemResponse{
			Message:    "A food item with similar details might already exist on this restaurant's menu.",
			StatusCode: 409,
			Error:      "Food item creation failed",
		}, nil
	}
	logger.Info("Restaurant item added successfully", zap.String("restaurantItemName", request.RestaurantItem.RestaurantItemName))
	// Return a success message if the restaurant item is created successfully
	restaurantItemResponse := addRestaurantItemResponseData(restaurant.ID,restaurantItem.ID, request.RestaurantItem)
	return &restaurantpb.AddRestaurantItemResponse{
		Data: &restaurantpb.AddRestaurantItemResponseData{
			RestaurantItem: restaurantItemResponse,
		},
		Message:    "Restaurant item added successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}

func addRestaurantItemResponseData(restaurantId uint,restaurantItemId uint ,restaurantItemData *restaurantpb.AddRestaurantItemData) *restaurantpb.RestaurantItem {
	var restaurantItemResponse = &restaurantpb.RestaurantItem{
		RestaurantItemId:          strconv.FormatUint(uint64(restaurantItemId), 10),
		RestaurantItemName:        restaurantItemData.RestaurantItemName,
		RestaurantItemPrice:       restaurantItemData.RestaurantItemPrice,
		RestaurantItemImageUrl:    restaurantItemData.RestaurantItemImageUrl,
		RestaurantItemCategory:    restaurantItemData.RestaurantItemCategory,
		RestaurantItemCuisineType: restaurantItemData.RestaurantItemCuisineType,
		RestaurantItemVeg:         restaurantItemData.RestaurantItemVeg,
		RestaurantId:              strconv.FormatUint(uint64(restaurantId), 10),
	}
	return restaurantItemResponse
}