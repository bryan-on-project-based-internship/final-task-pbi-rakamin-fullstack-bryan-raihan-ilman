package controllers

import (
	"net/http"
	"strconv"

	"user-profile-apis/app/models" // Replace with the actual path to your models package

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func UploadPhoto(c *gin.Context) {
	var photo models.Photo
	userID := c.MustGet("userID").(uint)

	if err := c.BindJSON(&photo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate photo data (optional)

	photo.UserID = userID

	// Save photo (consider using cloud storage for efficiency)
	if err := models.CreatePhoto(c.MustGet("db").(*gorm.DB), &photo); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo uploaded successfully", "photo": photo})
}

func GetPhotos(c *gin.Context) {
	var photos []models.Photo
	if err := models.GetPhotos(c.MustGet("db").(*gorm.DB), &photos); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"photos": photos})
}

func UpdatePhoto(c *gin.Context) {
	photoID := c.Param("photoId")
	var updatedPhoto models.Photo
	if err := c.BindJSON(&updatedPhoto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate photo data (optional)

	// Get photo by ID (using correct function call with db and ID)
	photo, err := models.GetPhotoByID(c.MustGet("db").(*gorm.DB), parsePhotoID(photoID)) // Assuming photoID is a string
	if err != nil {
		if err == models.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check authorization (ensure user owns the photo)
	if photo.UserID != c.MustGet("userID").(uint) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to update this photo"})
		return
	}

	// Update photo in database (using correct function call with all arguments)
	if err := models.UpdatePhoto(c.MustGet("db").(*gorm.DB), photo.ID, &updatedPhoto); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo updated successfully", "photo": updatedPhoto})
}

func DeletePhoto(c *gin.Context) {
	photoID := c.Param("photoId")

	// Get photo by ID (using correct function call with db and ID)
	photo, err := models.GetPhotoByID(c.MustGet("db").(*gorm.DB), parsePhotoID(photoID)) // Assuming photoID is a string
	if err != nil {
		if err == models.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check authorization (ensure user owns the photo)
	if photo.UserID != c.MustGet("userID").(uint) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this photo"})
		return
	}

	// Delete photo from database (using correct function call with ID)
	if err := models.DeletePhoto(c.MustGet("db").(*gorm.DB), photo.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

func parsePhotoID(photoID string) uint {
	parsedID, err := strconv.ParseUint(photoID, 10, 32)
	if err != nil {
		return 0
	}
	return uint(parsedID)
}
