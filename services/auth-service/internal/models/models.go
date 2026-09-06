package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email           string    `gorm:"uniqueIndex;not null" json:"email"`
	Password        string    `gorm:"not null" json:"-"`
	Name            string    `gorm:"not null" json:"name"`
	PhoneNumber     *string   `json:"phoneNumber"`
	ProfileImageURL *string   `json:"profileImageUrl"`
	IsActive        bool      `gorm:"default:true" json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
