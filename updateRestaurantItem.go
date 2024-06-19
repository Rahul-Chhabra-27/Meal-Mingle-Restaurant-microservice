package main

// import (
// 	"context"
// 	"restaurant-micro/config"
// 	"restaurant-micro/model"
// 	restaurantpb "restaurant-micro/proto/restaurant"

// 	"google.golang.org/grpc/codes"
// )

// func (*RestaurantService) UpdateRestaurantItem(ctx context.Context, request *restaurantpb.UpdateRestaurantItemRequest) (*restaurantpb.UpdateRestaurantItemResponse, error) {
// 	// get the user email from the context
// 	userEmail, ok := ctx.Value("userEmail").(string)
// 	if !ok {
// 		return &restaurantpb.UpdateRestaurantItemResponse{
// 			Message:    "",
// 			StatusCode: int64(codes.Internal),
// 			Error:      "Internal Server Error",
// 		}, nil
// 	}
// 	// fetch restaurant from restaurantDB
// 	var restaurant model.Restaurant
// 	primaryKeyRes := restaurantDBConnector.Where("name = ?", request.RestaurantName).First(&restaurant)
// 	// check if the restaurant is exist or nor
// 	if primaryKeyRes.Error != nil || restaurant.RestaurantOwnerMail != userEmail {
// 		return &restaurantpb.UpdateRestaurantItemResponse{
// 			Message:    "",
// 			StatusCode: int64(codes.NotFound),
// 			Error:      "Restaurant Does not exist OR you are not the owner of this restaurant",
// 		}, nil
// 	}

// 	var restaurantItem model.RestaurantItem
// 	primaryKey := restaurantItemDBConnector.Where("item_name = ? AND restaurant_id = ?", request.RestaurantItemName, restaurant.ID).First(&restaurantItem)
// 	if primaryKey.Error != nil {
// 		return &restaurantpb.UpdateRestaurantItemResponse{
// 			Message:    "",
// 			StatusCode: int64(codes.Internal),
// 			Error:      "Restaurant item does not exist",
// 		}, nil
// 	}
// 	if !config.ValidateRestaurantItemFields(request.RestaurantItemName, request.RestaurantItemImageUrl) {
// 		return &restaurantpb.UpdateRestaurantItemResponse{
// 			Message:    "",
// 			StatusCode: int64(codes.InvalidArgument),
// 			Error:      "Invalid restaurant item fields",
// 		}, nil
// 	}
// 	restaurantItem.ItemName = request.RestaurantItemName
// 	restaurantItem.ItemPrice = request.RestaurantItemPrice
// 	restaurantItem.ImageUrl = request.RestaurantItemImageUrl
// 	err := restaurantItemDBConnector.Save(&restaurantItem)
// 	if err.Error != nil {
// 		return &restaurantpb.UpdateRestaurantItemResponse{
// 			Message:    "",
// 			StatusCode: 500,
// 			Error:      "Failed to update restaurant item",
// 		}, nil
// 	}
// 	return &restaurantpb.UpdateRestaurantItemResponse{
// 		Message:    "Restaurant item updated successfully",
// 		StatusCode: 200,
// 		Error:      "",
// 	}, nil
// }
