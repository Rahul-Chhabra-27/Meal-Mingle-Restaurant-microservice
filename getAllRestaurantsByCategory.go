package main

import (
	"context"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"
)

func (*RestaurantService) GetRestaurantsByItemCategory(ctx context.Context,
	request *restaurantpb.GetRestaurantsByItemCategoryRequest) (*restaurantpb.GetRestaurantsByItemCategoryResponse, error) {

	if request.Category == "" {
		return &restaurantpb.GetRestaurantsByItemCategoryResponse{
			Data:       nil,
			Message:    "Invalid Field, category field is required",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}
	var restaurantItems []model.RestaurantItem
	err := restaurantItemDBConnector.Where("category = ?", request.Category).Find(&restaurantItems)
	if err.Error != nil {
		return &restaurantpb.GetRestaurantsByItemCategoryResponse{
			Data:       nil,
			Message:    "Failed to get restaurant items",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}
	// find all the restaurants that have this item
	// Extract unique RestaurantId's from restaurantItems
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
			return &restaurantpb.GetRestaurantsByItemCategoryResponse{
				Data:       nil,
				Message:    "Failed to get restaurant or restaurant address",
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
				City:       restaurantAddress.City,
				Country:    restaurantAddress.Country,
				Pincode:    restaurantAddress.Pincode,
				StreetName: restaurantAddress.StreetName,
			},
		})
		restaurantIds[restaurantItem.RestaurantId] = true
	}
	return &restaurantpb.GetRestaurantsByItemCategoryResponse{
		Data: &restaurantpb.GetRestaurantsByItemCategoryResponseData{
			TotalRestaurants: int64(len(restaurantsResponse)),
			Restaurants:      restaurantsResponse,
		},
		Message:    "Restaurants fetched successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
