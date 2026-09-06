package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var httpClient = &http.Client{Timeout: 6 * time.Second}

type envelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

func authBaseURL() string {
	if v := os.Getenv("AUTH_SERVICE_URL"); v != "" {
		return v
	}
	return "http://auth-service:8001"
}

type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUser fetches a user's display name from auth-service (best effort).
func GetUser(userID string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/internal/users/%s", authBaseURL(), userID)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var env envelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, err
	}
	if !env.Success {
		return nil, fmt.Errorf("user not found")
	}

	var u UserInfo
	if err := json.Unmarshal(env.Data, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
