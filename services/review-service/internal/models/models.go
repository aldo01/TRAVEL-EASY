package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID         string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LocationID string    `gorm:"type:uuid;not null;index;uniqueIndex:idx_review_user_location" json:"locationId"`
	UserID     string    `gorm:"type:uuid;not null;uniqueIndex:idx_review_user_location" json:"userId"`
	UserName   string    `json:"userName"`
	Rating     int       `gorm:"not null" json:"rating"`
	Comment    string    `gorm:"type:text" json:"comment"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Review{})
}
