package main

import (
	"context"
	"fmt"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"
	"strings"
)

func (*RestaurantService) GetRestaurantsByCity(ctx context.Context, request *restaurantpb.GetRestaurantsByCityRequest) (*restaurantpb.GetRestaurantsByCityResponse, error) {

	city := request.City
	city = strings.ReplaceAll(city, "-", " ")
	var restaurantAddress []model.Address
	err := restaurantAddressDBConnector.Where("city = ?", city).Find(&restaurantAddress)
	if err.Error != nil {
		fmt.Println("[ GetRestaurantsByCity ] Failed to get restaurants from the database", err.Error)
		return &restaurantpb.GetRestaurantsByCityResponse{
			Data:       nil,
			Message:    "",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}
	var restaurantsResponse []*restaurantpb.Restaurant
	totalRestaurants := 0
	for _, address := range restaurantAddress {
		totalRestaurants++
		// fetch all restaurant details from the database filter by city.
		var restaurant model.Restaurant
		restaurantErr := restaurantDBConnector.Where("id = ?", address.RestaurantId).First(&restaurant).Error
		if restaurantErr != nil {
			fmt.Println("[ GetRestaurantsByCity ] Failed to get restaurant from the database", restaurantErr)
			return &restaurantpb.GetRestaurantsByCityResponse{
				Data:       nil,
				Message:    "",
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
	return &restaurantpb.GetRestaurantsByCityResponse{
		Data: &restaurantpb.GetRestaurantsByCityData{
			TotalRestaurants: int64(totalRestaurants),
			Restaurants:       restaurantsResponse,
		},
		Message:    "Restaurants fetched successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
