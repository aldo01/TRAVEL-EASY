package transport

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"

	"notification-service/internal/model"
	"notification-service/internal/notifier"
)

// StartNATSSubscriber consumes booking.confirmed events and dispatches them
// to the notifier service. It runs in the background.
func StartNATSSubscriber(svc *notifier.Service) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	var nc *nats.Conn
	var err error
	for i := 0; i < 10; i++ {
		nc, err = nats.Connect(url,
			nats.RetryOnFailedConnect(true),
			nats.MaxReconnects(-1),
			nats.ReconnectWait(2*time.Second),
		)
		if err == nil {
			break
		}
		log.Printf("[nats] connect attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Printf("[nats] could not connect, event notifications disabled: %v", err)
		return
	}

	_, err = nc.Subscribe("booking.confirmed", func(msg *nats.Msg) {
		var n model.BookingNotification
		if err := json.Unmarshal(msg.Data, &n); err != nil {
			log.Printf("[nats] bad booking.confirmed payload: %v", err)
			return
		}
		if err := svc.HandleBookingConfirmed(n); err != nil {
			log.Printf("[nats] failed to process booking.confirmed: %v", err)
		}
	})
	if err != nil {
		log.Printf("[nats] failed to subscribe booking.confirmed: %v", err)
		return
	}

	log.Println("[nats] subscribed to booking.confirmed")
}
