package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"shop/database"
	"shop/entities"
	"strconv"
	"strings"
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
func GetProducts(c *gin.Context) {
	searchQuery := c.DefaultQuery("search", "")
	if searchQuery != "" {
		GetProductBySearch(c)
	} else {
		GetAllProductsByAdmin(c)
	}

}
func GetProductBySearch(c *gin.Context) {
	filter := bson.M{}
	projection := bson.M{"name_fuzzy": 0}

	searchQuery := c.DefaultQuery("search", "")
	if searchQuery != "" {

		keywords := strings.Fields(searchQuery)

		regexPatterns := make([]bson.M, len(keywords))
		for i, keyword := range keywords {
			regexPatterns[i] = bson.M{"name_fuzzy": bson.M{"$regex": keyword, "$options": "i"}}
		}

		filter["$and"] = regexPatterns
	}

	var products []entities.Products
	cur, err := proCollection.Find(c, filter, options.Find().SetProjection(projection))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cur.Close(c)
	for cur.Next(c) {
		var pro entities.Products
		err := cur.Decode(&pro)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		products = append(products, pro)
	}

	c.JSON(http.StatusOK, gin.H{"pages": 1, "docs": products})

}

func GetAllProductsByAdmin(c *gin.Context) {

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
	slug := c.Param("id")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid 'id' parameter"})
		return
	}

	// Define a filter to find the product based on the slug
	filter := bson.M{"slug": slug}
	var product entities.Products
	err := proCollection.FindOne(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	// Initialize categoryInfo to hold category details
	var categoryInfo []map[string]interface{}

	// Iterate over each category in the product's categories field
	for _, catID := range product.Category {
		var category entities.Category
		err := categoryCollection.FindOne(c, bson.M{"_id": catID}).Decode(&category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}

		// Create a map with the category details
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

		// Append the category map to categoryInfo
		categoryInfo = append(categoryInfo, categoryMap)
	}

	// Create a map with the product details, including categoryInfo
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

	// Respond with the product details
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "product", "body": productData})
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateShortID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 3)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateNameFuzzy(name string) []string {
	name = strings.ToLower(name)
	var substrings []string
	for i := 0; i < len(name); i++ {
		for j := i + 1; j <= len(name); j++ {
			substrings = append(substrings, name[i:j])
		}
	}
	return substrings
}

