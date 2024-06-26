package main

import (
	"context"
	"fmt"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"

	"go.uber.org/zap"
)

func (*RestaurantService) AddRestaurant(ctx context.Context, request *restaurantpb.AddRestaurantRequest) (*restaurantpb.AddRestaurantResponse, error) {
	logger.Info("Received AddRestaurant request",
		zap.String("restaurantName", request.Restaurant.RestaurantName))

	userEmail, emailCtxError := ctx.Value("userEmail").(string)
	userRole, roleCtxError := ctx.Value("userRole").(string)

	if !emailCtxError || !roleCtxError {
		logger.Error("Failed to get user email or role from context")
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.AddRestaurantResponse{
			Message:    "Failed to get user mail from context",
			Error:      "Internal Server Error",
			StatusCode: int64(500),
		}, nil
	}
	logger.Info("Context values retrieved", zap.String("userEmail", userEmail), zap.String("userRole", userRole))
	if userRole != model.AdminRole {
		logger.Warn("Permission denied", zap.String("userRole", userRole))
		return &restaurantpb.AddRestaurantResponse{
			Data:       nil,
			Message:    "You do not have permission to perform this action. Only admin can add a restaurant",
			StatusCode: 403,
			Error:      "Forbidden",
		}, nil
	}

	var restaurantAddress model.Address
	var restaurant model.Restaurant
	if request.Restaurant != nil && request.Restaurant.RestaurantAddress != nil {
		restaurantAddress.City = request.Restaurant.RestaurantAddress.City
		restaurantAddress.Country = request.Restaurant.RestaurantAddress.Country
		restaurantAddress.Pincode = request.Restaurant.RestaurantAddress.Pincode
		restaurantAddress.StreetName = request.Restaurant.RestaurantAddress.StreetName
	} else {
		logger.Warn("Invalid restaurant address data provided")
		return &restaurantpb.AddRestaurantResponse{
			Data:       nil,
			Message:    "Invalid restaurant address data provided. Some fields might be missing or invalid",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}
	restaurant.Name = request.Restaurant.RestaurantName
	restaurant.Availability = request.Restaurant.RestaurantAvailability
	restaurant.Phone = request.Restaurant.RestaurantPhoneNumber
	restaurant.Rating = request.Restaurant.RestaurantRating
	restaurant.ImageUrl = request.Restaurant.RestaurantImageUrl
	restaurant.OperationDays = request.Restaurant.RestaurantOperationDays
	restaurant.OperationHours = request.Restaurant.RestaurantOperationHours
	restaurant.RestaurantOwnerMail = userEmail

	logger.Info("Restaurant data populated", zap.String("restaurantName", restaurant.Name))

	if !config.ValidateRestaurantFields(restaurant.Name, restaurantAddress,
		restaurant.Phone, restaurant.Availability,
		restaurant.ImageUrl, restaurant.OperationDays,
		restaurant.OperationHours, restaurant.Rating) {
		logger.Warn("Invalid restaurant data provided", zap.String("restaurantName", restaurant.Name))
		
		return &restaurantpb.AddRestaurantResponse{
			Data:       nil,
			Message:    "Invalid restaurant data provided. Some fields might be missing or invalid",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}

	if !config.ValidateRestaurantPhone(restaurant.Phone) {
		logger.Warn("Invalid phone number format", zap.String("phone", restaurant.Phone))
		return &restaurantpb.AddRestaurantResponse{
			Data:       nil,
			Message:    "Invalid phone number format",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}
	var existingRestaurant model.Restaurant
	restaurantNotFoundErr := restaurantDBConnector.Where("name = ?", restaurant.Name).First(&existingRestaurant).Error
	if restaurantNotFoundErr == nil {
		return &restaurantpb.AddRestaurantResponse{
			Data:       nil,
			Message:    "Same name restaurant exists. Please check the restaurant name and try again.",
			StatusCode: int64(409),
			Error:      "Restaurant creation failed",
		}, nil
	}
	primaryKey := restaurantDBConnector.Create(&restaurant)
	if primaryKey.Error != nil {
		logger.Error("[ AddRestaurant ] Failed to add restaurant", zap.Error(primaryKey.Error))
		return &restaurantpb.AddRestaurantResponse{
			Data:       nil,
			Message:    "Failed to add restaurant",
			StatusCode: 409,
			Error:      "The provided phone number is already associated with an account",
		}, nil
	}
	restaurantAddress.RestaurantId = restaurant.ID
	err := restaurantAddressDBConnector.Create(&restaurantAddress)
	if err.Error != nil {
		logger.Error("[ AddRestaurant ] Failed to add restaurant address", zap.Error(err.Error))
		return &restaurantpb.AddRestaurantResponse{
			Data:       nil,
			Message:    "Failed to add restaurant address",
			StatusCode: 500,
			Error:      err.Error.Error(),
		}, nil
	}
	logger.Info("Restaurant added successfully", zap.String("restaurantId", strconv.FormatUint(uint64(restaurant.ID), 10)))
	RestaurantResponse := request.Restaurant
	RestaurantResponse.RestaurantId = strconv.FormatUint(uint64(restaurant.ID), 10)

	return &restaurantpb.AddRestaurantResponse{
		Data: &restaurantpb.AddRestaurantResponseData{
			Restaurant: RestaurantResponse,
		},
		Message:    "Restaurant added successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
