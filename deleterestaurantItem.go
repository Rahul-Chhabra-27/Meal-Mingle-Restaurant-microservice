package main

import (
	"context"
	"fmt"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strings"

	"google.golang.org/grpc/codes"
)

func (*RestaurantService) DeleteRestaurantItem(ctx context.Context, response *restaurantpb.DeleteRestaurantItemRequest) (*restaurantpb.DeleteRestaurantItemResponse, error) {
	// Get the user email from the context
	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.DeleteRestaurantItemResponse{
			Message:    "",
			StatusCode: int64(codes.Internal),
			Error:      "Internal Server Error",
		}, nil
	}
	fmt.Println(userEmail)
	restaurantName := strings.ReplaceAll(response.RestaurantName, "-", " ")
	restaurantItemName := strings.ReplaceAll(response.RestaurantItemName, "-", " ")

	// check if user own's this restaurant
	var restaurant model.Restaurant
	primaryKeyRes := restaurantDBConnector.Where("name = ?", restaurantName).First(&restaurant)
	if primaryKeyRes.Error != nil || restaurant.RestaurantOwnerMail != userEmail {
		return &restaurantpb.DeleteRestaurantItemResponse{
			Message:    "",
			StatusCode: int64(codes.NotFound),
			Error:      "Restaurant Does not exist OR you are not the owner of this restaurant",
		}, nil
	}

	// delete the restaurant item
	var restaurantItem model.RestaurantItem
	primaryKey := restaurantItemDBConnector.Where("item_name = ? AND restaurant_id = ?", restaurantItemName, restaurant.ID).First(&restaurantItem)
	if primaryKey.Error != nil {
		return &restaurantpb.DeleteRestaurantItemResponse{
			Message:    "",
			StatusCode: int64(codes.Internal),
			Error:      "Restaurant item does not exist",
		}, nil
	}
	err := restaurantItemDBConnector.Delete(&restaurantItem)
	if err.Error != nil {
		return &restaurantpb.DeleteRestaurantItemResponse{
			Message:    "",
			StatusCode: 500,
			Error:      "Failed to delete restaurant item",
		}, nil
	}
	return &restaurantpb.DeleteRestaurantItemResponse{
		Message:    "Restaurant item deleted successfully",
		StatusCode: int64(codes.OK),
		Error:      "",
	}, nil
}

