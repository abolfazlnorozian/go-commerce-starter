package services

import (
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var coupontCollection *mongo.Collection = database.GetCollection(database.DB, "coupons")

// func GetCouponByCouponCode(c *gin.Context) {
// 	tokenClaims, exists := c.Get("tokenClaims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims not found in context"})
// 		return
// 	}

// 	claims, ok := tokenClaims.(*auth.SignedUserDetails)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims type"})
// 		return
// 	}

// 	couponCode := c.Param("id")
// 	var coupon entities.Coupons
// 	err := coupontCollection.FindOne(c, bson.M{"couponCode": couponCode}).Decode(&coupon)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "coupon_not_exist"})
// 		return
// 	}
// 	var categoryStr *string
// 	if coupon.Category == primitive.NilObjectID {
// 		categoryStr = nil
// 	} else {
// 		categoryHex := coupon.Category.Hex()
// 		categoryStr = &categoryHex
// 	}

// 	// Handle User as a nullable string
// 	var userStr *string
// 	if coupon.User == "" {
// 		userStr = nil
// 	} else {
// 		userStr = &coupon.User
// 	}

// 	// Create response with adjusted Category and User
// 	couponResponse := entities.CouponResponse{
// 		ID:              coupon.ID.Hex(),
// 		IsExpired:       coupon.IsExpired,
// 		WithoutDiscount: coupon.WithoutDiscount,
// 		CouponCode:      coupon.CouponCode,
// 		Amount:          coupon.Amount,
// 		MinimumPurchase: coupon.MinimumPurchase,
// 		DateExpired:     coupon.DateExpired,
// 		To:              coupon.To,
// 		User:            userStr,
// 		Category:        categoryStr,
// 		CreatedAt:       coupon.CreatedAt,
// 		UpdatedAt:       coupon.UpdatedAt,
// 		V:               coupon.V,
// 	}
// 	if coupon.DateExpired.Before(time.Now()) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Coupon has expired"})
// 		return
// 	}

// 	// var categoryMatch bool
// 	filter := bson.M{"username": claims.Username}
// 	var cart entities.Catrs
// 	err = cartCollection.FindOne(c, filter).Decode(&cart)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
// 		return
// 	}
// 	totalPrice := 0

// 	if len(cart.Mix) == 0 {
// 		for _, product := range cart.Products {
// 			var retrievedProduct entities.Products
// 			err := prodCollection.FindOne(c, bson.M{"_id": product.ProductId}).Decode(&retrievedProduct)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product", "details": err.Error()})
// 				return
// 			}
// 			if retrievedProduct.ProductType == "variations" {
// 				// Calculate the total price for variations
// 				for _, variation := range retrievedProduct.Variations {
// 					if containsAllKeys(product.VariationsKey, variation.Keys) {
// 						totalPrice += variation.Price * product.Quantity
// 						break
// 					}
// 				}
// 			} else {
// 				// Calculate the total price for the simple product
// 				totalPrice += retrievedProduct.Price * product.Quantity
// 			}
// 			var categoryMatch bool
// 			if coupon.Category == primitive.NilObjectID {
// 				// Coupon category is empty, allow it to be applied to all products
// 				categoryMatch = true
// 			} else {
// 				// Coupon category is not empty, check if it matches any product category
// 				categoryMatch = false
// 				for _, catID := range retrievedProduct.Category {
// 					if catID == coupon.Category {
// 						categoryMatch = true
// 						break
// 					}
// 				}
// 			}

// 			// If no category match found, return error
// 			if !categoryMatch {
// 				c.JSON(http.StatusBadRequest, gin.H{"error": "Category mismatch between coupon and product"})
// 				return
// 			}

// 		}

// 		if totalPrice < coupon.MinimumPurchase {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Total price is lower than the minimum purchase required for this coupon"})
// 			return
// 		}

// 		update := bson.M{"$set": bson.M{"couponCode": couponCode}}
// 		_, err := cartCollection.UpdateOne(c, filter, update)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply coupon to cart", "details": err.Error()})
// 			return
// 		}

// 	} else if len(cart.Products) == 0 {

// 	} else if len(cart.Products) > 0 && len(cart.Mix) > 0 {
// 		for _, product := range cart.Products {
// 			var retrievedProduct entities.Products
// 			err := prodCollection.FindOne(c, bson.M{"_id": product.ProductId}).Decode(&retrievedProduct)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product", "details": err.Error()})
// 				return
// 			}
// 			totalPrice += retrievedProduct.Price * product.Quantity

// 		}

// 		if totalPrice < coupon.MinimumPurchase {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Total price is lower than the minimum purchase required for this coupon"})
// 			return
// 		}

// 		update := bson.M{"$set": bson.M{"couponCode": couponCode}}
// 		_, err := cartCollection.UpdateOne(c, filter, update)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply coupon to cart", "details": err.Error()})
// 			return
// 		}

// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "coupon", "body": couponResponse})

// }

// func containsAllKeys(variationsKey []int, keys []int) bool {
// 	keysMap := make(map[int]bool)
// 	for _, key := range keys {
// 		keysMap[key] = true
// 	}
// 	for _, key := range variationsKey {
// 		if !keysMap[key] {
// 			return false
// 		}
// 	}
// 	return true
// }
func GetCouponByCouponCode(c *gin.Context) {
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims not found in context"})
		return
	}

	claims, ok := tokenClaims.(*auth.SignedUserDetails)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims type"})
		return
	}

	couponCode := c.Param("id")
	var coupon entities.Coupons
	err := coupontCollection.FindOne(c, bson.M{"couponCode": couponCode}).Decode(&coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "coupon_not_exist"})
		return
	}
	var categoryStr *string
	if coupon.Category == primitive.NilObjectID {
		categoryStr = nil
	} else {
		categoryHex := coupon.Category.Hex()
		categoryStr = &categoryHex
	}

	// Handle User as a nullable string
	var userStr *string
	if coupon.User == "" {
		userStr = nil
	} else {
		userStr = &coupon.User
	}

	// Create response with adjusted Category and User
	couponResponse := entities.CouponResponse{
		ID:              coupon.ID.Hex(),
		IsExpired:       coupon.IsExpired,
		WithoutDiscount: coupon.WithoutDiscount,
		CouponCode:      coupon.CouponCode,
		Amount:          coupon.Amount,
		MinimumPurchase: coupon.MinimumPurchase,
		DateExpired:     coupon.DateExpired,
		To:              coupon.To,
		User:            userStr,
		Category:        categoryStr,
		CreatedAt:       coupon.CreatedAt,
		UpdatedAt:       coupon.UpdatedAt,
		V:               coupon.V,
	}
	if coupon.DateExpired.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coupon has expired"})
		return
	}

	// var categoryMatch bool
	filter := bson.M{"username": claims.Username}
	var cart entities.Catrs
	err = cartCollection.FindOne(c, filter).Decode(&cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	totalPrice := 0

	if len(cart.Mix) == 0 {
		var retrievedProducts []entities.Products
		for _, product := range cart.Products {
			var retrievedProduct entities.Products
			err := prodCollection.FindOne(c, bson.M{"_id": product.ProductId}).Decode(&retrievedProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product", "details": err.Error()})
				return
			}
			if retrievedProduct.ProductType == "variations" {
				// Calculate the total price for variations
				for _, variation := range retrievedProduct.Variations {
					if containsAllKeys(product.VariationsKey, variation.Keys) {
						totalPrice += variation.Price * product.Quantity
						break
					}
				}
			} else {
				// Calculate the total price for the simple product
				totalPrice += retrievedProduct.Price * product.Quantity
			}
			retrievedProducts = append(retrievedProducts, retrievedProduct)
		}

		var categoryMatch bool
		if coupon.Category == primitive.NilObjectID {
			// Coupon category is empty, allow it to be applied to all products
			categoryMatch = true
		} else {
			// Coupon category is not empty, check if it matches any product category
			categoryMatch = false
			for _, retrievedProduct := range retrievedProducts {
				for _, catID := range retrievedProduct.Category {
					if catID == coupon.Category {
						categoryMatch = true
						break
					}
				}
				if categoryMatch {
					break
				}
			}
		}

		// If no category match found, return error
		if !categoryMatch {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category mismatch between coupon and product"})
			return
		}

		if totalPrice < coupon.MinimumPurchase {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Total price is lower than the minimum purchase required for this coupon"})
			return
		}

		update := bson.M{"$set": bson.M{"couponCode": couponCode}}
		_, err := cartCollection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply coupon to cart", "details": err.Error()})
			return
		}
	} else if len(cart.Products) == 0 {
		// Handle case where only mix products are present (if applicable)
	} else if len(cart.Products) > 0 && len(cart.Mix) > 0 {
		var retrievedProducts []entities.Products
		for _, product := range cart.Products {
			var retrievedProduct entities.Products
			err := prodCollection.FindOne(c, bson.M{"_id": product.ProductId}).Decode(&retrievedProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product", "details": err.Error()})
				return
			}
			if retrievedProduct.ProductType == "variations" {
				// Calculate the total price for variations
				for _, variation := range retrievedProduct.Variations {
					if containsAllKeys(product.VariationsKey, variation.Keys) {
						totalPrice += variation.Price * product.Quantity
						break
					}
				}
			} else {
				// Calculate the total price for the simple product
				totalPrice += retrievedProduct.Price * product.Quantity
			}
			retrievedProducts = append(retrievedProducts, retrievedProduct)
		}

		var categoryMatch bool
		if coupon.Category == primitive.NilObjectID {
			// Coupon category is empty, allow it to be applied to all products
			categoryMatch = true
		} else {
			// Coupon category is not empty, check if it matches any product category
			categoryMatch = false
			for _, retrievedProduct := range retrievedProducts {
				for _, catID := range retrievedProduct.Category {
					if catID == coupon.Category {
						categoryMatch = true
						break
					}
				}
				if categoryMatch {
					break
				}
			}
		}

		// If no category match found, return error
		if !categoryMatch {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category mismatch between coupon and product"})
			return
		}

		if totalPrice < coupon.MinimumPurchase {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Total price is lower than the minimum purchase required for this coupon"})
			return
		}

		update := bson.M{"$set": bson.M{"couponCode": couponCode}}
		_, err := cartCollection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply coupon to cart", "details": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "coupon", "body": couponResponse})
}

func containsAllKeys(variationsKey []int, keys []int) bool {
	keysMap := make(map[int]bool)
	for _, key := range keys {
		keysMap[key] = true
	}
	for _, key := range variationsKey {
		if !keysMap[key] {
			return false
		}
	}
	return true
}
