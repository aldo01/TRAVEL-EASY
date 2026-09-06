package model

// BookingNotification is the payload for a booking confirmation, received both
// over the NATS "booking.confirmed" subject and the HTTP endpoint.
type BookingNotification struct {
	BookingNumber string  `json:"bookingNumber" binding:"required"`
	UserEmail     string  `json:"userEmail" binding:"required,email"`
	UserName      string  `json:"userName" binding:"required"`
	UserPhone     string  `json:"userPhone"`
	LocationName  string  `json:"locationName" binding:"required"`
	LocationAddr  string  `json:"locationAddress"`
	StartTime     string  `json:"startTime" binding:"required"`
	EndTime       string  `json:"endTime"`
	Bags          int     `json:"bags"`
	RateType      string  `json:"rateType"`
	TotalPrice    float64 `json:"totalPrice"`
	QRCode        string  `json:"qrCode"`
}
