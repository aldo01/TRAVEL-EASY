package events

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"location-service/config"
	"location-service/internal/models"

	"github.com/nats-io/nats.go"
)

// ReviewCreatedEvent is published by review-service when a new review lands.
type ReviewCreatedEvent struct {
	LocationID    string  `json:"locationId"`
	AverageRating float64 `json:"average"`
	TotalReviews  int     `json:"total"`
}

// SubscribeReviewEvents connects to NATS and keeps a location's rating in sync.
// It runs in the background and reconnects are handled by the NATS client.
func SubscribeReviewEvents() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	var nc *nats.Conn
	var err error
	// Retry connect: the broker may start slightly after this service.
	for i := 0; i < 10; i++ {
		nc, err = nats.Connect(url,
			nats.RetryOnFailedConnect(true),
			nats.MaxReconnects(-1),
			nats.ReconnectWait(2*time.Second),
		)
		if err == nil {
			break
		}
		log.Printf("[NATS] connect attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Printf("[NATS] could not connect, review rating sync disabled: %v", err)
		return
	}

	_, err = nc.Subscribe("review.created", func(msg *nats.Msg) {
		var ev ReviewCreatedEvent
		if err := json.Unmarshal(msg.Data, &ev); err != nil {
			log.Printf("[NATS] bad review.created payload: %v", err)
			return
		}
		if ev.LocationID == "" {
			return
		}

		if err := config.GetDB().Model(&models.Location{}).
			Where("id = ?", ev.LocationID).
			Updates(map[string]interface{}{
				"rating":        ev.AverageRating,
				"total_reviews": ev.TotalReviews,
			}).Error; err != nil {
			log.Printf("[NATS] failed to update location rating: %v", err)
			return
		}
		log.Printf("[NATS] updated rating for location %s -> %.2f (%d reviews)",
			ev.LocationID, ev.AverageRating, ev.TotalReviews)
	})
	if err != nil {
		log.Printf("[NATS] failed to subscribe review.created: %v", err)
		return
	}

	log.Println("[NATS] subscribed to review.created")
}
