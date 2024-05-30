package main

import (
	"context"
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
	if response.RestaurantName != "" {
		restaurant.Name = response.RestaurantName
	}
	if response.RestaurantAddress != "" {
		restaurant.Address = response.RestaurantAddress
	}
	if response.RestaurantPhone != "" {
		restaurant.Phone = response.RestaurantPhone
	}
	if response.RestaurantAvailability != "" {
		restaurant.Availability = response.RestaurantAvailability
	}

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
