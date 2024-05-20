package service

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"shop/auth"
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

var proCollection *mongo.Collection = database.GetCollection(database.DB, "products")

// func GetAllProductsByAdmin(c *gin.Context) {
// 	// Pagination parameters from the query
// 	pageStr := c.DefaultQuery("page", "1")
// 	limitStr := c.DefaultQuery("limit", "40")

// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page < 1 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
// 		return
// 	}

// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil || limit < 1 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
// 		return
// 	}

// 	// Your existing code for fetching and processing products...
// 	// The pagination parameters (page and limit) are now valid, continue with your existing logic...
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Calculate skip value for pagination
// 	skip := (page - 1) * limit

// 	// Calculate total number of documents in the collection
// 	totalDocs, err := proCollection.CountDocuments(ctx, bson.M{})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
// 		return
// 	}

// 	// Fetch products based on the constructed filter
// 	results, err := proCollection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
// 		return
// 	}

// 	defer results.Close(ctx)
// 	var products []entities.Products

// 	for results.Next(ctx) {
// 		var pro entities.Products
// 		err := results.Decode(&pro)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode products"})
// 			return
// 		}

// 		products = append(products, pro)
// 	}

// 	// Calculate total number of pages based on the limit
// 	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))

// 	// Determine if there are previous and next pages
// 	hasPrevPage := page > 1
// 	hasNextPage := page < totalPages

// 	// Prepare the custom response with selected fields
// 	var customProducts []gin.H
// 	for _, product := range products {
// 		customProduct := gin.H{
// 			"_id":             product.ID,
// 			"notExist":        product.NotExist,
// 			"amazing":         product.Amazing,
// 			"productType":     product.ProductType,
// 			"images":          product.Images,
// 			"name":            product.Name,
// 			"price":           product.Price,
// 			"discountPercent": product.DiscountPercent,
// 			"stock":           product.Stock,
// 			"slug":            product.Slug,
// 			"variations":      product.Variations,
// 			"salesNumber":     product.SalesNumber,
// 			"bannerUrl":       product.BannerUrl,
// 		}
// 		customProducts = append(customProducts, customProduct)
// 	}

// 	// Prepare the response with custom products and pagination information
// 	response := gin.H{
// 		"docs":          customProducts,
// 		"totalDocs":     totalDocs,
// 		"limit":         limit,
// 		"totalPages":    totalPages,
// 		"page":          page,
// 		"pagingCounter": skip + 1,
// 		"hasPrevPage":   hasPrevPage,
// 		"hasNextPage":   hasNextPage,
// 	}

// 	// Set prevPage and nextPage values based on the current page
// 	if hasPrevPage {
// 		response["prevPage"] = page - 1
// 	} else {
// 		response["prevPage"] = nil
// 	}

// 	if hasNextPage {
// 		response["nextPage"] = page + 1
// 	} else {
// 		response["nextPage"] = nil
// 	}

// 	c.JSON(http.StatusOK, response)
// }

