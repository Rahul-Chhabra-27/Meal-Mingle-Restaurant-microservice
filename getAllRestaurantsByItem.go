package main

import (
	"context"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"
)

func (*RestaurantService) GetRestaurantsByRestaurantsItem(ctx context.Context,
	request *restaurantpb.GetRestaurantsByRestaurantItemRequest) (*restaurantpb.GetRestaurantsByRestaurantItemResponse, error) {

	if request.RestaurantItemName == "" {
		return &restaurantpb.GetRestaurantsByRestaurantItemResponse{
			Data:       nil,
			Message:    "Invalid Field, restaurant item name is required",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}
	var restaurantItems []model.RestaurantItem
	err := restaurantItemDBConnector.Where("item_name = ?", request.RestaurantItemName).Find(&restaurantItems)
	if err.Error != nil {
		return &restaurantpb.GetRestaurantsByRestaurantItemResponse{
			Data:       nil,
			Message:    "Failed to get restaurant items",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}
	// find all the restaurants that have this item
	var restaurantsResponse []*restaurantpb.Restaurant
	for _, restaurantItem := range restaurantItems {
		var restaurant model.Restaurant
		var restaurantAddress model.Address
		restaurantError := restaurantDBConnector.Where("id = ?", restaurantItem.RestaurantId).First(&restaurant)
		restaurantAddressError := restaurantAddressDBConnector.Where("restaurant_id = ?", restaurant.ID).First(&restaurantAddress)

		if restaurantError.Error != nil || restaurantAddressError.Error != nil {
			return &restaurantpb.GetRestaurantsByRestaurantItemResponse{
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
	}
	return &restaurantpb.GetRestaurantsByRestaurantItemResponse{
		Data: &restaurantpb.GetRestaurantsByRestaurantItemResponseData{
			TotalRestaurants: int64(len(restaurantsResponse)),
			Restaurants:      restaurantsResponse,
		},
		Message:    "Restaurants fetched successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
