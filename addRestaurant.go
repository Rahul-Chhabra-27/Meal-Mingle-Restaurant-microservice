package main

import (
	"context"
	"fmt"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"

	"google.golang.org/grpc/codes"
)

func (*RestaurantService) AddRestaurant(ctx context.Context, response *restaurantpb.AddRestaurantRequest) (*restaurantpb.AddRestaurantResponse, error) {
	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.AddRestaurantResponse{Message: "", Error: "Internal Server Error", StatusCode: int64(codes.Internal)}, nil
	}
	if !config.ValidateRestaurantFields(response.RestaurantName, response.RestaurantCity,response.RestaurantAddress, response.RestaurantPhone, response.RestaurantAvailability) {

		return &restaurantpb.AddRestaurantResponse{
			Message:    "",
			StatusCode: int64(codes.InvalidArgument),
			Error:      "Invalid restaurant fields",
		}, nil
	}
	if !config.ValidateRestaurantPhone(response.RestaurantPhone) {
		return &restaurantpb.AddRestaurantResponse{
			Message:    "",
			StatusCode: int64(codes.InvalidArgument),
			Error:      "Invalid phone number",
		}, nil
	}
	var restaurant model.Restaurant
	restaurant.Name = response.RestaurantName
	restaurant.City = response.RestaurantCity
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