func GetAllProductsByAdmin(c *gin.Context) {
	// pages := c.Query("page")

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
	filter := bson.M{}
	// Pagination parameters from the query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "40"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Calculate total number of documents in the collection
	totalDocs, err := proCollection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count products"})
		return
	}

	// Fetch products based on the constructed filter
	results, err := proCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query products"})
		return
	}

	defer results.Close(ctx)
	var products []entities.Products

	for results.Next(ctx) {
		var pro entities.Products
		err := results.Decode(&pro)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message10": err.Error()})
			return
		}

		products = append(products, pro)
	}

	// Calculate total number of pages based on the limit
	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))

	// Determine if there are previous and next pages
	hasPrevPage := page > 1
	hasNextPage := page < totalPages

	// Prepare the custom response with selected fields
	var customProducts []gin.H
	for _, product := range products {
		customProduct := gin.H{
			"_id":             product.ID,
			"notExist":        product.NotExist,
			"amazing":         product.Amazing,
			"productType":     product.ProductType,
			"images":          product.Images,
			"name":            product.Name,
			"price":           product.Price,
			"discountPercent": product.DiscountPercent,
			"stock":           product.Stock,
			"slug":            product.Slug,
			"variations":      product.Variations,
			"salesNumber":     product.SalesNumber,
			"bannerUrl":       product.BannerUrl,
		}
		customProducts = append(customProducts, customProduct)
	}

	// Prepare the response with custom products and pagination information
	response := gin.H{
		"docs":          customProducts,
		"totalDocs":     totalDocs,
		"limit":         limit,
		"totalPages":    totalPages,
		"page":          page,
		"pagingCounter": skip + 1,
		"hasPrevPage":   hasPrevPage,
		"hasNextPage":   hasNextPage,
	}

	// Set prevPage and nextPage values based on the current page
	if hasPrevPage {
		response["prevPage"] = page - 1
	} else {
		response["prevPage"] = nil
	}

	if hasNextPage {
		response["nextPage"] = page + 1
	} else {
		response["nextPage"] = nil
	}

	c.JSON(http.StatusOK, response)

}
func GetProductByIdByAdmin(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid 'id' parameter"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Define a filter to find the product based on productId
	filter := bson.M{"_id": objectID}
	var product entities.Products
	err = proCollection.FindOne(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	var categoryInfo []map[string]interface{}
	for _, catID := range product.Category {
		var category entities.Category
		err := categoryCollection.FindOne(c, bson.M{"_id": catID}).Decode(&category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}
		categoryMap := map[string]interface{}{
			"_id":       category.ID,
			"image":     category.Images,
			"parent":    category.Parent,
			"name":      category.Name,
			"ancestors": category.Ancestors,
			"slug":      category.Slug,
			"__v":       category.V,
			"details":   category.Details,
			"faq":       category.Faq,
		}
		categoryInfo = append(categoryInfo, categoryMap)
	}
	productData := gin.H{
		"notExist":        product.NotExist,
		"amazing":         product.Amazing,
		"productType":     product.ProductType,
		"quantity":        product.Quantity,
		"comments":        product.Comment,
		"parent":          product.Parent,
		"categories":      categoryInfo,
		"tags":            product.Tags,
		"similarProducts": product.SimilarProducts,
		"_id":             product.ID,
		"images":          product.Images,
		"name":            product.Name,
		"price":           product.Price,
		"details":         product.Details,
		"discountPercent": product.DiscountPercent,
		"bannerUrl":       product.BannerUrl,
		"stock":           product.Stock,
		"categoryId":      product.CategoryID,
		"attributes":      product.Attributes,
		"slug":            product.Slug,
		"shortId":         product.ShortID,
		"dimensions":      product.Dimensions,
		"variations":      product.Variations,
		"createdAt":       product.CreatedAt,
		"updatedAt":       product.UpdatedAt,
		"__v":             product.V,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "product", "body": productData})
}

