package main

import (
	"context"
	"fmt"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
)

func (*RestaurantService) UpdateRestaurant(ctx context.Context, request *restaurantpb.UpdateRestaurantRequest) (*restaurantpb.UpdateRestaurantResponse, error) {
	// fetch the user email from the context
	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.UpdateRestaurantResponse{
			Data:       nil,
			Message:    "",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}
	// fetch restaurant from restaurantDB
	var restaurant model.Restaurant
	primaryKeyRes := restaurantDBConnector.Where("id = ?", request.Restaurant.RestaurantId).First(&restaurant)

	// check if the restaurant is exist or not
	if primaryKeyRes.Error != nil || restaurant.RestaurantOwnerMail != userEmail {
		return &restaurantpb.UpdateRestaurantResponse{
			Data: 	 nil,
			Message:    "The requested restaurant does not exist or you do not have permission to access it.",
			StatusCode: 404,
			Error:      "Resource not found or forbidden",
		}, nil
	}
	// fetch the restaurant Address
	var restaurantAddress model.Address
	restaurantAddress.City = request.Restaurant.RestaurantAddress.City
	restaurantAddress.Country = request.Restaurant.RestaurantAddress.Country
	restaurantAddress.Pincode = request.Restaurant.RestaurantAddress.Pincode
	restaurantAddress.StreetName = request.Restaurant.RestaurantAddress.StreetName

	err := restaurantAddressDBConnector.Where("restaurant_id = ?", restaurant.ID).First(&restaurantAddress)
	if err.Error != nil {
		fmt.Println("[ UpdateRestaurant ] Failed to get restaurant address from the database", err.Error)
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "Failed to get restaurant address from the database",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}
	if !config.ValidateRestaurantFields(restaurant.Name, restaurantAddress, restaurant.Phone, restaurant.Availability, restaurant.ImageUrl, restaurant.OperationHours, restaurant.OperationDays) {
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "Invalid restaurant data provided.Some fields might be missing or invalid",
			StatusCode: 400,
			Error:      "Invalid restaurant fields",
		}, nil
	}

	if !config.ValidateRestaurantPhone(restaurant.Phone) {
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "Invalid phone number format",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}

	restaurant.Name = request.Restaurant.RestaurantName
	restaurant.Availability = request.Restaurant.RestaurantAvailability
	restaurant.Phone = request.Restaurant.RestaurantPhoneNumber
	restaurant.Rating = request.Restaurant.RestaurantRating
	restaurant.ImageUrl = request.Restaurant.RestaurantImageUrl
	restaurant.OperationDays = request.Restaurant.RestaurantOperationDays
	restaurant.OperationHours = request.Restaurant.RestaurantOperationHours
	restaurant.RestaurantOwnerMail = userEmail

	updateRestaurantError := restaurantDBConnector.Save(&restaurant)
	updateAddressError := restaurantAddressDBConnector.Save(&restaurantAddress)

	if updateAddressError.Error != nil || updateRestaurantError.Error != nil {
		fmt.Println("[ UpdateRestaurant ] Failed to update restaurant", updateRestaurantError.Error)
		fmt.Println("[ UpdateRestaurant ] Failed to update restaurant address", updateAddressError.Error)
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "Failed to save the item to the database. Please try again later.",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}
	return &restaurantpb.UpdateRestaurantResponse{
		Data: &restaurantpb.UpdateRestaurantData {
			Restaurant: request.Restaurant,
		},
		Message:    "Restaurant updated successfully",
		StatusCode: 200,
		Error:      "",
	}, nil

}
