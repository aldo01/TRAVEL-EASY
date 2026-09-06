package clients

import (
	"encoding/json"
	"fmt"
	"os"
)

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
	Phone string `json:"phoneNumber"`
}

// GetUser fetches a user's public info from auth-service (best effort).
func GetUser(userID string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/internal/users/%s", authBaseURL(), userID)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("auth-service unavailable")
	}
	defer resp.Body.Close()

	var env envelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, fmt.Errorf("invalid response from auth-service")
	}
	if !env.Success {
		return nil, fmt.Errorf("user not found")
	}

	var u UserInfo
	if err := json.Unmarshal(env.Data, &u); err != nil {
		return nil, fmt.Errorf("invalid user payload")
	}
	return &u, nil
}
