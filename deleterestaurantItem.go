package main

import (
	"context"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func (*RestaurantService) DeleteRestaurantItem(ctx context.Context, request *restaurantpb.DeleteRestaurantItemRequest) (*restaurantpb.DeleteRestaurantItemResponse, error) {
	logger.Info("Received DeleteRestaurantItem request",
		zap.String("restaurantName", request.RestaurantName),
		zap.String("restaurantItemName", request.RestaurantItemName))

	// Get the user email from the context
	userEmail, emailCtxError := ctx.Value("userEmail").(string)
	userRole, roleCtxError := ctx.Value("userRole").(string)

	if !emailCtxError || !roleCtxError {
		logger.Error("Failed to get user email or role from context")
		return &restaurantpb.DeleteRestaurantItemResponse{
			Message:    "Failed to get user email and role from context",
			Error:      "Internal Server Error",
			StatusCode: int64(500),
		}, nil
	}

	logger.Info("Context values retrieved", zap.String("userEmail", userEmail), zap.String("userRole", userRole))
	if userRole != model.AdminRole {
		logger.Warn("Permission denied", zap.String("userRole", userRole))
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "You do not have permission to perform this action. Only admin can add a delete restaurant item",
			StatusCode: 403,
			Error:      "Forbidden",
		}, nil
	}

	// validate the restaurant item fields
	if request.RestaurantName == "" || request.RestaurantItemName == "" {
		logger.Warn("Invalid restaurant item data provided")
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Invalid restaurant item data provided",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}
	restaurantName := strings.ReplaceAll(request.RestaurantName, "-", " ")
	restaurantItemName := strings.ReplaceAll(request.RestaurantItemName, "-", " ")

	// check if user own's this restaurant
	var restaurant model.Restaurant
	primaryKeyRes := restaurantDBConnector.Where("name = ?", restaurantName).First(&restaurant)
	if primaryKeyRes.Error != nil || restaurant.RestaurantOwnerMail != userEmail {
		logger.Warn("Restaurant does not exist or user is not the owner",
			zap.String("userEmail", userEmail),
			zap.String("restaurantName", restaurantName))
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Restaurant does not exist OR you are not the owner of this restaurant",
			StatusCode: 404,
			Error:      "Resource not found or forbidden",
		}, nil
	}

	// delete the restaurant item
	var restaurantItem model.RestaurantItem
	primaryKey := restaurantItemDBConnector.Where("item_name = ? AND restaurant_id = ?", restaurantItemName, restaurant.ID).First(&restaurantItem)

	if primaryKey.Error != nil {
		logger.Warn("Restaurant item does not exist", zap.String("restaurantItemName", restaurantItemName))
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Restaurant item does not exist",
			StatusCode: 404,
			Error:      "Resource not found or forbidden",
		}, nil
	}
	err := restaurantItemDBConnector.Delete(&restaurantItem)
	if err.Error != nil {
		logger.Error("Failed to delete restaurant item", zap.String("restaurantItemName", restaurantItemName), zap.Error(err.Error))
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Failed to delete restaurant item",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}
	logger.Info("Restaurant item deleted successfully", zap.String("restaurantItemId", strconv.FormatUint(uint64(restaurantItem.ID), 10)))
	return &restaurantpb.DeleteRestaurantItemResponse{
		Data: &restaurantpb.DeleteRestaurantItemResponseData{
			RestaurantItemId: strconv.FormatUint(uint64(restaurantItem.ID), 10),
		},
		Message:    "Restaurant item deleted successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
