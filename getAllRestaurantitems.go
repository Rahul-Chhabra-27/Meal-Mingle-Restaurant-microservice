package main

import (
	"context"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"

	"go.uber.org/zap"
)

func (*RestaurantService) GetAllRestaurantItems(ctx context.Context, request *restaurantpb.GetAllRestaurantItemsRequest) (*restaurantpb.GetAllRestaurantItemsResponse, error) {
	logger.Info("Received GetAllRestaurantItems request",
	zap.String("restaurantId", request.RestaurantId))

	// validate the restaurant id
	if request.RestaurantId == "" {
		logger.Warn("Invalid restaurant id provided")
		return &restaurantpb.GetAllRestaurantItemsResponse{
			Data: &restaurantpb.GetAllRestaurantItemsResponseData{
				TotalRestaurantItems: 0,
				RestaurantItems:      nil,
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
		logger.Warn("Restaurant does not exist", zap.String("restaurantId", request.RestaurantId))
		return &restaurantpb.GetAllRestaurantItemsResponse{
			Data: &restaurantpb.GetAllRestaurantItemsResponseData{
				TotalRestaurantItems: 0,
				RestaurantItems:      nil,
			},
			Message:    "Restaurant Does not exist",
			StatusCode: int64(404),
			Error:      "Bad Request",
		}, nil
	}
	var restaurantItems []model.RestaurantItem
	err := restaurantItemDBConnector.Where("restaurant_id = ?", restaurant.ID).Find(&restaurantItems)
	if err.Error != nil {
		logger.Error("Failed to get restaurant items", zap.String("restaurantId", request.RestaurantId), zap.Error(err.Error))
		return &restaurantpb.GetAllRestaurantItemsResponse{
			Data: &restaurantpb.GetAllRestaurantItemsResponseData{
				TotalRestaurantItems: 0,
				RestaurantItems:      nil,
			},
			Message:    "Failed to get restaurant items",
			StatusCode: int64(500),
			Error:      "Internal Server Error",
		}, nil
	}
	restaurantItemsResponse := []*restaurantpb.RestaurantItem{}
	for _, restaurantItem := range restaurantItems {
		restaurantItemsResponse = append(restaurantItemsResponse, &restaurantpb.RestaurantItem{
			RestaurantItemId:          strconv.FormatUint(uint64(restaurantItem.ID), 10),
			RestaurantItemName:        restaurantItem.ItemName,
			RestaurantItemPrice:       restaurantItem.ItemPrice,
			RestaurantItemImageUrl:    restaurantItem.ImageUrl,
			RestaurantItemCategory:    restaurantItem.Category,
			RestaurantItemCuisineType: restaurantItem.CuisineType,
			RestaurantItemVeg:         restaurantItem.Veg,
			RestaurantName:            restaurant.Name,
		})
	}
	logger.Info("Restaurant items fetched successfully", zap.Int("totalItems", len(restaurantItems)))
	return &restaurantpb.GetAllRestaurantItemsResponse{
		Data: &restaurantpb.GetAllRestaurantItemsResponseData{
			TotalRestaurantItems: int64(len(restaurantItems)),
			RestaurantItems:      restaurantItemsResponse,
		},
		Message:    "Restaurant items fetched successfully",
		StatusCode: int64(200),
		Error:      "",
	}, nil
}
