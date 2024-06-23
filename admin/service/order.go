package service

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"shop/database"
	"shop/entities"
	"shop/helpers"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ordersCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

func GetOrdersByAdmin(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	limit := 15

	// Calculate offset
	offset := int64((page - 1) * limit)
	limit64 := int64(limit)

	// Get the status query parameter
	status := c.DefaultQuery("status", "")
	paymentStatus := c.DefaultQuery("paymentstatus", "")
	log.Printf("Received query params - page: %d, status: %s, paymentStatus: %s", page, status, paymentStatus)

	// Build the filter for the MongoDB query
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	if paymentStatus != "" {
		filter["paymentStatus"] = paymentStatus
	}
	log.Printf("Constructed filter: %+v", filter)
	// Query MongoDB with pagination and status filter
	cur, err := ordersCollection.Find(c, filter, &options.FindOptions{
		Limit: &limit64,
		Skip:  &offset,
		Sort:  bson.M{"createdAt": -1},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(c)

	var simplifiedOrders []map[string]interface{}

	for cur.Next(c) {
		var order bson.M
		if err := cur.Decode(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userId, ok := order["userId"].(primitive.ObjectID)
		if !ok {
			order["userId"] = nil
			simplifiedOrders = append(simplifiedOrders, order)
			continue
		}

		var user bson.M
		userFilter := bson.M{"_id": userId}
		err := userCollection.FindOne(c, userFilter).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				order["userId"] = nil
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user details", "details": err.Error()})
				return
			}
		} else {
			order["userId"] = user
		}

		simplifiedOrders = append(simplifiedOrders, order)
	}

	totalDocs, err := ordersCollection.CountDocuments(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))

	response := gin.H{
		"docs":          simplifiedOrders,
		"totalDocs":     totalDocs,
		"limit":         limit,
		"totalPages":    totalPages,
		"page":          page,
		"pagingCounter": offset + 1,
		"hasPrevPage":   page > 1,
		"hasNextPage":   page < totalPages,
		"prevPage":      page - 1,
		"nextPage":      page + 1,
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

	id := c.Param("id")
	fmt.Println("ID:", id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'id' query parameter"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
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
			// "paymentId": order["paymentId"],
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

func UpdateSatusOrderByAdmin(c *gin.Context) {
	var user entities.Users
	var orderData entities.Order

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'id' query parameter"})
		return
	}
	objectID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	filter := bson.M{"_id": objectID}

	var updateStatus struct {
		Status            string `json:"status" binding:"required"`
		PostalTrakingCode string `json:"postalTrackingCode" bson:"postalTrackingCode"`
	}
	err = c.ShouldBindJSON(&updateStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
		return
	}
	update := bson.M{
		"$set": bson.M{
			"status":             updateStatus.Status,
			"postalTrackingCode": updateStatus.PostalTrakingCode,
		},
	}
	_ = ordersCollection.FindOneAndUpdate(c, filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error3": err.Error()})
		return
	}
	// orderIdStr := strconv.Itoa(int(orderData.Id.(int32)))
	if updateStatus.Status == "pending" {
		err := ordersCollection.FindOne(c, filter).Decode(&orderData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error4": err.Error()})
			return
		}
		err = userCollection.FindOne(c, bson.M{"_id": orderData.UserId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error5": err.Error()})
			return
		}
		rawPhoneNumber := user.PhoneNumber

		phoneNumberWithZero := "0" + rawPhoneNumber
		helpers.SendSmstoCustomerForReadyProduct(strconv.Itoa(objectID), phoneNumberWithZero)

	} else if updateStatus.Status == "done" {
		err := ordersCollection.FindOne(c, filter).Decode(&orderData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error4": err.Error()})
			return
		}
		err = userCollection.FindOne(c, bson.M{"_id": orderData.UserId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error5": err.Error()})
			return
		}
		rawPhoneNumber := user.PhoneNumber

		phoneNumberWithZero := "0" + rawPhoneNumber
		helpers.SendSmstoCustomerForSendProduct(strconv.Itoa(objectID), orderData.PostalTrakingCode, phoneNumberWithZero)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "status_updated", "body": gin.H{}})

}

func GetAllOrders(c *gin.Context) {
	var orders int64
	count, err := ordersCollection.CountDocuments(c, bson.M{})
	if err != nil {
		// Handle error
		return
	}
	orders = count

	c.JSON(http.StatusOK, gin.H{"count": orders})
}
