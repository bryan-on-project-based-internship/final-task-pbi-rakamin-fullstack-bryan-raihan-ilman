// controllers/user.go:
package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"user-profile-apis/app/models"
	"user-profile-apis/helpers"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func UserRegister(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user data
	if err := validateUser(user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user in database
	if err := models.CreateUser(c.MustGet("db").(*gorm.DB), &user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token, err := helpers.GenerateJWT(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func UserLogin(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	fmt.Println("Attempting to bind JSON...")

	if err := c.BindJSON(&credentials); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("JSON bound successfully. Attempting to find user by email...")

	// Find user by email
	var user models.User
	user, err := models.GetUserByEmail(c.MustGet("db").(*gorm.DB), credentials.Email)
	if err != nil {
		fmt.Println("Error finding user by email:", err)
		if err == models.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Invalid credentials"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("User found successfully. Attempting to compare passwords...")

	if user.Password != credentials.Password {
		fmt.Println("Passwords do not match")
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Invalid credentials"})
		return
	}

	fmt.Println("Passwords match. Attempting to generate JWT token...")

	// Generate JWT token
	token, err := helpers.GenerateJWT(user.ID)
	if err != nil {
		fmt.Println("Error generating JWT token:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("JWT token generated successfully.")

	c.JSON(http.StatusOK, gin.H{"token": token, "userId": user.ID})
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("userId")
	fmt.Println("Received userID:", userID)

	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("JSON bound successfully. Checking user authorization...")

	currentUserID := c.MustGet("userID").(uint)
	parsedUserID, _ := strconv.Atoi(userID)
	if currentUserID != uint(parsedUserID) {
		fmt.Println("Unauthorized to update this user")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to update this user"})
		return
	}

	fmt.Println("User authorized. Attempting to update user in database...")

	// Update user in database
	if err := models.UpdateUser(c.MustGet("db").(*gorm.DB), &updatedUser); err != nil {
		fmt.Println("Error updating user in database:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("User updated successfully in database.")

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("userId")
	fmt.Println("Received userID:", userID)

	// Validate user ID from token
	currentUserID := c.MustGet("userID").(uint)
	parsedUserID, _ := strconv.Atoi(userID)
	if currentUserID != uint(parsedUserID) {
		fmt.Println("Unauthorized to delete this user")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this user"})
		return
	}

	fmt.Println("User authorized. Attempting to delete user from database...")

	// Delete user and related photos from database
	if err := models.DeleteUser(c.MustGet("db").(*gorm.DB), uint(parsedUserID)); err != nil {
		fmt.Println("Error deleting user from database:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("User deleted successfully from database.")

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func validateUser(user models.User) error {
	if len(user.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if !govalidator.IsEmail(user.Email) {
		return fmt.Errorf("invalid email format")
	}
	if len(user.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	return nil
}
