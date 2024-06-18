package main

import (
	"log"
	"net/http"
	"os"
	"shop/adminroute"
	"shop/database"
	"strings"

	"shop/router"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	"go.mongodb.org/mongo-driver/mongo"

	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func main() {
	// password := "A3Yx5hNT83"

	// hashedPassword, err := auth.GenerateBcryptHash(password)
	// if err != nil {
	// 	fmt.Println("Error generating hashed password:", err)
	// 	return
	// }

	// fmt.Println("Bcrypt hash of password:", hashedPassword)
	r := gin.Default()
	r.Use(corsMiddleware())
	// r.Use(handleDoubleSlash())
	r.LoadHTMLGlob("./assets/*")
	v1 := r.Group("api")

	v2 := r.Group("/")
	v3 := r.Group("/admin/api")
	v4 := r.Group("/admin")

	v2.Use(removeDoubleSlashesMiddleware)
	// if err := deleteAllDocumentsInCartCollection(); err != nil {
	// 	log.Fatalf("Error deleting documents in cart collection: %v", err)
	// }

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalln("error loading .env file")
	// }
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8000"
	// }
	// r.Use(corsMiddleware())

	// r.OPTIONS("/yourEndpointPath", func(c *gin.Context) {
	// 	// Set CORS headers
	// 	c.Header("Access-Control-Allow-Origin", "*")
	// 	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// 	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// 	// Respond with 204 No Content
	// 	c.Status(http.StatusNoContent)
	// })

	// v1.Use(corsMiddleware())
	// v2.Use(corsMiddleware())
	// v3.Use(corsMiddleware())

	// Configure CORS middleware
	// corsConfig := cors.DefaultConfig()
	// corsConfig.AllowOrigins = []string{os.Getenv("CORS_DOMAIN")}
	// corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// r.Use(cors.New(corsConfig))

	router.ApiGh(v2)
	router.ProRouter(v1)
	router.CategoryRouter(v1)
	// router.Uploader(v1)
	router.Downloader(v2)
	router.UserRoute(v1)
	router.OrderRouter(v1)
	router.CartRouter(v1)
	router.BrandRoute(v1)
	router.PageRoute(v1)
	router.CommentRoute(v1)
	router.FavoriteRoute(v1)
	router.AssetsRoute(v2)
	router.CouponRoute(v1)
	router.State(v2)

	//**********************
	adminroute.AdminCategoryRouter(v3)
	adminroute.AdminRoutes(v3)
	adminroute.AdminBrandRoute(v3)
	adminroute.AdminProductRoute(v3)
	adminroute.Comments(v3)
	adminroute.UserRouteByAdmin(v3)
	adminroute.PropertiesByAdmin(v3)
	adminroute.ImageRouteByAdmin(v3)
	adminroute.PageRouteByAdmin(v3)
	adminroute.OrderRouteByAdmin(v3)
	adminroute.CouponRoutesByAdmin(v3)
	adminroute.DownloaderByAdmin(v4)

	go func() {
		database.MD()

	}()

	// Middleware to conditionally handle /swagger/doc.json separately
	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/swagger/doc.json" {
			// Handle /swagger/doc.json separately if needed
			c.File("./docs/swagger.json")
			c.Abort()
			return
		}
		c.Next()
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + getPort())

}
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8000"

	}
	return port
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("error loading .env file")
		}
		cors := os.Getenv("CORS_DOMAIN")
		c.Writer.Header().Set("Access-Control-Allow-Origin", cors)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func removeDoubleSlashesMiddleware(c *gin.Context) {

	requestPath := c.Request.URL.Path
	cleanedPath := cleanPath(requestPath)

	if cleanedPath != requestPath {
		c.Redirect(http.StatusMovedPermanently, cleanedPath)
		return
	}

	c.Next()
}

func cleanPath(path string) string {
	parts := strings.Split(path, "/")
	cleanedParts := []string{}

	for _, part := range parts {
		if part != "" {
			cleanedParts = append(cleanedParts, part)
		}
	}

	return "/" + strings.Join(cleanedParts, "/")
}

var cartCollection2 *mongo.Collection = database.GetCollection(database.DB, "carts")

// func deleteAllDocumentsInCartCollection() error {
// 	_, err := cartCollection2.DeleteMany(context.Background(), bson.M{})
// 	if err != nil {
// 		return fmt.Errorf("delete error: %v", err)
// 	}
// 	fmt.Println("Deleted all documents in carts collection.")
// 	return nil
// }
