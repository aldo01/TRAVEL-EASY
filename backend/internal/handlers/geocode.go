package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"locker-storage-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type nominatimResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

type PlaceSuggestion struct {
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Title       string  `json:"title"`
	Subtitle    string  `json:"subtitle"`
	DisplayName string  `json:"displayName"`
}

func splitTitleSubtitle(displayName string) (string, string) {
	parts := strings.Split(displayName, ",")
	if len(parts) == 0 {
		return displayName, ""
	}
	title := strings.TrimSpace(parts[0])
	if len(parts) == 1 {
		return title, ""
	}
	subtitle := strings.TrimSpace(strings.Join(parts[1:], ","))
	return title, subtitle
}

// SuggestPlaces proxies OpenStreetMap Nominatim suggestions.
// Query: q (required), limit (optional, max 8)
func SuggestPlaces(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		utils.ErrorResponse(c, 400, "Query q is required")
		return
	}

	limit := 5
	if s := c.Query("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			if n < 1 {
				limit = 1
			} else if n > 8 {
				limit = 8
			} else {
				limit = n
			}
		}
	}

	u, _ := url.Parse("https://nominatim.openstreetmap.org/search")
	params := u.Query()
	params.Set("format", "json")
	params.Set("q", q)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("addressdetails", "1")
	u.RawQuery = params.Encode()

	client := &http.Client{Timeout: 6 * time.Second}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to create request")
		return
	}

	// Nominatim usage policy asks for a valid identifying User-Agent.
	req.Header.Set("User-Agent", "TRAVEL-EASY/1.0 (demo)")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorResponse(c, 502, "Geocoding provider unavailable")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		utils.ErrorResponse(c, 502, "Geocoding provider error")
		return
	}

	var results []nominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		utils.ErrorResponse(c, 502, "Failed to parse geocoding response")
		return
	}

	suggestions := make([]PlaceSuggestion, 0, len(results))
	for _, r := range results {
		lat, err1 := strconv.ParseFloat(r.Lat, 64)
		lon, err2 := strconv.ParseFloat(r.Lon, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		title, subtitle := splitTitleSubtitle(r.DisplayName)
		suggestions = append(suggestions, PlaceSuggestion{
			Lat:         lat,
			Lon:         lon,
			Title:       title,
			Subtitle:    subtitle,
			DisplayName: r.DisplayName,
		})
	}

	utils.SuccessResponse(c, 200, "Suggestions retrieved", gin.H{
		"count":       len(suggestions),
		"suggestions": suggestions,
	})
}