func PostProductByAdmin(c *gin.Context) {
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
	var product entities.Products
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	var newAttributes []entities.Attribute
	err = c.ShouldBindJSON(&newAttributes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range newAttributes {
		newAttributes[i].ID = primitive.NewObjectID()
	}
	product.Attributes = append(product.Attributes, newAttributes...)

	_, err = proCollection.InsertOne(c, bson.M{
		"_id":             product.ID,
		"notExist":        product.NotExist,
		"amazing":         product.Amazing,
		"isMillModel":     product.IsMillModel,
		"quantity":        product.Quantity,
		"comments":        product.Comment,
		"parent":          product.Parent,
		"categories":      product.Category,
		"tags":            product.Tags,
		"similarProducts": product.SimilarProducts,
		"name_fuzzy":      product.NameFuzzy,
		"productType":     product.ProductType,
		"images":          product.Images,
		"name":            product.Name,
		"price":           product.Price,
		"discountPercent": product.DiscountPercent,
		"details":         product.Details,
		"categoryId":      product.CategoryID,
		"attributes":      product.Attributes,
		"dimensions":      product.Dimensions,
		"stock":           product.Stock,
		"slug":            product.Slug,
		"variations":      product.Variations,
		"createdAt":       product.CreatedAt,
		"updatedAt":       product.UpdatedAt,
		"salesNumber":     product.SalesNumber,
		"bannerUrl":       product.BannerUrl,
		"__v":             product.V,
		"shortId":         product.ShortID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok_added", "body": gin.H{}})

}
func DeleteProductByAdmin(c *gin.Context) {
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
	filter := bson.M{"_id": objectID}
	var product entities.Products
	err = proCollection.FindOneAndDelete(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok_delete", "body": gin.H{}})

}
func UpdateProductByAdmin(c *gin.Context) {
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

	filter := bson.M{"_id": objectID}
	updateFields := make(map[string]interface{})

	// Extract fields to update from request JSON
	if err := c.ShouldBindJSON(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the updatedAt field to the current time
	updateFields["updatedAt"] = time.Now()

	// Construct the update query
	updateQuery := bson.M{"$set": updateFields}

	// Execute the update query
	_, err = proCollection.UpdateOne(c, filter, updateQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok_updated", "body": gin.H{}})
}

//*************************POST DIMENTION******************************************

func PostDimensionByAdmin(c *gin.Context) {
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

	// Define the filter to find the product by its ID
	filter := bson.M{"_id": objectID}

	// Bind the request body to the DimensionKey struct
	var dimensionKey struct {
		Key int `json:"dimensionKey"`
	}
	if err := c.ShouldBindJSON(&dimensionKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the product by ID
	var product entities.Products
	err = proCollection.FindOne(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Find the index of the existing dimension with the given key
	var existingIndex int = -1
	for i, dim := range product.Dimensions {
		if dim.Key == dimensionKey.Key {
			existingIndex = i
			break
		}
	}

	// If existing dimension found, update its values
	if existingIndex != -1 {
		product.Dimensions[existingIndex].Values = []int{} // Empty values array
	} else {
		// Add a new dimension with the provided key and empty values array
		product.Dimensions = append(product.Dimensions, entities.Dimension{
			ID:     primitive.NewObjectID(),
			Key:    dimensionKey.Key,
			Values: []int{},
		})
	}
	update := bson.M{
		"$set": bson.M{
			"dimensions": product.Dimensions,
		},
	}
	// Update the product in the database
	_, err = proCollection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "dimension_added", "body": gin.H{}})
}

//********************************************************************************8
func PostValuesByAdminToDimension(c *gin.Context) {
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

	// Define the filter to find the product by its ID
	filter := bson.M{"_id": objectID}

	// Find the product by ID
	var product entities.Products
	err = proCollection.FindOne(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Extract the dimension key from the URL parameter
	dimensionKey := c.Param("dimensionKey")
	if dimensionKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'dimensionKey' parameter"})
		return
	}
	key, err := strconv.Atoi(dimensionKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'dimensionKey' parameter"})
		return
	}

	// Find the index of the dimension with the provided key
	var dimensionIndex int = -1
	for i, dim := range product.Dimensions {
		if dim.Key == key {
			dimensionIndex = i
			break
		}
	}

	if dimensionIndex == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Dimension with key %d not found", key)})
		return
	}

	// Parse the JSON payload to get the values to be added
	var payload struct {
		Values []int `json:"values"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Overwrite the existing values with the new ones
	product.Dimensions[dimensionIndex].Values = payload.Values
	update := bson.M{
		"$set": bson.M{
			"dimensions": product.Dimensions,
		},
	}
	// Update the product in the database
	_, err = proCollection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "property_updated", "body": gin.H{}})
}

//***********************************POST  VARIATIONA*************************************************************
func PostVariationToProductByAdmin(c *gin.Context) {
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

	var variation entities.Variation
	err = c.ShouldBindJSON(&variation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	var product entities.Products
	filter := bson.M{"_id": objectID}
	err = proCollection.FindOne(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	product.Variations = append(product.Variations, entities.Variation{
		Id:              primitive.NewObjectID(),
		DiscountPercent: variation.DiscountPercent,
		Quantity:        variation.Quantity,
		Keys:            variation.Keys,
		Price:           variation.Price,
	})
	update := bson.M{
		"$set": bson.M{
			"variations": product.Variations,
		},
	}
	_, err = proCollection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "varaition_added", "body": gin.H{}})

}

type DimensionResponse struct {
	Values []int               `json:"values"`
	ID     string              `json:"_id,omitempty"`
	Key    *PropertiesResponse `json:"key"`
}
type VariationResponae struct {
	Keys []PropertiesResponse `json:"keys" bson:"keys"`
	// DiscountPercent int                   `json:"discountPercent" bson:"discountPercent"`
	// Quantity        int                   `json:"quantity" bson:"quantity"`
	// Id              primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	// Price           int                   `json:"price" bson:"price"`
}

type PropertiesResponse struct {
	ID        int       `json:"_id" bson:"_id"`
	Parent    *int      `json:"parent" bson:"parent"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	V         int       `json:"__v" bson:"__v"`
}

func GetProductByIDByVAlueFalse(c *gin.Context) {
	// Check token claims
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Token claims not found in context"})
		return
	}
	_, ok := tokenClaims.(*auth.SignedAdminDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Invalid token claims type"})
		return
	}

	// Check 'values' parameter
	values := c.Query("values")
	if values != "false" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid 'values' parameter"})
		return
	}

	// Get product ID from request parameters
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid 'id' parameter"})
		return
	}

	// Convert ID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid 'id' parameter"})
		return
	}

	// Fetch product from the database
	var product entities.Products
	filter := bson.M{"_id": objectID}
	err = proCollection.FindOne(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to fetch product"})
		return
	}

	// Prepare response dimensions
	var dimensionResponses []DimensionResponse
	for _, dim := range product.Dimensions {
		// Fetch properties document from the properties collection
		var properties entities.Properties
		err = propertiesCollection.FindOne(c, bson.M{"_id": dim.Key}).Decode(&properties)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch dimension key information"})
			return
		}

		// Create a DimensionResponse using the fetched properties
		dimensionResponse := DimensionResponse{
			Values: dim.Values,
			ID:     dim.ID.Hex(),
			Key: &PropertiesResponse{
				ID:        properties.ID,
				Parent:    properties.Parent,
				Name:      properties.Name,
				CreatedAt: properties.CreatedAt,
				UpdatedAt: properties.UpdatedAt,
				V:         properties.V,
			},
		}
		dimensionResponses = append(dimensionResponses, dimensionResponse)
	}
	var variationsResponses []VariationResponae
	for _, variat := range product.Variations {
		// Fetch properties documents from the properties collection
		cursor, err := propertiesCollection.Find(c, bson.M{"_id": bson.M{"$in": variat.Keys}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch variation keys information"})
			return
		}
		defer cursor.Close(c)

		var properties []*entities.Properties
		for cursor.Next(c) {
			var prop entities.Properties
			if err := cursor.Decode(&prop); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to decode properties document"})
				return
			}
			properties = append(properties, &prop)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Cursor error"})
			return
		}

		// Create VariationResponse using the fetched properties
		var variationKeys []PropertiesResponse
		for _, prop := range properties {
			variationKey := PropertiesResponse{
				ID:        prop.ID,
				Parent:    prop.Parent,
				Name:      prop.Name,
				CreatedAt: prop.CreatedAt,
				UpdatedAt: prop.UpdatedAt,
				V:         prop.V,
			}
			variationKeys = append(variationKeys, variationKey)
		}
		variationResponse := VariationResponae{
			Keys: variationKeys,
		}
		variationsResponses = append(variationsResponses, variationResponse)
	}
	if variationsResponses == nil {
		variationsResponses = make([]VariationResponae, 0)
	}
	// Construct the response JSON similar to the provided format
	responseBody := gin.H{
		"_id":             product.ID.Hex(),
		"dimensions":      dimensionResponses,
		"name":            product.Name,
		"productType":     product.ProductType,
		"price":           product.Price,
		"variations":      product.Variations,
		"categories":      product.Category,
		"categoryId":      product.CategoryID,
		"attributes":      product.Attributes,
		"tags":            product.Tags,
		"bannerUrl":       product.BannerUrl,
		"slug":            product.Slug,
		"discountPercent": product.DiscountPercent,
		"amazing":         product.Amazing,
		"images":          product.Images,
		"variationNames":  variationsResponses,
		"shortId":         product.ShortID,
	}

	// Send the response
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "product", "body": responseBody})
}

