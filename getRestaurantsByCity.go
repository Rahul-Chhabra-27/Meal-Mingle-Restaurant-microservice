package main

import (
	"context"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"

	"go.uber.org/zap"
)

func (*RestaurantService) GetRestaurantsByCity(ctx context.Context, request *restaurantpb.GetRestaurantsByCityRequest) (*restaurantpb.GetRestaurantsByCityResponse, error) {

	logger.Info("Received GetRestaurantsByCity request",
		zap.String("City", request.City))

	if request.City == "" {
		logger.Warn("City field is required but not provided")
		return &restaurantpb.GetRestaurantsByCityResponse{
			Data:       nil,
			Message:    "Invalid Field, city field is required",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}
	city := request.City
	var restaurantAddress []model.Address
	err := restaurantAddressDBConnector.Where("city = ?", city).Find(&restaurantAddress)
	if err.Error != nil {
		logger.Error("Failed to get restaurants from the database", zap.Error(err.Error))
		return &restaurantpb.GetRestaurantsByCityResponse{
			Data:       nil,
			Message:    "Failed to get restaurants from the database",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}
	logger.Info("Fetched restaurant addresses", zap.Int("Count", len(restaurantAddress)))

	var restaurantsResponse []*restaurantpb.Restaurant
	for _, address := range restaurantAddress {
		// fetch all restaurant details from the database filter by city.
		var restaurant model.Restaurant
		restaurantErr := restaurantDBConnector.Where("id = ?", address.RestaurantId).First(&restaurant).Error
		if restaurantErr != nil {
			logger.Error("Failed to get restaurant from the database", zap.Error(restaurantErr))
			return &restaurantpb.GetRestaurantsByCityResponse{
				Data:       nil,
				Message:    "Failed to get restaurant from the database",
				StatusCode: 500,
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
			RestaurantAddress: &restaurantpb.Address{
				City:       address.City,
				Country:    address.Country,
				Pincode:    address.Pincode,
				StreetName: address.StreetName,
			},
		})
	}
	logger.Info("Fetched restaurants by city successfully", zap.Int("TotalRestaurants", len(restaurantsResponse)))
	return &restaurantpb.GetRestaurantsByCityResponse{
		Data: &restaurantpb.GetRestaurantsByCityResponseData{
			TotalRestaurants: int64(len(restaurantsResponse)),
			Restaurants:      restaurantsResponse,
		},
		Message:    "Restaurants fetched successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
