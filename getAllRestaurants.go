package main

import (
	"context"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"

	"google.golang.org/grpc/codes"
)

func (*RestaurantService) GetAllRestaurants(context.Context, *restaurantpb.GetAllRestaurantsRequest) (*restaurantpb.GetAllRestaurantsResponse, error) {
	var restaurants []model.Restaurant
	err := restaurantDBConnector.Find(&restaurants)
	if err.Error != nil {
		return &restaurantpb.GetAllRestaurantsResponse{
			Data: &restaurantpb.RestaurantData{
				TotalRestaurants: 0,
				Restaurants: nil,
			},
			Message:     "",
			StatusCode:  int64(codes.Internal),
			Error:       "Failed to get restaurants",
		}, nil
	}
	restaurantsResponse := []*restaurantpb.Restaurant{}
	totalRestaurants := 0;
	for _, restaurant := range restaurants {
		totalRestaurants++
		restaurantsResponse = append(restaurantsResponse, &restaurantpb.Restaurant{
			RestaurantName:         restaurant.Name,
			RestaurantAddress:      restaurant.Address,
			RestaurantPhone:        restaurant.Phone,
			RestaurantAvailability: restaurant.Availability,
			RestaurantRating:       restaurant.Rating,
			RestaurantOwnerMail: restaurant.RestaurantOwnerMail,
			RestaurantCity: 	   restaurant.City,
		})
	}
	return &restaurantpb.GetAllRestaurantsResponse{
		Data: &restaurantpb.RestaurantData{
			TotalRestaurants: int64(totalRestaurants),
			Restaurants: restaurantsResponse,
		},
		Message:     "Restaurants fetched successfully",
		StatusCode:  200,
		Error:       "",
	}, nil
}
