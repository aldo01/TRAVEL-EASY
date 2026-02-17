m,.?m     package models

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

	Bookings []Booking `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
}

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
	OpeningTime   string       `json:"openingTime"` // "09:00"
	ClosingTime   string       `json:"closingTime"` // "21:00"
	IsOpen24Hours bool         `gorm:"default:false" json:"isOpen24Hours"`
	PhoneNumber   *string      `json:"phoneNumber"`
	Email         *string      `json:"email"`
	Description   *string      `json:"description"`
	ImageURL      *string      `json:"imageUrl"`
	Amenities     string       `gorm:"type:text" json:"amenities"` // JSON array as string
	HourlyRate    float64      `gorm:"default:0" json:"hourlyRate"`
	DailyRate     float64      `gorm:"default:0" json:"dailyRate"`
	IsActive      bool         `gorm:"default:true;index" json:"isActive"`
	Rating        float64      `gorm:"default:0" json:"rating"`
	TotalReviews  int          `gorm:"default:0" json:"totalReviews"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`

	Lockers  []Locker  `gorm:"foreignKey:LocationID" json:"lockers,omitempty"`
	Bookings []Booking `gorm:"foreignKey:LocationID" json:"bookings,omitempty"`
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

	Location Location  `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Bookings []Booking `gorm:"foreignKey:LockerID" json:"bookings,omitempty"`
}

type RateType string

const (
	RateHourly   RateType = "HOURLY"
	RateDaily    RateType = "DAILY"
	RateMultiDay RateType = "MULTI_DAY"
)

type BookingStatus string

const (
	BookingPending   BookingStatus = "PENDING"
	BookingConfirmed BookingStatus = "CONFIRMED"
	BookingActive    BookingStatus = "ACTIVE"
	BookingCompleted BookingStatus = "COMPLETED"
	BookingCancelled BookingStatus = "CANCELLED"
	BookingExpired   BookingStatus = "EXPIRED"
)

type PaymentStatus string

const (
	PaymentPending       PaymentStatus = "PENDING"
	PaymentPaid          PaymentStatus = "PAID"
	PaymentPartiallyPaid PaymentStatus = "PARTIALLY_PAID"
	PaymentRefunded      PaymentStatus = "REFUNDED"
	PaymentFailed        PaymentStatus = "FAILED"
)

type Booking struct {
	ID                 string        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID             string        `gorm:"type:uuid;not null;index" json:"userId"`
	LocationID         string        `gorm:"type:uuid;not null;index" json:"locationId"`
	LockerID           string        `gorm:"type:uuid;not null;index" json:"lockerId"`
	BookingNumber      string        `gorm:"uniqueIndex;not null" json:"bookingNumber"`
	StartTime          time.Time     `gorm:"not null" json:"startTime"`
	EndTime            *time.Time    `json:"endTime"`
	RateType           RateType      `gorm:"type:varchar(20);not null" json:"rateType"`
	BasePrice          float64       `gorm:"not null" json:"basePrice"`
	AdditionalCharges  float64       `gorm:"default:0" json:"additionalCharges"`
	TotalPrice         float64       `gorm:"not null" json:"totalPrice"`
	Status             BookingStatus `gorm:"type:varchar(20);default:'PENDING';index" json:"status"`
	QRCode             string        `gorm:"uniqueIndex;not null" json:"qrCode"`
	QRCodeImageURL     *string       `json:"qrCodeImageUrl"`
	CheckInTime        *time.Time    `json:"checkInTime"`
	CheckOutTime       *time.Time    `json:"checkOutTime"`
	PaymentStatus      PaymentStatus `gorm:"type:varchar(20);default:'PENDING'" json:"paymentStatus"`
	PaymentMethod      *string       `json:"paymentMethod"`
	TransactionID      *string       `json:"transactionId"`
	Notes              *string       `json:"notes"`
	CancellationReason *string       `json:"cancellationReason"`
	CreatedAt          time.Time     `json:"createdAt"`
	UpdatedAt          time.Time     `json:"updatedAt"`

	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Location   Location    `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Locker     Locker      `gorm:"foreignKey:LockerID" json:"locker,omitempty"`
	AccessLogs []AccessLog `gorm:"foreignKey:BookingID" json:"accessLogs,omitempty"`
}

type AccessAction string

const (
	AccessLock   AccessAction = "LOCK"
	AccessUnlock AccessAction = "UNLOCK"
	AccessOpen   AccessAction = "OPEN"
	AccessClose  AccessAction = "CLOSE"
)

type AccessLog struct {
	ID         string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BookingID  string       `gorm:"type:uuid;not null;index" json:"bookingId"`
	Action     AccessAction `gorm:"type:varchar(20);not null" json:"action"`
	Timestamp  time.Time    `gorm:"default:CURRENT_TIMESTAMP;index" json:"timestamp"`
	Method     *string      `json:"method"` // "QR_SCAN", "MANUAL_OVERRIDE"
	DeviceInfo *string      `json:"deviceInfo"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Location{},
		&Locker{},
		&Booking{},
		&AccessLog{},
	)
}
