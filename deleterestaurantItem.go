package main

import (
	"context"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"

	"go.uber.org/zap"
)

func (*RestaurantService) DeleteRestaurantItem(ctx context.Context, request *restaurantpb.DeleteRestaurantItemRequest) (*restaurantpb.DeleteRestaurantItemResponse, error) {
	logger.Info("Received DeleteRestaurantItem request",
		zap.String("restaurantId", request.RestaurantId),
		zap.String("restaurantItemId", request.RestaurantItemId))

	// Get the user email from the context
	userEmail, emailCtxError := ctx.Value("userEmail").(string)
	userRole, roleCtxError := ctx.Value("userRole").(string)

	if !emailCtxError || !roleCtxError {
		logger.Error("Failed to get user email or role from context")
		return &restaurantpb.DeleteRestaurantItemResponse{
			Message:    "Failed to get user email and role from context",
			Error:      "Internal Server Error",
			StatusCode: StatusInternalServerError,
		}, nil
	}

	logger.Info("Context values retrieved", zap.String("userEmail", userEmail), zap.String("userRole", userRole))
	if userRole != model.AdminRole {
		logger.Warn("Permission denied", zap.String("userRole", userRole))
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "You do not have permission to perform this action. Only admin can add a delete restaurant item",
			StatusCode: StatusForbidden,
			Error:      "Forbidden",
		}, nil
	}
	//convert string to int64
	restaurantId, _ := strconv.ParseUint(request.RestaurantId, 10, 64)
	restaurantItemId, _ := strconv.ParseUint(request.RestaurantItemId, 10, 64)
	// validate the restaurant item fields
	if restaurantId <= 0 || restaurantItemId <= 0 {
		logger.Warn("Invalid restaurant item data provided")
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Invalid restaurant item data provided",
			StatusCode: StatusBadRequest,
			Error:      "Bad Request",
		}, nil
	}

	// check if user own's this restaurant
	var restaurant model.Restaurant
	primaryKeyRes := restaurantDBConnector.Where("id = ?", restaurantId).First(&restaurant)
	if primaryKeyRes.Error != nil || restaurant.RestaurantOwnerMail != userEmail {
		logger.Warn("Restaurant does not exist or user is not the owner",
			zap.String("userEmail", userEmail),
			zap.String("restaurantId", request.RestaurantId))
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Restaurant does not exist OR you are not the owner of this restaurant",
			StatusCode: StatusForbidden,
			Error:      "Resource not found or forbidden",
		}, nil
	}

	// delete the restaurant item
	var restaurantItem model.RestaurantItem
	primaryKey := restaurantItemDBConnector.Where("id = ? AND restaurant_id = ?",
		restaurantItemId, restaurant.ID).First(&restaurantItem)

	if primaryKey.Error != nil {
		logger.Warn("Restaurant item does not exist", zap.String("restaurantItemId", request.RestaurantItemId))
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Restaurant item does not exist",
			StatusCode: StatusForbidden,
			Error:      "Resource not found or forbidden",
		}, nil
	}
	err := restaurantItemDBConnector.Delete(&restaurantItem)
	if err.Error != nil {
		logger.Error("Failed to delete restaurant item", zap.String("restaurantItemId", request.RestaurantItemId), zap.Error(err.Error))
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Failed to delete restaurant item",
			StatusCode: StatusInternalServerError,
			Error:      "Internal Server Error",
		}, nil
	}
	logger.Info("Restaurant item deleted successfully", zap.String("restaurantItemId", strconv.FormatUint(uint64(restaurantItem.ID), 10)))
	return &restaurantpb.DeleteRestaurantItemResponse{
		Data: &restaurantpb.DeleteRestaurantItemResponseData{
			RestaurantItemId: strconv.FormatUint(uint64(restaurantItem.ID), 10),
		},
		Message:    "Restaurant item deleted successfully",
		StatusCode: StatusOK,
		Error:      "",
	}, nil
}
