package main

import (
	"context"
	"fmt"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
)

func (*RestaurantService) AddRestaurant(ctx context.Context, request *restaurantpb.AddRestaurantRequest) (*restaurantpb.AddRestaurantResponse, error) {
	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.AddRestaurantResponse{Message: "", Error: "Internal Server Error", StatusCode: int64(500)}, nil
	}
	var restaurantAddress model.Address;
	var restaurant model.Restaurant

	restaurant.Name = request.Restaurant.RestaurantName
	restaurant.Availability = request.Restaurant.RestaurantAvailability
	restaurant.Phone = request.Restaurant.RestaurantPhoneNumber
	restaurant.Rating = request.Restaurant.RestaurantRating
	restaurant.ImageUrl = request.Restaurant.RestaurantImageUrl
	restaurant.OperationDays = request.Restaurant.RestaurantOperationDays
	restaurant.OperationHours = request.Restaurant.RestaurantOperationHours
	restaurant.RestaurantOwnerMail = userEmail
	restaurantAddress.City = request.Restaurant.RestaurantAddress.City
	restaurantAddress.Country = request.Restaurant.RestaurantAddress.Country
	restaurantAddress.Pincode = request.Restaurant.RestaurantAddress.Pincode
	restaurantAddress.StreetName = request.Restaurant.RestaurantAddress.StreetName

	if !config.ValidateRestaurantFields(restaurant.Name, restaurantAddress, restaurant.Phone, restaurant.Availability, restaurant.ImageUrl) {

		return &restaurantpb.AddRestaurantResponse{
			Message:    "Invalid restaurant data provided.Some fields might be missing or invalid",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}
	if !config.ValidateAddressFields(restaurantAddress.City, restaurantAddress.Country, restaurantAddress.Pincode, restaurantAddress.StreetName) {
		return &restaurantpb.AddRestaurantResponse{
			Message:    "Invalid address format. Please check the address details.Some fields might be missing or invalid.",
			StatusCode: 400,
			Error:   "Bad Request",
		}, nil
	}

	if !config.ValidateRestaurantPhone(restaurant.Phone) {
		return &restaurantpb.AddRestaurantResponse{
			Message:    "Invalid phone number format",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}
	
	restaurantNotFoundErr := restaurantDBConnector.Where("name = ?", restaurant.Name).First(&restaurant).Error

	if restaurantNotFoundErr == nil {
		return &restaurantpb.AddRestaurantResponse{
			Message:    "A restaurant with similar details might already exist. Please check the restaurant name and try again.",
			StatusCode: int64(409),
			Error:      "Restaurant creation failed",
		}, nil
	}
	primaryKey := restaurantDBConnector.Create(&restaurant)
	if primaryKey.Error != nil {
		return &restaurantpb.AddRestaurantResponse{
			Message:    "Failed to add restaurant",
			StatusCode: 500,
			Error:      "The provided phone number is already associated with an account",
		}, nil
	}
	restaurantAddress.RestaurantId = restaurant.ID
	err := restaurantAddressDBConnector.Create(&restaurantAddress)
	if err.Error != nil {
		fmt.Println("[ AddRestaurant ] Failed to add restaurant address", err.Error)
		return &restaurantpb.AddRestaurantResponse{
			Message:    "Failed to add restaurant address",
			StatusCode: 500,
			Error:      err.Error.Error(),
		}, nil
	}
	return &restaurantpb.AddRestaurantResponse{
		Data: 	 &restaurantpb.AddRestaurantData{
			Restaurant: &restaurantpb.Restaurant{
				RestaurantName: 	   restaurant.Name,
				RestaurantPhoneNumber: restaurant.Phone,
				RestaurantRating:      restaurant.Rating,
				RestaurantImageUrl:    restaurant.ImageUrl,
				RestaurantOperationDays: restaurant.OperationDays,
				RestaurantOperationHours: restaurant.OperationHours,
				RestaurantAvailability: restaurant.Availability,
				RestaurantOwnerMail:  restaurant.RestaurantOwnerMail,
				RestaurantAddress: &restaurantpb.Address{
					City:      restaurantAddress.City,
					Country:   restaurantAddress.Country,
					Pincode:   restaurantAddress.Pincode,
					StreetName: restaurantAddress.StreetName,
				},	
			},
		},
		Message:    "Restaurant added successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
