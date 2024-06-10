package service

import (
	"context"
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
			"_id":     brand.Id,
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

func UpdateBrandByAdmin(c *gin.Context) {
	// Parse the input JSON into a map
	updateFields := make(map[string]interface{})
	err := c.ShouldBindJSON(&updateFields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	// Get the brand ID from the URL parameter
	brandId := c.Param("id")
	intID, err := strconv.Atoi(brandId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error0": err.Error()})
		return
	}

	// Ensure the _id field is not included in the updateFields map
	delete(updateFields, "_id")

	// Add the updatedAt field to the updateFields map
	updateFields["updatedAt"] = time.Now()

	// Create filter to match integer _id
	filter := bson.M{"_id": intID}
	updateQuery := bson.M{"$set": updateFields}

	// Perform the update operation
	_, err = brandCollection.UpdateOne(c, filter, updateQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "brand_updated", "body": gin.H{}})
}
