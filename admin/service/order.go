package service

import (
	"math"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ordersCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

func GetOrdersByAdmin(c *gin.Context) {
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	_, ok := tokenClaims.(*auth.SignedAdminDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	limit := 20

	// Calculate offset
	offset := int64((page - 1) * limit)
	limit64 := int64(limit)

	// // Query MongoDB with pagination
	// var orders []entities.Order
	cur, err := ordersCollection.Find(c, bson.M{}, &options.FindOptions{
		Limit: &limit64,
		Skip:  &offset,
		Sort:  bson.M{"createdAt": -1},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	defer cur.Close(c)
	// Initialize a slice to hold simplified orders
	var simplifiedOrders []map[string]interface{}

	// Iterate over each order and construct simplified order object
	for cur.Next(c) {
		var order map[string]interface{}

		if err := cur.Decode(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}

		// Construct simplified order object
		// simplified := map[string]interface{}{
		// 	"isCoupon":      order.IsCoupon,
		// 	"startDate":     order.StartDate,
		// 	"status":        order.Status,
		// 	"paymentStatus": order.PaymentStatus,
		// 	"message":       order.Message,
		// 	"_id":           order.Id,
		// 	"totalPrice":    order.TotalPrice,
		// 	"totalDiscount": order.TotalDiscount,
		// 	"totalQuantity": order.TotalQuantity,
		// 	"postalCost":    order.PostalCost,
		// 	// Convert ObjectID to string
		// 	"userId":             order.UserId,
		// 	"products":           order.Products,
		// 	"jStartDate":         order.JStartDate,
		// 	"address":            order.Address,
		// 	"createdAt":          order.CreatedAt,
		// 	"updatedAt":          order.UpdatedAt,
		// 	"__v":                order.V,
		// 	"paymentId":          order.PaymentId,
		// 	"postalTrackingCode": order.PostalTrakingCode,
		// 	"mix":                order.Mix,
		// }
		// if len(order.Mix) > 0 {
		// 	simplified["mix"] = order.Mix
		// } else {
		// 	simplified["mix"] = []primitive.ObjectID{} // Empty slice
		// }

		simplifiedOrders = append(simplifiedOrders, order)
	}

	// Count total documents for pagination info
	totalDocs, err := ordersCollection.CountDocuments(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))

	// Prepare pagination response
	hasNextPage := page < totalPages
	nextPage := page + 1
	hasPrevPage := page > 1
	prevPage := page - 1

	// Response
	response := gin.H{
		"docs":          simplifiedOrders,
		"totalDocs":     totalDocs,
		"limit":         limit,
		"totalPages":    totalPages,
		"page":          page,
		"pagingCounter": offset + 1,
		"hasPrevPage":   hasPrevPage,
		"hasNextPage":   hasNextPage,
		"prevPage":      prevPage,
		"nextPage":      nextPage,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "orders", "body": response})
}

type product struct {
	Quantity        int                   `json:"quantity" bson:"quantity"`
	VariationKey    []entities.Properties `json:"variationsKey" bson:"variationsKey"`
	Id              primitive.ObjectID    `json:"_id" bson:"_id"`
	Name            string                `json:"name" bson:"name"`
	Price           int                   `json:"price" bson:"price"`
	DiscountPercent float64               `json:"discountPercent" bson:"discountPercent"`
	ProductId       primitive.ObjectID    `json:"productId" bson:"productId"`
}

func GetOrderByIdByAdmin(c *gin.Context) {
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	_, ok := tokenClaims.(*auth.SignedAdminDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	// Get the 'id' parameter from the query string
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'id' query parameter"})
		return
	}

	// Parse the id to an integer
	objectID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	filter := bson.M{"_id": objectID}
	var order entities.Order
	err = ordersCollection.FindOne(c, filter).Decode(&order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	userID, err := primitive.ObjectIDFromHex(order.UserId.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var user entities.Users
	err = userCollection.FindOne(c, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
		return
	}
	var products []product
	for _, orderProduct := range order.Products {
		var variationDetails []entities.Properties
		for _, key := range orderProduct.VariationKey {
			var variationDetail entities.Properties
			err = propertiesCollection.FindOne(c, bson.M{"_id": key}).Decode(&variationDetail)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error3": err.Error()})
				return
			}
			variationDetails = append(variationDetails, variationDetail)
		}

		products = append(products, product{
			Id:              orderProduct.Id,
			Quantity:        orderProduct.Quantity,
			VariationKey:    variationDetails,
			Name:            orderProduct.Name,
			Price:           orderProduct.Price,
			DiscountPercent: orderProduct.DiscountPercent,
			ProductId:       orderProduct.ProductId,
		})

	}

	// Construct the response
	response := gin.H{
		"_id":           order.Id,
		"startDate":     order.StartDate,
		"jStartDate":    order.JStartDate,
		"status":        order.Status,
		"paymentStatus": order.PaymentStatus,
		"totalPrice":    order.TotalPrice,
		"totalDiscount": order.TotalDiscount,
		"message":       order.Message,
		"totalQuantity": order.TotalQuantity,
		"postalCost":    order.PostalCost,
		"user": gin.H{
			"_id":         user.Id.Hex(),
			"phoneNumber": user.PhoneNumber,
			"lastname":    user.LastName,
			"name":        user.Name,
		},
		"products":           products,
		"address":            order.Address,
		"postalTrackingCode": order.PostalTrakingCode,
		// "mix":                order.Mix,
		"mixproducts": []primitive.ObjectID{},
	}
	if len(order.Mix) > 0 {
		response["mix"] = order.Mix
	} else {
		response["mix"] = []primitive.ObjectID{}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "order", "body": response})
}
func GetInformationUserByAdmin(c *gin.Context) {
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	_, ok := tokenClaims.(*auth.SignedAdminDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'id' query parameter"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	filter := bson.M{"_id": objectID}
	var user entities.Users
	err = userCollection.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	// var orders []gin.H
	var orders []map[string]interface{}
	cursor, err := ordersCollection.Find(c, bson.M{"userId": objectID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(c)
	for cursor.Next(c) {
		var order map[string]interface{}
		// var order entities.Order
		if err = cursor.Decode(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}
		orderMap := map[string]interface{}{
			"_id":           order["_id"],
			"totalPrice":    order["totalPrice"],
			"totalDiscount": order["totalDiscount"],
			"address": gin.H{
				"_id":        order["address"].(map[string]interface{})["_id"],
				"address":    order["address"].(map[string]interface{})["address"],
				"city":       order["address"].(map[string]interface{})["city"],
				"postalCode": order["address"].(map[string]interface{})["postalCode"],
				"state":      order["address"].(map[string]interface{})["state"],
			},
			"createdAt": order["createdAt"],
			"paymentId": order["paymentId"],
		}
		orders = append(orders, orderMap)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
		return
	}

	response := gin.H{
		"user": gin.H{
			"_id":                  user.Id.Hex(),
			"activeSession":        user.ActiveSession,
			"fcmRegistrationToken": user.FcmRegistratinToken,
			"favoritesProducts":    user.Favorites,
			"username":             user.Username,
			"phoneNumber":          user.PhoneNumber,
			"sex":                  user.Sex,
			"role":                 user.Role,
			"addresses":            user.Address,
			"createdAt":            user.CreatedAt,
		},
		"orders": orders,
	}

	c.JSON(http.StatusOK, response)
}
