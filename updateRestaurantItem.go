package main

import (
	"context"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"

	"go.uber.org/zap"
)

func (*RestaurantService) UpdateRestaurantItem(ctx context.Context, request *restaurantpb.UpdateRestaurantItemRequest) (*restaurantpb.UpdateRestaurantItemResponse, error) {

	logger.Info("Received UpdateRestaurantItem request",
	zap.Any("Request", request))

	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		logger.Error("Failed to get user email from context")
		return &restaurantpb.UpdateRestaurantItemResponse{
			Message:    "Failed to get user email from context",
			StatusCode: StatusInternalServerError,
			Error:      "Internal Server Error",
		}, nil
	}
	// validate the restaurant item fields
	if request.RestaurantItem == nil || !config.ValidateRestaurantItemFields(request.RestaurantItem.RestaurantItemName,
		strconv.FormatInt(request.RestaurantItem.RestaurantItemPrice, 10),
		request.RestaurantItem.RestaurantItemPrice, request.RestaurantItem.GetRestaurantItemCategory(),
		request.RestaurantItem.RestaurantItemCuisineType, request.RestaurantItem.RestaurantItemImageUrl) {
			logger.Warn("Invalid restaurant item data provided", 
			zap.Any("RestaurantItem", request.RestaurantItem))

		return &restaurantpb.UpdateRestaurantItemResponse{
			Message:    "Invalid restaurant item data provided.Some fields might be missing or invalid",
			StatusCode: StatusBadRequest,
			Error:      "Invalid restaurant item fields",
		}, nil
	}
	// fetch restaurant from restaurantDB
	var restaurant model.Restaurant
	primaryKeyRes := restaurantDBConnector.Where("id = ?", request.RestaurantItem.RestaurantId).First(&restaurant)
	// check if the restaurant is exist or nor
	if primaryKeyRes.Error != nil || restaurant.RestaurantOwnerMail != userEmail {
		logger.Warn("Unauthorized update attempt", 
		zap.String("RestaurantOwnerMail", restaurant.RestaurantOwnerMail), 
		zap.String("UserEmail", userEmail))

		return &restaurantpb.UpdateRestaurantItemResponse{
			Message:    "You do not have permission to perform this action. Only restaurant owner can update the restaurant item",
			StatusCode: StatusForbidden,
			Error:      "Resource not found or forbidden",
		}, nil
	}

	var restaurantItem model.RestaurantItem
	primaryKey := restaurantItemDBConnector.Where("id = ? AND restaurant_id = ?", request.RestaurantItem.RestaurantItemId,
		restaurant.ID).First(&restaurantItem)
	if primaryKey.Error != nil {
		logger.Error("Failed to fetch restaurant item from database", 
		zap.Error(primaryKey.Error))
		return &restaurantpb.UpdateRestaurantItemResponse{
			Message:    "Restaurant Item does not exist",
			StatusCode: StatusNotFound,
			Error:      "Not Found",
		}, nil
	}
	restaurantItem.ItemName = request.RestaurantItem.RestaurantItemName
	restaurantItem.ItemPrice = request.RestaurantItem.RestaurantItemPrice
	restaurantItem.Category = request.RestaurantItem.GetRestaurantItemCategory()
	restaurantItem.CuisineType = request.RestaurantItem.RestaurantItemCuisineType
	restaurantItem.Veg = request.RestaurantItem.RestaurantItemVeg
	restaurantItem.ImageUrl = request.GetRestaurantItem().GetRestaurantItemImageUrl()

	err := restaurantItemDBConnector.Save(&restaurantItem)
	if err.Error != nil {
		logger.Error("Failed to update restaurant item", zap.Error(err.Error))
		return &restaurantpb.UpdateRestaurantItemResponse{
			Message:    "Failed to update restaurant item, this can be due to same item name already exist in the restaurant or some other issue.",
			StatusCode: StatusInternalServerError,
			Error:      "Internal Server Error",
		}, nil
	}
	logger.Info("Restaurant item updated successfully", zap.Any("RestaurantItem", request.RestaurantItem))
	return &restaurantpb.UpdateRestaurantItemResponse{
		Data: &restaurantpb.UpdateRestaurantItemResponseData{
			RestaurantItem: request.RestaurantItem,
		},
		Message:    "Restaurant item updated successfully",
		StatusCode: StatusOK,
		Error:      "",
	}, nil
}