func PostProductByAdmin(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Helper functions
	getBool := func(key string, defaultValue bool) bool {
		if value, ok := input[key]; ok {
			if boolVal, ok := value.(bool); ok {
				return boolVal
			}
		}
		return defaultValue
	}

	getString := func(key string) string {
		if value, ok := input[key]; ok {
			if strVal, ok := value.(string); ok {
				return strVal
			}
		}
		return ""
	}

	getInt := func(key string) int {
		if value, ok := input[key]; ok {
			switch v := value.(type) {
			case float64:
				return int(v)
			case string:
				if intVal, err := strconv.Atoi(v); err == nil {
					return intVal
				}
			}
		}
		return 0
	}

	getStringArray := func(key string) []string {
		if value, ok := input[key]; ok {
			if arr, ok := value.([]interface{}); ok {
				strArr := make([]string, len(arr))
				for i, v := range arr {
					if strVal, ok := v.(string); ok {
						strArr[i] = strVal
					}
				}
				return strArr
			}
		}
		return []string{}
	}

	getObjectIDArray := func(key string) []primitive.ObjectID {
		if value, ok := input[key]; ok {
			if arr, ok := value.([]interface{}); ok {
				objIDArr := make([]primitive.ObjectID, len(arr))
				for i, v := range arr {
					if strVal, ok := v.(string); ok {
						if objID, err := primitive.ObjectIDFromHex(strVal); err == nil {
							objIDArr[i] = objID
						}
					}
				}
				return objIDArr
			}
		}
		return []primitive.ObjectID{}
	}

	// getFloat64 := func(key string) float64 {
	// 	if value, ok := input[key]; ok {
	// 		if floatVal, ok := value.(float64); ok {
	// 			return floatVal
	// 		}
	// 	}
	// 	return 0
	// }

	getDimensions := func(key string) []entities.Dimension {
		if value, ok := input[key]; ok {
			if arr, ok := value.([]interface{}); ok {
				dimensions := make([]entities.Dimension, len(arr))
				for i, v := range arr {
					if _, ok := v.(map[string]interface{}); ok {
						dimensions[i] = entities.Dimension{
							// populate dimension fields from dimMap
						}
					}
				}
				return dimensions
			}
		}
		return []entities.Dimension{}
	}

	getVariations := func(key string) []entities.Variation {
		if value, ok := input[key]; ok {
			if arr, ok := value.([]interface{}); ok {
				variations := make([]entities.Variation, len(arr))
				for i, v := range arr {
					if _, ok := v.(map[string]interface{}); ok {
						variations[i] = entities.Variation{
							// populate variation fields from varMap
						}
					}
				}
				return variations
			}
		}
		return []entities.Variation{}
	}

	imageData := input["images"].([]interface{})
	// Create a slice to hold the ImagePro objects
	var images []entities.ImagePro
	// Iterate over the image data and create ImagePro objects
	for _, img := range imageData {
		imgMap := img.(map[string]interface{})
		// Extract ID and URL from the map
		id := imgMap["_id"].(string) // Assuming ID is stored as a string
		url := imgMap["url"].(string)
		// Convert ID to ObjectID
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			// Handle error
		}
		// Create ImagePro object
		image := entities.ImagePro{
			ID:  &objID,
			Url: url,
		}
		images = append(images, image)
	}

	// Populate product fields
	product := entities.Products{
		ID:              primitive.NewObjectID(),
		Amazing:         getBool("amazing", false),
		IsMillModel:     getBool("isMillModel", false),
		ProductType:     getString("productType"),
		Quantity:        getInt("quantity"),
		Comment:         getStringArray("comments"),
		Parent:          primitive.NilObjectID, // default value
		Category:        getObjectIDArray("categories"),
		Tags:            getStringArray("tags"),
		SimilarProducts: getStringArray("similarProducts"),
		NameFuzzy:       generateNameFuzzy(getString("name")),
		Images:          images,
		Name:            getString("name"),
		Price:           getInt("price"),
		Details:         getString("details"),
		DiscountPercent: float64(getInt("discountPercent")),
		Stock:           getInt("stock"),
		CategoryID:      getString("categoryId"),
		Attributes:      []*entities.Attribute{},
		Slug:            getString("slug"),
		Dimensions:      getDimensions("dimensions"),
		Variations:      getVariations("variations"),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		V:               getInt("__v"),
		ShortID:         getString("shortId"),
		NotExist:        getBool("notExist", false),
		BannerUrl:       getString("bannerUrl"),
		SalesNumber:     getInt("salesNumber"),
	}
	// Handle attributes separately if they need to be populated
	if attributes, ok := input["attributes"].([]interface{}); ok {
		for _, attr := range attributes {
			attrMap := attr.(map[string]interface{})
			attribute := entities.Attribute{
				ID:    primitive.NewObjectID(),
				Key:   attrMap["key"].(string),
				Value: attrMap["value"].(string),
			}
			product.Attributes = append(product.Attributes, &attribute)
		}
	}
	if len(product.Dimensions) > 0 || len(product.Variations) > 0 {
		// Product has dimensions or variations
		product.ProductType = "variations"
	} else {
		// Product does not have dimensions or variations
		product.ProductType = "none"
	}

	if product.ShortID == "" {
		product.ShortID = generateShortID()
	}

	// Set parent to categoryId if provided
	productCategoryID := getString("categoryId")
	if productCategoryID != "" {
		product.Parent, _ = primitive.ObjectIDFromHex(productCategoryID)
	}

	productCategoryID = getString("categoryId")
	if productCategoryID != "" {
		product.Parent, _ = primitive.ObjectIDFromHex(productCategoryID)

		// Fetch category and its ancestors
		categoryIDs, err := fetchCategoryAndAncestors(product.Parent, categoryCollection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Add fetched categories to the product's categories field
		product.Category = categoryIDs
	}

	// Insert the product into the collection
	_, err := proCollection.InsertOne(c, bson.M{
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

func fetchCategoryAndAncestors(categoryID primitive.ObjectID, categoryCollection *mongo.Collection) ([]primitive.ObjectID, error) {
	var categoryIDs []primitive.ObjectID

	for {
		var category entities.Category
		err := categoryCollection.FindOne(context.Background(), bson.M{"_id": categoryID}).Decode(&category)
		if err != nil {
			return nil, err
		}

		categoryIDs = append(categoryIDs, *category.ID)

		// Check if parent ID exists and convert it to primitive.ObjectID
		parentID, err := getParentID(category.Parent)
		if err != nil {
			return nil, err
		}

		// If the parent ID is nil, we've reached the root category
		if parentID.IsZero() {
			break
		}

		// Move to the parent category
		categoryID = parentID
	}

	return categoryIDs, nil
}

func getParentID(parent interface{}) (primitive.ObjectID, error) {
	switch parentValue := parent.(type) {
	case primitive.ObjectID:
		return parentValue, nil
	case nil:
		return primitive.NilObjectID, nil
	case map[string]interface{}:
		oid, ok := parentValue["$oid"].(string)
		if !ok {
			return primitive.NilObjectID, errors.New("parent ID not found or invalid")
		}
		objID, err := primitive.ObjectIDFromHex(oid)
		if err != nil {
			return primitive.NilObjectID, err
		}
		return objID, nil
	default:
		return primitive.NilObjectID, errors.New("parent ID assertion failed")
	}
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

	filter := bson.M{"_id": objectID}
	var product entities.Products
	err = proCollection.FindOneAndDelete(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok_delete", "body": gin.H{}})

}

func mergeAttributes(existing []entities.Attribute, updated interface{}) []entities.Attribute {
	updatedAttributes := []entities.Attribute{}

	// Convert the updated attributes interface to []entities.Attribute
	if updatedAttrs, ok := updated.([]interface{}); ok {
		for _, attr := range updatedAttrs {
			if attrMap, ok := attr.(map[string]interface{}); ok {
				var id primitive.ObjectID
				if idStr, ok := attrMap["_id"].(string); ok {
					parsedID, err := primitive.ObjectIDFromHex(idStr)
					if err == nil {
						id = parsedID
					}
				}

				updatedAttributes = append(updatedAttributes, entities.Attribute{
					ID:    id,
					Key:   attrMap["key"].(string),
					Value: attrMap["value"].(string),
				})
			}
		}
	}

	// Merge existing and updated attributes
	for _, attr := range existing {
		// Check if the attribute with the same key exists in the updated attributes
		exists := false
		for _, updatedAttr := range updatedAttributes {
			if attr.Key == updatedAttr.Key {
				exists = true
				break
			}
		}
		// If the attribute doesn't exist in the updated attributes, retain it
		if !exists {
			updatedAttributes = append(updatedAttributes, attr)
		}
	}

	return updatedAttributes
}

func UpdateProductByAdmin(c *gin.Context) {
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
	updateFields := make(map[string]interface{})

	// Extract fields to update from request JSON
	if err := c.ShouldBindJSON(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove _id from the updateFields if it exists
	if _, exists := updateFields["_id"]; exists {
		delete(updateFields, "_id")
	}
	if _, exists := updateFields["createdAt"]; exists {
		delete(updateFields, "createdAt")
	}

	getBool := func(key string) bool {
		if value, ok := updateFields[key]; ok {
			if boolVal, ok := value.(bool); ok {
				return boolVal
			}
		}
		return false
	}

	getString := func(key string) string {
		if value, ok := updateFields[key]; ok {
			if strVal, ok := value.(string); ok {
				return strVal
			}
		}
		return ""
	}

	getInt := func(value interface{}) int {
		switch v := value.(type) {
		case float64:
			return int(v)
		case string:
			if intVal, err := strconv.Atoi(v); err == nil {
				return intVal
			}
		}
		return 0
	}

	getStringArray := func(key string) []string {
		if value, ok := updateFields[key]; ok {
			if arr, ok := value.([]interface{}); ok {
				strArr := make([]string, len(arr))
				for i, v := range arr {
					if strVal, ok := v.(string); ok {
						strArr[i] = strVal
					}
				}
				return strArr
			}
		}
		return []string{}
	}

	getObjectIDArray := func(key string) []primitive.ObjectID {
		if value, ok := updateFields[key]; ok {
			if arr, ok := value.([]interface{}); ok {
				objIDArr := make([]primitive.ObjectID, len(arr))
				for i, v := range arr {
					if strVal, ok := v.(string); ok {
						if objID, err := primitive.ObjectIDFromHex(strVal); err == nil {
							objIDArr[i] = objID
						}
					}
				}
				return objIDArr
			}
		}
		return []primitive.ObjectID{}
	}

	convertInterfaceArrayToIntArray := func(arr []interface{}) []int {
		intArr := make([]int, len(arr))
		for i, v := range arr {
			if num, ok := v.(float64); ok {
				intArr[i] = int(num)
			}
		}
		return intArr
	}

	// Update fields processing
	if amazing := getBool("amazing"); amazing {
		updateFields["amazing"] = amazing
	}
	if isMillModel := getBool("isMillModel"); isMillModel {
		updateFields["isMillModel"] = isMillModel
	}
	if productType := getString("productType"); productType != "" {
		updateFields["productType"] = productType
	}
	if quantity, exists := updateFields["quantity"]; exists {
		updateFields["quantity"] = getInt(quantity)
	}
	if comments := getStringArray("comments"); len(comments) > 0 {
		updateFields["comments"] = comments
	}
	if categories := getObjectIDArray("categories"); len(categories) > 0 {
		updateFields["categories"] = categories
	}
	if tags := getStringArray("tags"); len(tags) > 0 {
		updateFields["tags"] = tags
	}
	if similarProducts := getStringArray("similarProducts"); len(similarProducts) > 0 {
		updateFields["similarProducts"] = similarProducts
	}
	if nameFuzzy := getStringArray("name_fuzzy"); len(nameFuzzy) > 0 {
		updateFields["name_fuzzy"] = nameFuzzy
	}
	if name := getString("name"); name != "" {
		updateFields["name"] = name
	}
	if price, exists := updateFields["price"]; exists {
		updateFields["price"] = getInt(price)
	}
	if details := getString("details"); details != "" {
		updateFields["details"] = details
	}
	if discountPercent, exists := updateFields["discountPercent"]; exists {
		updateFields["discountPercent"] = float64(getInt(discountPercent))
	}
	if stock, exists := updateFields["stock"]; exists {
		updateFields["stock"] = getInt(stock)
	}
	if categoryId := getString("categoryId"); categoryId != "" {
		if parentObjectID, err := primitive.ObjectIDFromHex(categoryId); err == nil {
			updateFields["parent"] = parentObjectID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'categoryId' parameter"})
			return
		}
	}
	if slug := getString("slug"); slug != "" {
		updateFields["slug"] = slug
	}
	if bannerUrl := getString("bannerUrl"); bannerUrl != "" {
		updateFields["bannerUrl"] = bannerUrl
	}
	if salesNumber, exists := updateFields["salesNumber"]; exists {
		updateFields["salesNumber"] = getInt(salesNumber)
	}
	if shortId := getString("shortId"); shortId != "" {
		updateFields["shortId"] = generateShortID()
	}
	if images, ok := updateFields["images"].([]interface{}); ok {
		imageArray := make([]entities.ImagePro, len(images))
		for i, v := range images {
			if imgMap, ok := v.(map[string]interface{}); ok {
				idStr, _ := imgMap["_id"].(string)
				var imgID *primitive.ObjectID
				if idStr != "" {
					objID, err := primitive.ObjectIDFromHex(idStr)
					if err == nil {
						imgID = &objID
					}
				}
				imageArray[i] = entities.ImagePro{
					ID:  imgID,
					Url: imgMap["url"].(string),
				}
			}
		}
		updateFields["images"] = imageArray
	}

	// Set the updatedAt field to the current time
	updateFields["updatedAt"] = time.Now()

	// Handle dimensions
	if dimensions, ok := updateFields["dimensions"].([]interface{}); ok {
		dimensionsArray := make([]entities.Dimension, len(dimensions))
		for i, v := range dimensions {
			if dimMap, ok := v.(map[string]interface{}); ok {
				// Ensure values is an array
				values, valuesOk := dimMap["values"].([]interface{})
				if !valuesOk {
					values = []interface{}{}
				}
				// Populate dimension fields from dimMap
				key, keyOk := dimMap["key"].(float64)
				if !keyOk {
					key = 0
				}
				dimensionsArray[i] = entities.Dimension{
					ID:     primitive.NewObjectID(),
					Key:    int(key),
					Values: convertInterfaceArrayToIntArray(values),
				}
			}
		}
		updateFields["dimensions"] = dimensionsArray
	}

	if variations, ok := updateFields["variations"].([]interface{}); ok {
		variationsArray := make([]entities.Variation, len(variations))
		for i, v := range variations {
			if varMap, ok := v.(map[string]interface{}); ok {
				// Ensure keys is an array
				keys, keysOk := varMap["keys"].([]interface{})
				if !keysOk {
					keys = []interface{}{}
				}
				variationsArray[i] = entities.Variation{
					Id:              primitive.NewObjectID(),
					Keys:            convertInterfaceArrayToIntArray(keys),
					DiscountPercent: getInt(varMap["discountPercent"]),
					Quantity:        getInt(varMap["quantity"]),
					Price:           getInt(varMap["price"]),
				}
			}
		}
		updateFields["variations"] = variationsArray
	}

	existingAttributes := []entities.Attribute{}
	if err := proCollection.FindOne(c, filter).Decode(&existingAttributes); err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error50": err.Error()})
		return
	}

	// Merge existing attributes with updated attributes
	updateAttributes := mergeAttributes(existingAttributes, updateFields["attributes"])

	// Update the 'attributes' field with the merged attributes
	updateFields["attributes"] = updateAttributes

	// if attributes, ok := updateFields["attributes"].([]interface{}); ok {
	// 	attributesArray := make([]entities.Attribute, len(attributes))
	// 	for i, attr := range attributes {
	// 		if attrMap, ok := attr.(map[string]interface{}); ok {
	// 			attributesArray[i] = entities.Attribute{
	// 				ID:    attrMap["_id"].(string),
	// 				Key:   attrMap["key"].(string),
	// 				Value: attrMap["value"].(string),
	// 			}
	// 		}
	// 	}
	// 	updateFields["attributes"] = attributesArray
	// } else {
	// 	// If attributes is not provided or empty, set it to an empty array
	// 	updateFields["attributes"] = []entities.Attribute{}
	// }

	// Fetch category and its ancestors if categoryId is provided
	if categoryId := getString("categoryId"); categoryId != "" {
		if parentObjectID, err := primitive.ObjectIDFromHex(categoryId); err == nil {
			categoryIDs, err := fetchCategoryAndAncestors(parentObjectID, categoryCollection)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			updateFields["categories"] = categoryIDs
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'categoryId' parameter"})
			return
		}
	}

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

	product.ProductType = "variations"
	update := bson.M{
		"$set": bson.M{
			"dimensions":  product.Dimensions,
			"productType": product.ProductType,
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

	var variationMap map[string]interface{}
	if err := c.ShouldBindJSON(&variationMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	// Convert the fields to the appropriate types
	variation := entities.Variation{}

	if keys, ok := variationMap["keys"].([]interface{}); ok {
		for _, key := range keys {
			if intKey, ok := key.(float64); ok {
				variation.Keys = append(variation.Keys, int(intKey))
			}
		}
	}

	if price, ok := variationMap["price"].(string); ok {
		if intPrice, err := strconv.Atoi(price); err == nil {
			variation.Price = intPrice
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'price' parameter"})
			return
		}
	}

	if discountPercent, ok := variationMap["discountPercent"].(string); ok {
		if intDiscountPercent, err := strconv.Atoi(discountPercent); err == nil {
			variation.DiscountPercent = intDiscountPercent
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'discountPercent' parameter"})
			return
		}
	}

	if quantity, ok := variationMap["quantity"].(string); ok {
		if intQuantity, err := strconv.Atoi(quantity); err == nil {
			variation.Quantity = intQuantity
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'quantity' parameter"})
			return
		}
	}

	var product entities.Products
	filter := bson.M{"_id": objectID}
	if err := proCollection.FindOne(c, filter).Decode(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	variation.Id = primitive.NewObjectID()
	product.Variations = append(product.Variations, variation)

	update := bson.M{
		"$set": bson.M{
			"variations": product.Variations,
		},
	}

	if _, err := proCollection.UpdateOne(c, filter, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "variation_added", "body": gin.H{}})
}
func UpdateVariationByAdmin(c *gin.Context) {
	productId := c.Param("id")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	productObjectID, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'productId' parameter"})
		return
	}

	variationID := c.Param("variationId")
	if variationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'variationId' parameter"})
		return
	}

	variationObjectID, err := primitive.ObjectIDFromHex(variationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'variationId' parameter"})
		return
	}

	filter := bson.M{"_id": productObjectID}
	var product entities.Products
	err = proCollection.FindOne(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	var updatedFields struct {
		Keys            []int           `json:"keys"`
		DiscountPercent json.RawMessage `json:"discountPercent"`
		Quantity        json.RawMessage `json:"quantity"`
		Price           json.RawMessage `json:"price"`
	}

	if err := c.ShouldBindJSON(&updatedFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
		return
	}

	var updatedVariation entities.Variation
	var isUpdated bool

	if len(updatedFields.DiscountPercent) > 0 {
		var discountPercent int
		if err := json.Unmarshal(updatedFields.DiscountPercent, &discountPercent); err != nil {
			var discountPercentStr string
			if err := json.Unmarshal(updatedFields.DiscountPercent, &discountPercentStr); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'discountPercent' value"})
				return
			}
			discountPercent, err = strconv.Atoi(discountPercentStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'discountPercent' value"})
				return
			}
		}
		updatedVariation.DiscountPercent = discountPercent
		isUpdated = true
	}

	if len(updatedFields.Quantity) > 0 {
		var quantity int
		if err := json.Unmarshal(updatedFields.Quantity, &quantity); err != nil {
			var quantityStr string
			if err := json.Unmarshal(updatedFields.Quantity, &quantityStr); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'quantity' value"})
				return
			}
			quantity, err = strconv.Atoi(quantityStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'quantity' value"})
				return
			}
		}
		updatedVariation.Quantity = quantity
		isUpdated = true
	}

	if len(updatedFields.Price) > 0 {
		var price int
		if err := json.Unmarshal(updatedFields.Price, &price); err != nil {
			var priceStr string
			if err := json.Unmarshal(updatedFields.Price, &priceStr); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'price' value"})
				return
			}
			price, err = strconv.Atoi(priceStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'price' value"})
				return
			}
		}
		updatedVariation.Price = price
		isUpdated = true
	}

	if !isUpdated {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	var found bool
	for i, variation := range product.Variations {
		if variation.Id == variationObjectID {

			if len(updatedFields.DiscountPercent) > 0 {
				product.Variations[i].DiscountPercent = updatedVariation.DiscountPercent
			}
			if len(updatedFields.Quantity) > 0 {
				product.Variations[i].Quantity = updatedVariation.Quantity
			}
			if len(updatedFields.Price) > 0 {
				product.Variations[i].Price = updatedVariation.Price
			}
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Variation not found"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"variations": product.Variations,
			"updatedAt":  time.Now(),
		},
	}

	_, err = proCollection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "variations_edited", "body": gin.H{}})
}

type DimensionResponse struct {
	Values []int               `json:"values"`
	ID     string              `json:"_id,omitempty"`
	Key    *PropertiesResponse `json:"key"`
}
type VariationResponse struct {
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

type AttributeResponse struct {
	ID    string `json:"_id,omitempty"` // The ID as a string
	Key   string `json:"key"`           // The attribute key
	Value string `json:"value"`         // The attribute value
}
type CategoryResponse struct {
	ID        string             `json:"_id"`       // The ID as a string
	Images    ImageResponse      `json:"image"`     // Image response structure
	Parent    interface{}        `json:"parent"`    // Keeping this as interface{} to support various types
	Name      string             `json:"name"`      // The category name
	Ancestors []AncestorResponse `json:"ancestors"` // List of ancestor responses
	Slug      string             `json:"slug"`      // Slug for the category
	Details   string             `json:"details"`   // Category details
	Faq       []NewFaqResponse   `json:"faq"`       // List of FAQ responses
}

type ImageResponse struct {
	Url string `json:"url"` // The URL for the image
}

type NewFaqResponse struct {
	ID       string `json:"_id,omitempty"` // The ID as a string
	Answer   string `json:"answer"`        // The answer text
	Complete bool   `json:"completed"`     // The completion status
	Question string `json:"question"`      // The question text
}

type AncestorResponse struct {
	ID   string `json:"_id,omitempty"` // The ID as a string
	Name string `json:"name"`          // The name of the ancestor
	Slug string `json:"slug"`          // The slug for the ancestor
}

type VariationName struct {
	ID        int       `json:"_id"`
	Parent    int       `json:"parent"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Version   int       `json:"__v"`
}

func GetProductByIDByVAlueFalse(c *gin.Context) {

	// Check 'values' parameter
	values := c.Query("values")
	if values != "false" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid 'values' parameter"})
		return
	}

	// Get product slug from request parameters
	slug := c.Param("id")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid 'id' parameter"})
		return
	}

	// Fetch product from the database
	var product entities.Products
	filter := bson.M{"slug": slug}
	err := proCollection.FindOne(c, filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to fetch product"})
		return
	}

	// Prepare dimensions response
	var dimensionResponses []DimensionResponse
	for _, dim := range product.Dimensions {
		var properties entities.Properties
		err = propertiesCollection.FindOne(c, bson.M{"_id": dim.Key}).Decode(&properties)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch dimension key information"})
			return
		}
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
	var variationNames []VariationName
	for _, variat := range product.Variations {
		// Find properties based on the variation's keys
		cursor, err := propertiesCollection.Find(c, bson.M{"_id": bson.M{"$in": variat.Keys}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch variation keys information"})
			return
		}
		defer cursor.Close(c)

		// Process each property
		for cursor.Next(c) {
			var prop entities.Properties
			if err := cursor.Decode(&prop); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to decode properties document"})
				return
			}

			// Create a VariationName from the property
			variationName := VariationName{
				ID:        prop.ID,
				Parent:    *prop.Parent,
				Name:      prop.Name,
				CreatedAt: prop.CreatedAt,
				UpdatedAt: prop.UpdatedAt,
				Version:   prop.V,
			}

			// Add to the slice
			variationNames = append(variationNames, variationName)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Cursor error"})
			return
		}
	}

	// Ensure variationNames is not nil
	if variationNames == nil {
		variationNames = make([]VariationName, 0)
	}

	var categoriesResponses []CategoryResponse
	for _, categoryID := range product.Category {
		var category entities.Category
		err = categoryCollection.FindOne(c, bson.M{"_id": categoryID}).Decode(&category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch category information"})
			return
		}

		// Map category to CategoryResponse
		categoryResponse := CategoryResponse{
			ID:      category.ID.Hex(),
			Name:    category.Name,
			Slug:    category.Slug,
			Details: category.Details,
			Images:  ImageResponse{Url: category.Images.Url},
		}

		// Map ancestors
		for _, ancestor := range category.Ancestors {
			ancestorResponse := AncestorResponse{
				ID:   ancestor.ID.Hex(),
				Name: ancestor.Name,
				Slug: ancestor.Slug,
			}
			categoryResponse.Ancestors = append(categoryResponse.Ancestors, ancestorResponse)
		}

		// Map FAQs
		for _, faq := range category.Faq {
			faqResponse := NewFaqResponse{
				ID:       faq.ID.Hex(),
				Answer:   faq.Answer,
				Complete: faq.Complete,
				Question: faq.Question,
			}
			categoryResponse.Faq = append(categoryResponse.Faq, faqResponse)
		}

		categoriesResponses = append(categoriesResponses, categoryResponse)
	}
	// Prepare attributes response
	var attributesResponses []AttributeResponse
	for _, attr := range product.Attributes {
		attributeResponse := AttributeResponse{
			ID:    attr.ID.Hex(),
			Key:   attr.Key,
			Value: attr.Value,
		}
		attributesResponses = append(attributesResponses, attributeResponse)
	}

	// Construct the response JSON
	responseBody := gin.H{
		"_id":             product.ID.Hex(),
		"dimensions":      dimensionResponses,
		"name":            product.Name,
		"productType":     product.ProductType,
		"price":           product.Price,
		"variations":      product.Variations,
		"categories":      categoriesResponses,
		"categoryId":      product.CategoryID,
		"attributes":      attributesResponses,
		"tags":            product.Tags,
		"bannerUrl":       product.BannerUrl,
		"slug":            product.Slug,
		"discountPercent": product.DiscountPercent,
		"amazing":         product.Amazing,
		"images":          product.Images,
		"variationNames":  variationNames,
		"shortId":         product.ShortID,
	}

	// Send the response
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "product", "body": responseBody})
}

// func GetProductByIDByVAlueFalse(c *gin.Context) {

// 	// Check 'values' parameter
// 	values := c.Query("values")
// 	if values != "false" {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid 'values' parameter"})
// 		return
// 	}

// 	// Get product ID from request parameters
// 	slug := c.Param("id")
// 	if slug == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid 'id' parameter"})
// 		return
// 	}

// 	// Convert ID string to ObjectID
// 	// objectID, err := primitive.ObjectIDFromHex(id)
// 	// if err != nil {
// 	// 	c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid 'id' parameter"})
// 	// 	return
// 	// }

// 	// Fetch product from the database
// 	var product entities.Products
// 	filter := bson.M{"slug": slug}
// 	err := proCollection.FindOne(c, filter).Decode(&product)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to fetch product"})
// 		return
// 	}

// 	// Prepare response dimensions
// 	var dimensionResponses []DimensionResponse
// 	for _, dim := range product.Dimensions {
// 		// Fetch properties document from the properties collection
// 		var properties entities.Properties
// 		err = propertiesCollection.FindOne(c, bson.M{"_id": dim.Key}).Decode(&properties)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch dimension key information"})
// 			return
// 		}

// 		// Create a DimensionResponse using the fetched properties
// 		dimensionResponse := DimensionResponse{
// 			Values: dim.Values,
// 			ID:     dim.ID.Hex(),
// 			Key: &PropertiesResponse{
// 				ID:        properties.ID,
// 				Parent:    properties.Parent,
// 				Name:      properties.Name,
// 				CreatedAt: properties.CreatedAt,
// 				UpdatedAt: properties.UpdatedAt,
// 				V:         properties.V,
// 			},
// 		}
// 		dimensionResponses = append(dimensionResponses, dimensionResponse)
// 	}
// 	var variationsResponses []VariationResponae
// 	for _, variat := range product.Variations {
// 		// Fetch properties documents from the properties collection
// 		cursor, err := propertiesCollection.Find(c, bson.M{"_id": bson.M{"$in": variat.Keys}})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch variation keys information"})
// 			return
// 		}
// 		defer cursor.Close(c)

// 		var properties []*entities.Properties
// 		for cursor.Next(c) {
// 			var prop entities.Properties
// 			if err := cursor.Decode(&prop); err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to decode properties document"})
// 				return
// 			}
// 			properties = append(properties, &prop)
// 		}
// 		if err := cursor.Err(); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Cursor error"})
// 			return
// 		}

// 		// Create VariationResponse using the fetched properties
// 		var variationKeys []PropertiesResponse
// 		for _, prop := range properties {
// 			variationKey := PropertiesResponse{
// 				ID:        prop.ID,
// 				Parent:    prop.Parent,
// 				Name:      prop.Name,
// 				CreatedAt: prop.CreatedAt,
// 				UpdatedAt: prop.UpdatedAt,
// 				V:         prop.V,
// 			}
// 			variationKeys = append(variationKeys, variationKey)
// 		}
// 		variationResponse := VariationResponae{
// 			Keys: variationKeys,
// 		}
// 		variationsResponses = append(variationsResponses, variationResponse)
// 	}
// 	if variationsResponses == nil {
// 		variationsResponses = make([]VariationResponae, 0)
// 	}
// 	// Construct the response JSON similar to the provided format
// 	responseBody := gin.H{
// 		"_id":             product.ID.Hex(),
// 		"dimensions":      dimensionResponses,
// 		"name":            product.Name,
// 		"productType":     product.ProductType,
// 		"price":           product.Price,
// 		"variations":      product.Variations,
// 		"categories":      product.Category,
// 		"categoryId":      product.CategoryID,
// 		"attributes":      product.Attributes,
// 		"tags":            product.Tags,
// 		"bannerUrl":       product.BannerUrl,
// 		"slug":            product.Slug,
// 		"discountPercent": product.DiscountPercent,
// 		"amazing":         product.Amazing,
// 		"images":          product.Images,
// 		"variationNames":  variationsResponses,
// 		"shortId":         product.ShortID,
// 	}

// 	// Send the response
// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "product", "body": responseBody})
// }

///*******************************DELETE VARIATION*************************************************
func DeleteVariationFromProductByAdmin(c *gin.Context) {

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

	var mix entities.MixProducts
	err := c.ShouldBindJSON(&mix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	mix.ID = primitive.NewObjectID()
	mix.CreatedAt = time.Now()
	mix.UpdatedAt = time.Now()

	_, err = mixProductCollection.InsertOne(c, mix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "mix_product"})
}

func DeleteMixBYAdmin(c *gin.Context) {

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
