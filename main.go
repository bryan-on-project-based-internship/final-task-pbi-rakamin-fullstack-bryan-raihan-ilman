package main

import (
	_ "github.com/lib/pq"

	"fmt"
	"log"

	"user-profile-apis/app/handlers"
	"user-profile-apis/controllers"
	"user-profile-apis/database"
	"user-profile-apis/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// Database connection
	db, err := database.Connect()
	if err != nil || db == nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize router
	router := gin.Default()
	router.Use(middlewares.DBMiddleware(db))

	// User Endpoints
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", handlers.UserRegister)
		userGroup.POST("/login", handlers.UserLogin)

		// Apply middleware to the following routes
		authorized := userGroup.Group("/")
		authorized.Use(middlewares.JWTAuthMiddleware())
		{
			authorized.PUT("/:userId", handlers.UpdateUser)
			authorized.DELETE("/:userId", handlers.DeleteUser)
		}
	}

	// Photo Endpoints (assuming implemented functions in controllers/photo.go)
	photoGroup := router.Group("/photos")
	photoGroup.Use(middlewares.JWTAuthMiddleware())
	{
		photoGroup.POST("", controllers.UploadPhoto) // Assuming UploadPhoto is implemented
		photoGroup.GET("", controllers.GetPhotos)   // Assuming GetPhotos is implemented
		photoGroup.PUT("/:photoId", controllers.UpdatePhoto) // Assuming UpdatePhoto is implemented
		photoGroup.DELETE("/:photoId", controllers.DeletePhoto) // Assuming DeletePhoto is implemented
	}

	// Start server
	port := ":8080" // replace with your desired port
	fmt.Println("Server listening on port", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
