package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Photo struct {
	gorm.Model
	Title     string  `json:"title"`
	Caption   string  `json:"caption"`
	PhotoUrl  string  `json:"photoUrl"`
	UserID    uint    `json:"userId" gorm:"index"`
	User      *User   `gorm:"foreignKey:UserID"` // One-to-One relationship with User
}

// CreatePhoto creates a new photo record in the database
func CreatePhoto(db *gorm.DB, photo *Photo) error {
	return db.Create(photo).Error
}

// GetPhotos retrieves all photo records from the database
func GetPhotos(db *gorm.DB, photos *[]Photo) error {
	return db.Find(photos).Error
}

// GetPhotoByID retrieves a specific photo record based on its ID
func GetPhotoByID(db *gorm.DB, photoID uint) (*Photo, error) {
	var photo Photo
	if err := db.First(&photo, photoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &photo, nil
}

// UpdatePhoto updates an existing photo record with the provided data
func UpdatePhoto(db *gorm.DB, photoID uint, updatedPhoto *Photo) error {
	return db.Model(&Photo{Model: gorm.Model{ID: photoID}}).Updates(updatedPhoto).Error
}

// DeletePhoto deletes a photo record from the database based on its ID
func DeletePhoto(db *gorm.DB, photoID uint) error {
	var photo Photo
	photo.ID = photoID
	return db.Delete(&photo).Error
}

// ErrRecordNotFound is a custom error for record not found cases
var ErrRecordNotFound = errors.New("record not found")