///*******************************DELETE VARIATION*************************************************
func DeleteVariationFromProductByAdmin(c *gin.Context) {
	// Check token claims
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

	// Get product ID from URL parameter
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Convert ID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Define filter to find product by ID
	filter := bson.M{"_id": objectID}

	// Find the product by ID
	var product entities.Products
	err = proCollection.FindOne(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get variation ID from URL parameter
	variatID := c.Param("variatId")
	if variatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'variatId' parameter"})
		return
	}

	// Convert variation ID string to ObjectID
	variationObjectID, err := primitive.ObjectIDFromHex(variatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'variatId' parameter"})
		return
	}

	// Loop through variations to find and remove the one with matching ID
	var updatedVariations []entities.Variation
	for _, variation := range product.Variations {
		if variation.Id != variationObjectID {
			updatedVariations = append(updatedVariations, variation)
		}
	}

	product.Variations = updatedVariations
	var update bson.M
	if len(updatedVariations) == 0 {

		update = bson.M{"$set": bson.M{"variations": []entities.Variation{}}}
	} else {

		update = bson.M{"$set": bson.M{"variations": updatedVariations}}
	}

	_, err = proCollection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "varaition_deleted", "body": gin.H{}})
}

//***********************DELETE DIMENSION KEY FORM DIMENSION*************************************
func DeleteDimensionFromProductByAdmin(c *gin.Context) {
	// Check token claims
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

	// Get product ID from URL parameter
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Convert ID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Define filter to find product by ID
	filter := bson.M{"_id": objectID}

	// Find the product by ID
	var product entities.Products
	err = proCollection.FindOne(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	KeyID := c.Param("keyID")
	if KeyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'variatId' parameter"})
		return
	}

	intID, err := strconv.Atoi(KeyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Loop through variations to find and remove the one with matching ID
	var updatedDimesion []entities.Dimension
	for _, dimension := range product.Dimensions {
		if dimension.Key != intID {
			updatedDimesion = append(updatedDimesion, dimension)
		}
	}

	product.Dimensions = updatedDimesion
	var update bson.M
	if len(updatedDimesion) == 0 {

		update = bson.M{"$set": bson.M{"dimensions": []entities.Dimension{}}}
	} else {

		update = bson.M{"$set": bson.M{"dimensions": updatedDimesion}}
	}

	_, err = proCollection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "dimension_deleted", "body": gin.H{}})
}

