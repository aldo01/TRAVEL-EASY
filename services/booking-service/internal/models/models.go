package models

import (
	"time"

	"gorm.io/gorm"
)

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

type AccessAction string

const (
	AccessLock   AccessAction = "LOCK"
	AccessUnlock AccessAction = "UNLOCK"
	AccessOpen   AccessAction = "OPEN"
	AccessClose  AccessAction = "CLOSE"
)

// Booking owns booking data. Location/locker details are denormalized snapshots
// captured at booking time (this service does not share the location DB).
type Booking struct {
	ID            string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        string     `gorm:"type:uuid;not null;index" json:"userId"`
	LocationID    string     `gorm:"type:uuid;not null;index" json:"locationId"`
	LockerID      string     `gorm:"type:uuid;not null;index" json:"lockerId"`
	BookingNumber string     `gorm:"uniqueIndex;not null" json:"bookingNumber"`
	StartTime     time.Time  `gorm:"not null" json:"startTime"`
	EndTime       *time.Time `json:"endTime"`
	RateType      RateType   `gorm:"type:varchar(20);not null" json:"rateType"`
	Bags          int        `gorm:"default:1" json:"bags"`
	BasePrice     float64    `gorm:"not null" json:"basePrice"`
	TotalPrice    float64    `gorm:"not null" json:"totalPrice"`

	Status         BookingStatus `gorm:"type:varchar(20);default:'PENDING';index" json:"status"`
	QRCode         string        `gorm:"uniqueIndex;not null" json:"qrCode"`
	CheckInTime    *time.Time    `json:"checkInTime"`
	CheckOutTime   *time.Time    `json:"checkOutTime"`
	PaymentStatus  PaymentStatus `gorm:"type:varchar(20);default:'PENDING'" json:"paymentStatus"`
	PaymentMethod  *string       `json:"paymentMethod"`
	TransactionID  *string       `json:"transactionId"`
	Notes          *string       `json:"notes"`
	CancellationReason *string   `json:"cancellationReason"`

	// Denormalized snapshot fields (from location-service at booking time).
	LocationName    string `json:"locationName"`
	LocationAddress string `json:"locationAddress"`
	LocationCity    string `json:"locationCity"`
	LockerNumber    string `json:"lockerNumber"`
	LockerSize      string `json:"lockerSize"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	AccessLogs []AccessLog `gorm:"foreignKey:BookingID" json:"accessLogs,omitempty"`

	// Computed, not persisted — provides a shape compatible with the frontend.
	Location *BookingLocation `gorm:"-" json:"location,omitempty"`
	Locker   *BookingLocker   `gorm:"-" json:"locker,omitempty"`
}

// BookingLocation is the nested, frontend-friendly view built from snapshots.
type BookingLocation struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
}

type BookingLocker struct {
	ID           string `json:"id"`
	LockerNumber string `json:"lockerNumber"`
	Size         string `json:"size"`
}

// Enrich populates the computed nested objects from the snapshot columns.
func (b *Booking) Enrich() {
	b.Location = &BookingLocation{
		ID:      b.LocationID,
		Name:    b.LocationName,
		Address: b.LocationAddress,
		City:    b.LocationCity,
	}
	b.Locker = &BookingLocker{
		ID:           b.LockerID,
		LockerNumber: b.LockerNumber,
		Size:         b.LockerSize,
	}
}

type AccessLog struct {
	ID         string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BookingID  string       `gorm:"type:uuid;not null;index" json:"bookingId"`
	Action     AccessAction `gorm:"type:varchar(20);not null" json:"action"`
	Timestamp  time.Time    `gorm:"default:CURRENT_TIMESTAMP;index" json:"timestamp"`
	Method     *string      `json:"method"`
	DeviceInfo *string      `json:"deviceInfo"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Booking{}, &AccessLog{})
}
