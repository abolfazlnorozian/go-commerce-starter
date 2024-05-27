package service

import (
	"context"
	"fmt"
	"net/http"
	"shop/database"
	"shop/entities"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var brandCollection *mongo.Collection = database.GetCollection(database.DB, "brands")

//@Summary Get Brands
//@Description Get All Brands
//@Tags brands
//@Accept json
//@produce json
//@Success 201
//@Router /api/brands [get]
func GetBrandsByAdmin(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var brands []map[string]interface{} // Use slice of maps to hold specific fields

	results, err := brandCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to find collection"})
		return
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var brand entities.Brands
		err := results.Decode(&brand)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// Create a map with the desired fields from the 'entities.Brands' struct
		brandData := map[string]interface{}{
			"id":      brand.Id,
			"name":    brand.Name,
			"details": brand.Details,
			"image":   brand.Image,
		}

		brands = append(brands, brandData)
	}

	c.JSON(http.StatusCreated, gin.H{"body": brands, "message": "brands", "success": true})
}
func getNextID() (int, error) {
	var result struct {
		MaxID int `bson:"_id"`
	}

	// Find the maximum integer ID currently in use in the collection
	err := brandCollection.FindOne(context.Background(), bson.D{}, options.FindOne().SetSort(bson.D{{"_id", -1}})).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return 1, nil
		}

		return 0, err
	}

	return result.MaxID + 1, nil
}

func PostBrandsByAdmin(c *gin.Context) {

	var brand entities.Brands
	err := c.ShouldBindJSON(&brand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	id, err := getNextID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}
	brand.Id = id
	brand.CreatedAt = time.Now()
	brand.UpdatedAt = time.Now()

	_, err = brandCollection.InsertOne(c, bson.M{
		"_id":       id,
		"name":      brand.Name,
		"details":   brand.Details,
		"image":     brand.Image,
		"createdAt": brand.CreatedAt,
		"updatedAt": brand.UpdatedAt,
		"__v":       brand.V,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "brand_added", "body": gin.H{}})

}

func DeleteBrandsByAdmin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid 'id' parameter"})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Create filter to match integer _id
	filter := bson.M{"_id": intID}

	var brand entities.Brands2
	err = brandCollection.FindOneAndDelete(c, filter).Decode(&brand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "brand_updated", "body": gin.H{}})

}

// func UpdateBrandByAdmin(c *gin.Context) {
// 	// Extract the id from the URL parameter
// 	id := c.Param("id")
// 	fmt.Println("Extracted id from URL:", id)

// 	// Convert the id to an integer
// 	intID, err := strconv.Atoi(id)
// 	if err != nil {
// 		fmt.Println("Error converting id to integer:", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter, must be an integer"})
// 		return
// 	}
// 	fmt.Println("Converted id to integer:", intID)

// 	// Create a filter to match the integer _id
// 	filter := bson.M{"_id": intID}

// 	// Parse the incoming JSON payload into a struct
// 	var updateBrand struct {
// 		Name      string    `json:"name"`
// 		Details   string    `json:"details"`
// 		UpdatedAt time.Time `json:"updatedAt"`
// 	}
// 	if err := c.BindJSON(&updateBrand); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
// 		return
// 	}

// 	// Set the updatedAt field to the current time
// 	updateBrand.UpdatedAt = time.Now()

// 	// Prepare the update query
// 	updateQuery := bson.M{"$set": updateBrand}

// 	// Execute the update operation
// 	_, err = brandCollection.UpdateOne(c, filter, updateQuery)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Return a success response
// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "brand_updated", "body": gin.H{}})
// }

func UpdateBrandByAdmin(c *gin.Context) {

	// var brand entities.Brands
	updateFields := make(map[string]interface{})
	err := c.ShouldBindJSON(&updateFields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	barnId := c.Param("_id")
	fmt.Println("a", barnId)

	intID, err := strconv.Atoi(barnId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error0": err.Error()})
		return
	}
	fmt.Println("intId", intID)
	// Create filter to match integer _id
	filter := bson.M{"_id": intID}
	updateFields["updatedAt"] = time.Now()
	// brand.CreatedAt = time.Now()
	// brand.UpdatedAt = time.Now()
	updateQuery := bson.M{"$set": updateFields}
	_, err = brandCollection.UpdateOne(c, filter, updateQuery)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	// _, err = brandCollection.UpdateOne(c, filter, bson.M{
	// 	"$set": bson.M{
	// 		// "_id":       brand.Id,
	// 		"name":      brand.Name,
	// 		"details":   brand.Details,
	// 		"image":     brand.Image,
	// 		"createdAt": brand.CreatedAt,
	// 		"updatedAt": brand.UpdatedAt,
	// 	},
	// })
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
	// 	return
	// }
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "brand_updated", "body": gin.H{}})

}
