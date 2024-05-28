package service

import (
	"math"
	"net/http"
	"shop/database"
	"shop/entities"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var couponCollection *mongo.Collection = database.GetCollection(database.DB, "coupons")

func GetCouponsByAdmin(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	limit := 15

	// Calculate offset
	offset := int64((page - 1) * limit)
	limit64 := int64(limit)

	// // Query MongoDB with pagination
	// var orders []entities.Order
	cur, err := couponCollection.Find(c, bson.M{}, &options.FindOptions{
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
	var coupons []entities.Coupons

	// Iterate over each order and construct simplified order object
	for cur.Next(c) {
		var coupon entities.Coupons

		if err := cur.Decode(&coupon); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}
		coupons = append(coupons, coupon)

	}

	// Count total documents for pagination info
	totalDocs, err := couponCollection.CountDocuments(c, bson.M{})
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
		"docs":          coupons,
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
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "coupons", "body": response})

}

func extractCouponFields(payload map[string]interface{}) (entities.Coupons, error) {
	var coupon entities.Coupons
	var err error

	// Handle ID
	if id, ok := payload["_id"].(string); ok {
		coupon.ID, _ = primitive.ObjectIDFromHex(id)
	} else {
		coupon.ID = primitive.NewObjectID()
	}

	// Handle IsExpired
	if isExpired, ok := payload["isExpired"].(bool); ok {
		coupon.IsExpired = isExpired
	}

	// Handle WithoutDiscount
	if withoutDiscount, ok := payload["withoutDiscount"].(bool); ok {
		coupon.WithoutDiscount = withoutDiscount
	}

	// Handle CouponCode
	if couponCode, ok := payload["couponCode"].(string); ok {
		coupon.CouponCode = couponCode
	}

	// Handle Amount
	if amount, ok := payload["amount"].(float64); ok {
		coupon.Amount = int(amount)
	} else if amount, ok := payload["amount"].(string); ok {
		coupon.Amount, err = strconv.Atoi(amount)
		if err != nil {
			return coupon, err
		}
	}

	// Handle MinimumPurchase
	if minimumPurchase, ok := payload["minimumPurchase"].(float64); ok {
		coupon.MinimumPurchase = int(minimumPurchase)
	} else if minimumPurchase, ok := payload["minimumPurchase"].(string); ok {
		coupon.MinimumPurchase, err = strconv.Atoi(minimumPurchase)
		if err != nil {
			return coupon, err
		}
	}

	// Handle DateExpired
	if dateExpired, ok := payload["dateExpired"].(string); ok {
		coupon.DateExpired, err = time.Parse(time.RFC3339, dateExpired)
		if err != nil {
			return coupon, err
		}
	}

	// Handle To
	if to, ok := payload["to"].(string); ok {
		coupon.To = to
	}

	// Handle User
	if user, ok := payload["user"].(string); ok {
		coupon.User = user
	}

	// Handle Category
	if category, ok := payload["category"].(string); ok {
		coupon.Category, err = primitive.ObjectIDFromHex(category)
		if err != nil {
			return coupon, err
		}
	}

	// Set UpdatedAt (and CreatedAt if creating a new coupon)
	coupon.UpdatedAt = time.Now()

	return coupon, nil
}

func UpdateCouponByAdmin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	filter := bson.M{"_id": objectID}
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract and convert the fields
	coupon, err := extractCouponFields(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only the fields that are provided in the payload
	updateFields := bson.M{
		"isExpired":       coupon.IsExpired,
		"withoutDiscount": coupon.WithoutDiscount,
		"couponCode":      coupon.CouponCode,
		"amount":          coupon.Amount,
		"minimumPurchase": coupon.MinimumPurchase,
		"dateExpired":     coupon.DateExpired,
		"to":              coupon.To,
		"user":            coupon.User,
		"category":        coupon.Category,
		"updatedAt":       coupon.UpdatedAt,
	}

	updateQuery := bson.M{"$set": updateFields}

	_, err = couponCollection.UpdateOne(c, filter, updateQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "coupon_edited", "body": gin.H{}})
}

func PostCouponByAdmin(c *gin.Context) {
	// Parse the JSON payload into a map
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract and convert the fields
	coupon := entities.Coupons{}
	var err error

	// Handle ID
	if id, ok := payload["_id"].(string); ok {
		coupon.ID, _ = primitive.ObjectIDFromHex(id)
	} else {
		coupon.ID = primitive.NewObjectID()
	}

	// Handle IsExpired
	if isExpired, ok := payload["isExpired"].(bool); ok {
		coupon.IsExpired = isExpired
	}

	// Handle WithoutDiscount
	if withoutDiscount, ok := payload["withoutDiscount"].(bool); ok {
		coupon.WithoutDiscount = withoutDiscount
	}

	// Handle CouponCode
	if couponCode, ok := payload["couponCode"].(string); ok {
		coupon.CouponCode = couponCode
	}

	// Handle Amount
	if amount, ok := payload["amount"].(string); ok {
		coupon.Amount, err = strconv.Atoi(amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
			return
		}
	}

	// Handle MinimumPurchase
	if minimumPurchase, ok := payload["minimumPurchase"].(string); ok {
		coupon.MinimumPurchase, err = strconv.Atoi(minimumPurchase)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minimumPurchase"})
			return
		}
	}

	// Handle DateExpired
	if dateExpired, ok := payload["dateExpired"].(string); ok {
		coupon.DateExpired, err = time.Parse("2006-01-02", dateExpired)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dateExpired"})
			return
		}
	}

	// Handle To
	if to, ok := payload["to"].(string); ok {
		coupon.To = to
	}

	// Handle User
	if user, ok := payload["user"].(string); ok {
		coupon.User = user
	}

	// Handle Category
	if category, ok := payload["category"].(string); ok {
		coupon.Category, err = primitive.ObjectIDFromHex(category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
			return
		}
	}

	// Set CreatedAt and UpdatedAt
	coupon.CreatedAt = time.Now()
	coupon.UpdatedAt = time.Now()

	// Insert the coupon into the collection
	_, err = couponCollection.InsertOne(c, coupon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "coupon_saved", "body": gin.H{}})
}

func DeleteCouponByAdmin(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid 'id' parameter"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}
	filter := bson.M{"_id": objectID}
	var coupon entities.Coupons
	err = couponCollection.FindOneAndDelete(c, filter).Decode(&coupon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "coupon_deleted", "body": gin.H{}})
}
