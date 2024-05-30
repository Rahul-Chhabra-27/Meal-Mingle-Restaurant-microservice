package main

import (
	"context"
	"log"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"

	"google.golang.org/grpc/codes"
)

func (*RestaurantService) AddRestaurant(ctx context.Context, response *restaurantpb.AddRestaurantRequest) (*restaurantpb.AddRestaurantResponse, error) {
	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		log.Fatalf("Failed to get user email from context")
		return &restaurantpb.AddRestaurantResponse{Message: "", Error: "Internal Server Error", StatusCode: int64(codes.Internal)}, nil
	}
	var restaurant model.Restaurant
	restaurant.Name = response.RestaurantName
	restaurant.Address = response.RestaurantAddress
	restaurant.Phone = response.RestaurantPhone
	restaurant.Availability = response.RestaurantAvailability
	restaurant.Rating = response.RestaurantRating
	restaurant.RestaurantOwnerMail = userEmail
	restaurantNotFoundErr := restaurantDBConnector.Where("name = ?", restaurant.Name).First(&restaurant).Error

	if restaurantNotFoundErr == nil {
		return &restaurantpb.AddRestaurantResponse{
			Message:    "Restaurant already exists",
			StatusCode: int64(codes.AlreadyExists),
			Error:      "",
		}, nil
	}
	primaryKey := restaurantDBConnector.Create(&restaurant)
	if primaryKey.Error != nil {
		return &restaurantpb.AddRestaurantResponse{
			Message:    "Failed to add restaurant",
			StatusCode: int64(codes.Internal),
			Error:      primaryKey.Error.Error(),
		}, nil
	}
	return &restaurantpb.AddRestaurantResponse{
		Message:    "Restaurant added successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
