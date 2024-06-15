package services

import (
	"fmt"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var coupontCollection *mongo.Collection = database.GetCollection(database.DB, "coupons")

func GetCouponeByCouponCod(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}
	var categoryStr *string
	if coupon.To == "all" {
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

	filter := bson.M{"username": claims.Username}
	var cart entities.Catrs
	err = cartCollection.FindOne(c, filter).Decode(&cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	totalPrice := 0
	if len(cart.Mix) == 0 {
		for _, product := range cart.Products {
			var retrievedProduct entities.Products
			err := prodCollection.FindOne(c, bson.M{"_id": product.ProductId}).Decode(&retrievedProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product", "details": err.Error()})
				return
			}
			totalPrice += retrievedProduct.Price * product.Quantity
		}

		if totalPrice < coupon.MinimumPurchase {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Total price is lower than the minimum purchase required for this coupon"})
			fmt.Println(":", err)
			return
		}

		update := bson.M{"$set": bson.M{"couponCode": couponCode}}
		_, err := cartCollection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply coupon to cart", "details": err.Error()})
			fmt.Println("tt:", err)
			return
		}

	} else if len(cart.Products) == 0 {

	} else if len(cart.Products) > 0 && len(cart.Mix) > 0 {
		for _, product := range cart.Products {
			var retrievedProduct entities.Products
			err := prodCollection.FindOne(c, bson.M{"_id": product.ProductId}).Decode(&retrievedProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product", "details": err.Error()})
				return
			}
			totalPrice += retrievedProduct.Price * product.Quantity
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