//**************************************************************************************************

var mixProductCollection *mongo.Collection = database.GetCollection(database.DB, "mixproducts")

func GetMixProductsByAdmin(c *gin.Context) {

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

	var mixProducts []entities.MixProducts
	cur, err := mixProductCollection.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": "Not Find Collection"})
		return
	}
	defer cur.Close(c)
	for cur.Next(c) {
		var mix entities.MixProducts
		err := cur.Decode(&mix)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}
		mixProducts = append(mixProducts, mix)

	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "mix_products", "body": mixProducts})

}

func PostMixByAdmin(c *gin.Context) {
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

	var mix entities.MixProducts
	err := c.ShouldBindJSON(&mix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	mix.ID = primitive.NewObjectID()
	mix.CreatedAt = time.Now()
	mix.UpdatedAt = time.Now()
	_, err = mixProductCollection.InsertOne(c, bson.M{
		"_id":       mix.ID,
		"name":      mix.Name,
		"image":     mix.Images,
		"price":     mix.Price,
		"createdAt": mix.CreatedAt,
		"updatedAt": mix.UpdatedAt,
		"__v":       mix.V,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "mix_product"})
}

func DeleteMixBYAdmin(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid 'id' parameter"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}
	filter := bson.M{"_id": objectID}
	var mix entities.MixProducts
	err = mixProductCollection.FindOneAndDelete(c, filter).Decode(&mix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "mix_product"})

}
