package main

import (
	"context"
	"restaurant-micro/config"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"

	"google.golang.org/grpc/codes"
)

func (*RestaurantService) UpdateRestaurant(ctx context.Context, response *restaurantpb.UpdateRestaurantRequest) (*restaurantpb.UpdateRestaurantResponse, error) {
	// fetch the user email from the context
	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}
	// fetch restaurant from restaurantDB
	var restaurant model.Restaurant
	primaryKeyRes := restaurantDBConnector.Where("name = ?", response.RestaurantName).First(&restaurant)
	// check if the restaurant is exist or not
	if primaryKeyRes.Error != nil || restaurant.RestaurantOwnerMail != userEmail {
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "",
			StatusCode: 404,
			Error:      "Restaurant Does not exist OR you are not the owner of this restaurant",
		}, nil
	}
	if !config.ValidateRestaurantFields(response.RestaurantName, response.RestaurantCity, response.RestaurantAddress, response.RestaurantPhone, response.RestaurantAvailability) {
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "",
			StatusCode: int64(codes.InvalidArgument),
			Error:      "Invalid restaurant fields",
		}, nil
	}
	
	if !config.ValidateRestaurantPhone(response.RestaurantPhone) {
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "",
			StatusCode: int64(codes.InvalidArgument),
			Error:      "Invalid phone number",
		}, nil
	}
	restaurant.Name = response.RestaurantName
	restaurant.Address = response.RestaurantAddress
	restaurant.Phone = response.RestaurantPhone
	restaurant.Availability = response.RestaurantAvailability
	restaurant.City = response.RestaurantCity

	err := restaurantDBConnector.Save(&restaurant)
	if err.Error != nil {
		return &restaurantpb.UpdateRestaurantResponse{
			Message:    "",
			StatusCode: int64(codes.Internal),
			Error:      err.Error.Error(),
		}, nil
	}
	return &restaurantpb.UpdateRestaurantResponse{
		Message:    "Restaurant updated successfully",
		StatusCode: 200,
		Error:      "",
	}, nil

}
