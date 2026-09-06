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

// ReviewCreatedEvent is consumed by location-service to update ratings.
type ReviewCreatedEvent struct {
	LocationID    string  `json:"locationId"`
	AverageRating float64 `json:"average"`
	TotalReviews  int     `json:"total"`
}

// PublishReviewCreated emits the review.created event (best effort).
func PublishReviewCreated(ev ReviewCreatedEvent) {
	if nc == nil {
		log.Println("[NATS] publisher not connected; skipping review.created")
		return
	}
	data, err := json.Marshal(ev)
	if err != nil {
		log.Printf("[NATS] marshal review.created failed: %v", err)
		return
	}
	if err := nc.Publish("review.created", data); err != nil {
		log.Printf("[NATS] publish review.created failed: %v", err)
		return
	}
	log.Printf("[NATS] published review.created for location %s", ev.LocationID)
}
