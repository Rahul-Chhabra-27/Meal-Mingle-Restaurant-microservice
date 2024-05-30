package main

import (
	"context"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strings"

	"google.golang.org/grpc/codes"
)

func (*RestaurantService) GetAllRestaurantItems(ctx context.Context, response *restaurantpb.GetAllRestaurantItemsRequest) (*restaurantpb.GetAllRestaurantItemsResponse, error) {
	// fetch restaurant from restaurantDB.
	var restaurant model.Restaurant
	restaurantName := strings.ReplaceAll(response.RestaurantName, "-", " ")
	primaryKeyRes := restaurantDBConnector.Where("name = ?", restaurantName).First(&restaurant)
	// check if the restaurant is exist or not.
	if primaryKeyRes.Error != nil {
		return &restaurantpb.GetAllRestaurantItemsResponse{
			RestaurantItems: nil,
			Message:         "",
			StatusCode:      int64(codes.NotFound),
			Error:           "Restaurant Does not exist",
		}, nil
	}
	var restaurantItems []model.RestaurantItem
	err := restaurantItemDBConnector.Where("restaurant_id = ?", restaurant.ID).Find(&restaurantItems)
	if err.Error != nil {
		return &restaurantpb.GetAllRestaurantItemsResponse{
			RestaurantItems: nil,
			Message:         "",
			StatusCode:      int64(codes.Internal),
			Error:           "Failed to get restaurant items",
		}, nil
	}
	restaurantItemsResponse := []*restaurantpb.RestaurantItem{}
	for _, restaurantItem := range restaurantItems {
		restaurantItemsResponse = append(restaurantItemsResponse, &restaurantpb.RestaurantItem{
			RestaurantItemName:     restaurantItem.ItemName,
			RestaurantItemPrice:    restaurantItem.ItemPrice,
			RestaurantItemImageUrl: restaurantItem.ImageUrl,
			RestaurantName:         restaurantName,
		})
	}
	return &restaurantpb.GetAllRestaurantItemsResponse{
		RestaurantItems: restaurantItemsResponse,
		Message:         "Restaurant items fetched successfully",
		StatusCode:      int64(codes.OK),
		Error:           "",
	}, nil
}
