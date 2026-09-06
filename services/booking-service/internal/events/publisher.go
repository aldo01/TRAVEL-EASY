package events

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	nc   *nats.Conn
	once sync.Once
)

func connect() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	var err error
	nc, err = nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		log.Printf("[NATS] publisher connect failed: %v", err)
		nc = nil
		return
	}
	log.Println("[NATS] publisher connected")
}

// Init establishes the shared publisher connection.
func Init() {
	once.Do(connect)
}

// BookingConfirmedEvent is consumed by notification-service.
type BookingConfirmedEvent struct {
	BookingNumber   string  `json:"bookingNumber"`
	UserEmail       string  `json:"userEmail"`
	UserName        string  `json:"userName"`
	UserPhone       string  `json:"userPhone"`
	LocationName    string  `json:"locationName"`
	LocationAddress string  `json:"locationAddress"`
	StartTime       string  `json:"startTime"`
	EndTime         string  `json:"endTime"`
	Bags            int     `json:"bags"`
	RateType        string  `json:"rateType"`
	TotalPrice      float64 `json:"totalPrice"`
	QRCode          string  `json:"qrCode"`
}

// PublishBookingConfirmed emits the event (best-effort, non-blocking on failure).
func PublishBookingConfirmed(ev BookingConfirmedEvent) {
	if nc == nil {
		log.Println("[NATS] publisher not connected; skipping booking.confirmed")
		return
	}
	data, err := json.Marshal(ev)
	if err != nil {
		log.Printf("[NATS] marshal booking.confirmed failed: %v", err)
		return
	}
	if err := nc.Publish("booking.confirmed", data); err != nil {
		log.Printf("[NATS] publish booking.confirmed failed: %v", err)
		return
	}
	log.Printf("[NATS] published booking.confirmed for %s", ev.BookingNumber)
}
