package main

import (
	"context"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"
)

func (*RestaurantService) GetAllRestaurantItems(ctx context.Context, request *restaurantpb.GetAllRestaurantItemsRequest) (*restaurantpb.GetAllRestaurantItemsResponse, error) {

	// validate the restaurant id
	if request.RestaurantId == "" {
		return &restaurantpb.GetAllRestaurantItemsResponse{
			Data: &restaurantpb.GetAllRestaurantItemsResponseData{
				TotalRestaurantItems: 0,
				RestaurantItems: nil,
			},
			Message:    "Invalid restaurant id provided",
			StatusCode: int64(400),
			Error:      "Bad Request",
		}, nil
	}	
	// fetch restaurant from restaurantDB.
	var restaurant model.Restaurant
	primaryKeyRes := restaurantDBConnector.Where("id = ?", request.RestaurantId).First(&restaurant)
	
	// check if the restaurant is exist or not.
	if primaryKeyRes.Error != nil {
		return &restaurantpb.GetAllRestaurantItemsResponse{
			Data: &restaurantpb.GetAllRestaurantItemsResponseData{
				TotalRestaurantItems: 0,
				RestaurantItems: nil,
			},
			Message:    "Restaurant Does not exist",
			StatusCode: int64(404),
			Error:      "Bad Request",
		}, nil
	}
	var restaurantItems []model.RestaurantItem
	err := restaurantItemDBConnector.Where("restaurant_id = ?", restaurant.ID).Find(&restaurantItems)
	if err.Error != nil {
		return &restaurantpb.GetAllRestaurantItemsResponse{
			Data: &restaurantpb.GetAllRestaurantItemsResponseData{
				TotalRestaurantItems: 0,
				RestaurantItems: nil,
			},
			Message:    "Failed to get restaurant items",
			StatusCode: int64(500),
			Error:      "Internal Server Error",
		}, nil
	}
	restaurantItemsResponse := []*restaurantpb.RestaurantItem{}
	for _, restaurantItem := range restaurantItems {
		restaurantItemsResponse = append(restaurantItemsResponse, &restaurantpb.RestaurantItem{
			RestaurantItemId:       strconv.FormatUint(uint64(restaurantItem.ID), 10),
			RestaurantItemName:     restaurantItem.ItemName,
			RestaurantItemPrice:    restaurantItem.ItemPrice,
			RestaurantItemImageUrl: restaurantItem.ImageUrl,
			RestaurantItemCategory: restaurantItem.Category,
			RestaurantItemCuisineType: restaurantItem.CuisineType,
			RestaurantItemVeg: restaurantItem.Veg,
			RestaurantName: restaurant.Name,
		})
	}
	return &restaurantpb.GetAllRestaurantItemsResponse{
		Data: &restaurantpb.GetAllRestaurantItemsResponseData{
			TotalRestaurantItems: int64(len(restaurantItems)),
			RestaurantItems: restaurantItemsResponse,
		},
		Message:    "Restaurant items fetched successfully",
		StatusCode: int64(200),
		Error:      "",
	}, nil
}
