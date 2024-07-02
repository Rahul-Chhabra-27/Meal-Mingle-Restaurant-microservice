package main

import (
	"context"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"

	"go.uber.org/zap"
)

func (*RestaurantService) UpdateRestaurant(ctx context.Context, request *restaurantpb.UpdateRestaurantRequest) (*restaurantpb.UpdateRestaurantResponse, error) {
	logger.Info("Received UpdateRestaurant request",
		zap.Any("Request", request))

	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		logger.Error("Failed to get user email from context")

		return &restaurantpb.UpdateRestaurantResponse{
			Data:       nil,
			Message:    "",
			StatusCode: StatusInternalServerError,
			Error:      "Internal Server Error",
		}, nil
	}
	logger.Info("Fetched user email from context", zap.String("userEmail", userEmail))
	var restaurant model.Restaurant
	var restaurantAddress model.Address

	if request.Restaurant == nil || request.Restaurant.GetRestaurantAddress() == nil {
		logger.Warn("Invalid restaurant data provided")
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "Invalid restaurant data provided",
			StatusCode: StatusBadRequest,
			Error:      "Bad Request",
		}, nil
	} else {
		restaurantAddress.City = request.Restaurant.RestaurantAddress.City
		restaurantAddress.Country = request.Restaurant.RestaurantAddress.Country
		restaurantAddress.Pincode = request.Restaurant.RestaurantAddress.Pincode
		restaurantAddress.StreetName = request.Restaurant.RestaurantAddress.StreetName
	}
	if !config.ValidateRestaurantFields(
		request.Restaurant.RestaurantName, restaurantAddress,
		request.Restaurant.RestaurantPhoneNumber,
		request.Restaurant.RestaurantAvailability, request.Restaurant.RestaurantImageUrl,
		request.Restaurant.RestaurantOperationHours,
		request.Restaurant.RestaurantOperationDays,
		request.Restaurant.RestaurantRating, 
		request.Restaurant.RestaurantMinimumOrderAmount, 
		request.Restaurant.RestaurantDiscountPercentage) ||
		!config.ValidateRestaurantPhone(request.Restaurant.RestaurantPhoneNumber) {

		logger.Warn("Invalid restaurant fields", zap.Any("Restaurant", request.Restaurant))
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "Invalid restaurant data provided.Some fields might be missing, empty or invalid, and make sure phone number contains only numbers and is 10 digits long",
			StatusCode: StatusBadRequest,
			Error:      "Invalid restaurant fields",
		}, nil
	}
	// fetch restaurant from restaurantDB
	primaryKeyRes := restaurantDBConnector.Where("id = ?", request.Restaurant.RestaurantId).First(&restaurant)
	// check if the restaurant is exist or not
	if primaryKeyRes.Error != nil || restaurant.RestaurantOwnerMail != userEmail {
		logger.Warn("Unauthorized update attempt", zap.String("RestaurantOwnerMail",
			restaurant.RestaurantOwnerMail),
			zap.String("UserEmail", userEmail))
		return &restaurantpb.UpdateRestaurantResponse{
			Data:       nil,
			Message:    "You do not have permission to perform this action. Only restaurant owner can update the restaurant",
			StatusCode: StatusUnauthorized,
			Error:      "Unauthorized",
		}, nil
	}
	err := restaurantAddressDBConnector.Where("restaurant_id = ?", restaurant.ID).First(&restaurantAddress)
	if err.Error != nil {
		logger.Error("Failed to get restaurant address from the database", zap.Error(err.Error))
		return &restaurantpb.UpdateRestaurantResponse{
			Data:       nil,
			Message:    "Failed to get restaurant address from the database",
			StatusCode: StatusInternalServerError,
			Error:      "Internal Server Error",
		}, nil
	}
	request.Restaurant.RestaurantOwnerMail = userEmail
	restaurant.Name = request.Restaurant.RestaurantName
	restaurant.Availability = request.Restaurant.RestaurantAvailability
	restaurant.Phone = request.Restaurant.RestaurantPhoneNumber
	restaurant.Rating = request.Restaurant.RestaurantRating
	restaurant.ImageUrl = request.Restaurant.RestaurantImageUrl
	restaurant.OperationDays = request.Restaurant.RestaurantOperationDays
	restaurant.OperationHours = request.Restaurant.RestaurantOperationHours
	restaurant.RestaurantOwnerMail = userEmail
	restaurant.RestaurantMinimumOrderAmount = request.Restaurant.RestaurantMinimumOrderAmount
	restaurant.RestaurantDiscountPercentage = request.Restaurant.RestaurantDiscountPercentage
	restaurantAddress.City = request.Restaurant.RestaurantAddress.City
	restaurantAddress.Country = request.Restaurant.RestaurantAddress.Country
	restaurantAddress.Pincode = request.Restaurant.RestaurantAddress.Pincode
	restaurantAddress.StreetName = request.Restaurant.RestaurantAddress.StreetName

	updateRestaurantError := restaurantDBConnector.Save(&restaurant)
	updateAddressError := restaurantAddressDBConnector.Save(&restaurantAddress)

	if updateAddressError.Error != nil || updateRestaurantError.Error != nil {
		logger.Error("Failed to update restaurant", zap.Error(updateRestaurantError.Error))
		logger.Error("Failed to update restaurant address", zap.Error(updateAddressError.Error))
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "A restaurant with the same name or phone number already exists.",
			StatusCode: StatusConflict,
			Error:      "Conflict",
		}, nil
	}
	logger.Info("Restaurant updated successfully", zap.Any("Restaurant", request.Restaurant))
	return &restaurantpb.UpdateRestaurantResponse{
		Data: &restaurantpb.UpdateRestaurantData{
			Restaurant: request.Restaurant,
		},
		Message:    "Restaurant updated successfully",
		StatusCode: StatusOK,
		Error:      "",
	}, nil

}
