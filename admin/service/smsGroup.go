package service

import (
	"context"
	"log"
	"net/http"
	"shop/entities"
	"shop/helpers"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func SendMessageToAllUsers(c *gin.Context) {
	var text entities.TextSms
	err := c.ShouldBindJSON(&text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	var users []entities.Users
	res, err := userCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
		return
	}
	defer res.Close(context.TODO())

	for res.Next(context.TODO()) {
		var user entities.Users
		err := res.Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error3": err.Error()})
			return
		}
		users = append(users, user)
	}

	log.Printf("Sending SMS to %d users", len(users))

	var phoneNumbers []string
	var messages []string

	for _, user := range users {
		phone := "+98" + strings.TrimPrefix(user.PhoneNumber, "0")
		personalizedMessage := user.Name + " عزیز\n" + text.Text
		phoneNumbers = append(phoneNumbers, phone)
		messages = append(messages, personalizedMessage)
	}

	for i, phoneNumber := range phoneNumbers {
		err = helpers.SendSmstoAllUsers(messages[i], []string{phoneNumber})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error4": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "process_running", "body": gin.H{}})
}

// func SendMessageToAllUsers(c *gin.Context) {
// 	var text entities.TextSms
// 	err := c.ShouldBindJSON(&text)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
// 		return
// 	}

// 	var users []string
// 	res, err := userCollection.Find(c, bson.M{})
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
// 		return
// 	}
// 	defer res.Close(c)

// 	for res.Next(c) {
// 		var user entities.Users
// 		err := res.Decode(&user)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error3": err.Error()})
// 			return
// 		}
// 		// Assuming the phone number starts with '0' and should be prefixed with '+98'
// 		phone := "0" + user.PhoneNumber[0:]
// 		users = append(users, phone)
// 	}

// 	fmt.Println("phone:", users)

// 	err = helpers.SendSmstoAllUsers(text.Text, users)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error4": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "process_running", "body": gin.H{}})
// }
