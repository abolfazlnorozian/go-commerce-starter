package upload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var imgCollection *mongo.Collection = database.GetCollection(database.DB, "uploadschemas")

func Uploadpath(ctx *gin.Context) {
	if err := auth.CheckUserType(ctx, "admin"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename

	out, err := os.Create("public/images/" + filename)
	if err != nil {
		log.Fatal(err)

	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)

	}

	filepath := "uploads/" + filename
	var img entities.Images
	if err := ctx.Bind(&img); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	img.Url = &filepath
	img.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	img.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	_, err = imgCollection.InsertOne(ctx, &img)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": &img})

}
