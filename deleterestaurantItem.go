package main

import (
	"context"
	"fmt"
	"restaurant-micro/model"
	restaurantpb "restaurant-micro/proto/restaurant"
	"strconv"
	"strings"
)

func (*RestaurantService) DeleteRestaurantItem(ctx context.Context, request *restaurantpb.DeleteRestaurantItemRequest) (*restaurantpb.DeleteRestaurantItemResponse, error) {
	// Get the user email from the context
	userEmail, emailCtxError := ctx.Value("userEmail").(string)
	userRole, roleCtxError := ctx.Value("userRole").(string)

	if !emailCtxError || !roleCtxError {
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.DeleteRestaurantItemResponse{
			Message:    "Failed to get user email and role from context",
			Error:      "Internal Server Error",
			StatusCode: int64(500),
		}, nil
	}

	if userRole != model.AdminRole {
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "You do not have permission to perform this action. Only admin can add a delete restaurant item",
			StatusCode: 403,
			Error:      "Forbidden",
		}, nil
	}

	// validate the restaurant item fields
	if request.RestaurantName == "" || request.RestaurantItemName == "" {
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Invalid restaurant item data provided",
			StatusCode: 400,
			Error:      "Bad Request",
		}, nil
	}
	restaurantName := strings.ReplaceAll(request.RestaurantName, "-", " ")
	restaurantItemName := strings.ReplaceAll(request.RestaurantItemName, "-", " ")

	// check if user own's this restaurant
	var restaurant model.Restaurant
	primaryKeyRes := restaurantDBConnector.Where("name = ?", restaurantName).First(&restaurant)
	if primaryKeyRes.Error != nil || restaurant.RestaurantOwnerMail != userEmail {
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Restaurant does not exist OR you are not the owner of this restaurant",
			StatusCode: 404,
			Error:      "Resource not found or forbidden",
		}, nil
	}

	// delete the restaurant item
	var restaurantItem model.RestaurantItem
	primaryKey := restaurantItemDBConnector.Where("item_name = ? AND restaurant_id = ?", restaurantItemName, restaurant.ID).First(&restaurantItem)
	if primaryKey.Error != nil {
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Restaurant item does not exist",
			StatusCode: 404,
			Error:      "Resource not found or forbidden",
		}, nil
	}
	err := restaurantItemDBConnector.Delete(&restaurantItem)
	if err.Error != nil {
		return &restaurantpb.DeleteRestaurantItemResponse{
			Data:       nil,
			Message:    "Failed to delete restaurant item",
			StatusCode: 500,
			Error:      "Internal Server Error",
		}, nil
	}
	return &restaurantpb.DeleteRestaurantItemResponse{
		Data: &restaurantpb.DeleteRestaurantItemResponseData{
			RestaurantItemId: strconv.FormatUint(uint64(restaurantItem.ID), 10),
		},
		Message:    "Restaurant item deleted successfully",
		StatusCode: 200,
		Error:      "",
	}, nil
}
