package services

import (
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

	couponCode := c.Param("id")
	var coupon entities.Coupons
	err := coupontCollection.FindOne(c, bson.M{"couponCode": couponCode}).Decode(&coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims not found in context"})
		return
	}

	_, ok := tokenClaims.(*auth.SignedUserDetails)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims type"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "coupon", "body": coupon})

}
