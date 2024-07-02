package main

import (
	"context"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"

	"go.uber.org/zap"
)

func (*RestaurantService) GetRestaurantsByItemCategory(ctx context.Context,
	request *restaurantpb.GetRestaurantsByItemCategoryRequest) (*restaurantpb.GetRestaurantsByItemCategoryResponse, error) {
	
	logger.Info("Received GetRestaurantsByItemCategory request", 
	zap.String("Category", request.Category))

	if request.Category == "" {
		logger.Warn("Category field is required but not provided")

		return &restaurantpb.GetRestaurantsByItemCategoryResponse{
			Data:       nil,
			Message:    "Invalid Field, category field is required",
			StatusCode: StatusBadRequest,
			Error:      "Bad Request",
		}, nil
	}
	var restaurantItems []model.RestaurantItem
	err := restaurantItemDBConnector.Where("category = ?", request.Category).Find(&restaurantItems)
	if err.Error != nil {
		logger.Error("Failed to get restaurant items from the database", zap.Error(err.Error))
		return &restaurantpb.GetRestaurantsByItemCategoryResponse{
			Data:       nil,
			Message:    "Failed to get restaurant items",
			StatusCode: StatusInternalServerError,
			Error:      "Internal Server Error",
		}, nil
	}
	logger.Info("Fetched restaurant items", zap.Int("Count", len(restaurantItems)))

	restaurantIds := make(map[uint]bool)
	var restaurantsResponse []*restaurantpb.Restaurant
	for _, restaurantItem := range restaurantItems {
		if restaurantIds[restaurantItem.RestaurantId] {
			continue
		}
		var restaurant model.Restaurant
		var restaurantAddress model.Address
		restaurantError := restaurantDBConnector.Where("id = ?", restaurantItem.RestaurantId).First(&restaurant)
		restaurantAddressError := restaurantAddressDBConnector.Where("restaurant_id = ?", restaurant.ID).First(&restaurantAddress)

		if restaurantError.Error != nil || restaurantAddressError.Error != nil {
			logger.Error("Failed to get restaurant or restaurant address from the database", 
			zap.Error(restaurantError.Error), 
			zap.Error(restaurantAddressError.Error))

			return &restaurantpb.GetRestaurantsByItemCategoryResponse{
				Data:       nil,
				Message:    "Failed to get restaurant or restaurant address",
				StatusCode: StatusInternalServerError,
				Error:      "Internal Server Error",
			}, nil
		}
		restaurantsResponse = append(restaurantsResponse, &restaurantpb.Restaurant{
			RestaurantId:             strconv.FormatUint(uint64(restaurant.ID), 10),
			RestaurantName:           restaurant.Name,
			RestaurantAvailability:   restaurant.Availability,
			RestaurantPhoneNumber:    restaurant.Phone,
			RestaurantRating:         restaurant.Rating,
			RestaurantImageUrl:       restaurant.ImageUrl,
			RestaurantOperationDays:  restaurant.OperationDays,
			RestaurantOperationHours: restaurant.OperationHours,
			RestaurantOwnerMail:      restaurant.RestaurantOwnerMail,
			RestaurantMinimumOrderAmount: restaurant.RestaurantMinimumOrderAmount,
			RestaurantDiscountPercentage: restaurant.RestaurantDiscountPercentage,
			RestaurantAddress: &restaurantpb.Address{
				City:       restaurantAddress.City,
				Country:    restaurantAddress.Country,
				Pincode:    restaurantAddress.Pincode,
				StreetName: restaurantAddress.StreetName,
			},
		})
		restaurantIds[restaurantItem.RestaurantId] = true
	}
	logger.Info("Fetched restaurants by item category successfully", zap.Int("TotalRestaurants", len(restaurantsResponse)))
	return &restaurantpb.GetRestaurantsByItemCategoryResponse{
		Data: &restaurantpb.GetRestaurantsByItemCategoryResponseData{
			TotalRestaurants: int64(len(restaurantsResponse)),
			Restaurants:      restaurantsResponse,
		},
		Message:    "Restaurants fetched successfully",
		StatusCode: StatusOK,
		Error:      "",
	}, nil
}
