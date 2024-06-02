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
	"golang.org/x/crypto/bcrypt"
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

	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = string(hashedPassword)

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

	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	user, err := models.GetUserByEmail(c.MustGet("db").(*gorm.DB), credentials.Email)
	if err != nil {
		if err == models.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Invalid credentials"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Invalid credentials"})
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

func UpdateUser(c *gin.Context) {
	userID := c.Param("userId")
	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID := c.MustGet("userID").(uint)
	parsedUserID, _ := strconv.Atoi(userID)
	if currentUserID != uint(parsedUserID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to update this user"})
		return
	}

	// Update user in database
	if err := models.UpdateUser(c.MustGet("db").(*gorm.DB), &updatedUser); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("userId")

	// Validate user ID from token
	currentUserID := c.MustGet("userID").(uint)
	parsedUserID, _ := strconv.Atoi(userID)
	if currentUserID != uint(parsedUserID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this user"})
		return
	}

	// Delete user and related photos from database
	if err := models.DeleteUser(c.MustGet("db").(*gorm.DB), uint(parsedUserID)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func validateUser(user models.User) error {
	if len(user.Username) < 3 {
		return fmt.Errorf("Username must be at least 3 characters long")
	}
	if !govalidator.IsEmail(user.Email) {
		return fmt.Errorf("Invalid email format")
	}
	if len(user.Password) < 6 {
		return fmt.Errorf("Password must be at least 6 characters long")
	}
	return nil
}
