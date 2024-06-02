package handlers

import (
	"user-profile-apis/controllers"

	"github.com/gin-gonic/gin"
)

func UploadPhoto(c *gin.Context) {
	controllers.UploadPhoto(c) // Assuming a controller function for photo upload exists
}

func GetPhotos(c *gin.Context) {
	controllers.GetPhotos(c) // Assuming a controller function for fetching photos exists
}

func UpdatePhoto(c *gin.Context) {
	controllers.UpdatePhoto(c) // Assuming a controller function for photo update exists
}

func DeletePhoto(c *gin.Context) {
	controllers.DeletePhoto(c) // Assuming a controller function for photo deletion exists
}
