package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type envelope struct {
	Success bool            `json:"success"`
	Error   string          `json:"error"`
	Data    json.RawMessage `json:"data"`
}

func locationBaseURL() string {
	if v := os.Getenv("LOCATION_SERVICE_URL"); v != "" {
		return v
	}
	return "http://location-service:8002"
}

var httpClient = &http.Client{Timeout: 8 * time.Second}

// LockerContext mirrors location-service's internal context payload.
type LockerContext struct {
	LockerID     string  `json:"lockerId"`
	LockerNumber string  `json:"lockerNumber"`
	LockerSize   string  `json:"lockerSize"`
	LockerStatus string  `json:"lockerStatus"`
	LocationID   string  `json:"locationId"`
	LocationName string  `json:"locationName"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	HourlyRate   float64 `json:"hourlyRate"`
	DailyRate    float64 `json:"dailyRate"`
}

func GetLockerContext(lockerID string) (*LockerContext, error) {
	url := fmt.Sprintf("%s/internal/lockers/%s/context", locationBaseURL(), lockerID)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("location-service unavailable")
	}
	defer resp.Body.Close()

	var env envelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, fmt.Errorf("invalid response from location-service")
	}
	if !env.Success {
		if env.Error == "" {
			env.Error = "locker not found"
		}
		return nil, fmt.Errorf(env.Error)
	}

	var ctx LockerContext
	if err := json.Unmarshal(env.Data, &ctx); err != nil {
		return nil, fmt.Errorf("invalid locker context")
	}
	return &ctx, nil
}

func postLocker(lockerID, action string) error {
	url := fmt.Sprintf("%s/internal/lockers/%s/%s", locationBaseURL(), lockerID, action)
	resp, err := httpClient.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("location-service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	var env envelope
	_ = json.NewDecoder(resp.Body).Decode(&env)
	if env.Error == "" {
		env.Error = fmt.Sprintf("locker %s failed (status %d)", action, resp.StatusCode)
	}
	return fmt.Errorf(env.Error)
}

// ReserveLocker atomically reserves an available locker. Returns an error if
// the locker is not available (409).
func ReserveLocker(lockerID string) error { return postLocker(lockerID, "reserve") }

// ReleaseLocker frees a locker back to AVAILABLE.
func ReleaseLocker(lockerID string) error { return postLocker(lockerID, "release") }

// OccupyLocker marks a locker OCCUPIED.
func OccupyLocker(lockerID string) error { return postLocker(lockerID, "occupy") }
