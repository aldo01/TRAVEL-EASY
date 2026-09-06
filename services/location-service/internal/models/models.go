package models

import (
	"time"

	"gorm.io/gorm"
)

type LocationType string

const (
	LocationKiosk  LocationType = "KIOSK"
	LocationGym    LocationType = "GYM"
	LocationClub   LocationType = "CLUB"
	LocationSchool LocationType = "SCHOOL"
	LocationHotel  LocationType = "HOTEL"
	LocationCafe   LocationType = "CAFE"
	LocationStore  LocationType = "STORE"
	LocationOther  LocationType = "OTHER"
)

type Location struct {
	ID            string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name          string       `gorm:"not null" json:"name"`
	Type          LocationType `gorm:"type:varchar(20);not null" json:"type"`
	Address       string       `gorm:"not null" json:"address"`
	City          string       `gorm:"not null;index" json:"city"`
	State         *string      `json:"state"`
	ZipCode       *string      `json:"zipCode"`
	Country       string       `gorm:"not null" json:"country"`
	Latitude      float64      `gorm:"not null;index" json:"latitude"`
	Longitude     float64      `gorm:"not null;index" json:"longitude"`
	OpeningTime   string       `json:"openingTime"`
	ClosingTime   string       `json:"closingTime"`
	IsOpen24Hours bool         `gorm:"default:false" json:"isOpen24Hours"`
	PhoneNumber   *string      `json:"phoneNumber"`
	Email         *string      `json:"email"`
	Description   *string      `json:"description"`
	ImageURL      *string      `json:"imageUrl"`
	Amenities     string       `gorm:"type:text" json:"amenities"`
	HourlyRate    float64      `gorm:"default:0" json:"hourlyRate"`
	DailyRate     float64      `gorm:"default:0" json:"dailyRate"`
	IsActive      bool         `gorm:"default:true;index" json:"isActive"`
	Rating        float64      `gorm:"default:0" json:"rating"`
	TotalReviews  int          `gorm:"default:0" json:"totalReviews"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`

	Lockers []Locker `gorm:"foreignKey:LocationID" json:"lockers,omitempty"`
}

type LockerSize string

const (
	LockerSmall  LockerSize = "SMALL"
	LockerMedium LockerSize = "MEDIUM"
	LockerLarge  LockerSize = "LARGE"
	LockerXLarge LockerSize = "XLARGE"
)

type LockerStatus string

const (
	LockerAvailable    LockerStatus = "AVAILABLE"
	LockerOccupied     LockerStatus = "OCCUPIED"
	LockerReserved     LockerStatus = "RESERVED"
	LockerMaintenance  LockerStatus = "MAINTENANCE"
	LockerOutOfService LockerStatus = "OUT_OF_SERVICE"
)

type Locker struct {
	ID            string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LocationID    string       `gorm:"type:uuid;not null;index" json:"locationId"`
	LockerNumber  string       `gorm:"not null" json:"lockerNumber"`
	Size          LockerSize   `gorm:"type:varchar(20);not null" json:"size"`
	Floor         *string      `json:"floor"`
	Section       *string      `json:"section"`
	Status        LockerStatus `gorm:"type:varchar(20);default:'AVAILABLE';index" json:"status"`
	IsOperational bool         `gorm:"default:true" json:"isOperational"`
	Height        *int         `json:"height"`
	Width         *int         `json:"width"`
	Depth         *int         `json:"depth"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Location{}, &Locker{})
}
