package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username  string `json:"username" gorm:"unique"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	Photo     *Photo `gorm:"foreignKey:UserID"` // One-to-One relationship with Photo
}

// CreateUser creates a new user record in the database.
func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

// GetUserByID retrieves a user by their ID from the database.
func GetUserByID(db *gorm.DB, id uint) (User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, errors.New("user not found")
		}
		return User{}, result.Error
	}
	return user, nil
}

// GetUserByEmail retrieves a user by their email address from the database.
func GetUserByEmail(db *gorm.DB, email string) (User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, errors.New("user not found")
		}
		return User{}, result.Error
	}
	return user, nil
}

// UpdateUser updates an existing user record in the database.
func UpdateUser(db *gorm.DB, user *User) error {
	return db.Save(user).Error
}

// DeleteUser deletes a user record from the database based on their ID.
func DeleteUser(db *gorm.DB, id uint) error {
	user := &User{}
	user.ID = id
	fmt.Println("Attempting to delete user with ID:", id)
	result := db.Delete(user)
	if result.Error != nil {
		fmt.Println("Error deleting user:", result.Error)
		return result.Error
	}
	fmt.Println("User deleted successfully")
	return nil
}