package services

import (
	"locker-storage-api/internal/models"
	"math"
)

// Haversine formula to calculate distance between two points (in kilometers)
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

type LocationWithDistance struct {
	models.Location
	Distance         float64 `json:"distance"`
	AvailableLockers int     `json:"availableLockers"`
}

func FindNearbyLocations(userLat, userLon float64, radiusKm float64, locations []models.Location) []LocationWithDistance {
	var nearby []LocationWithDistance

	for _, loc := range locations {
		distance := CalculateDistance(userLat, userLon, loc.Latitude, loc.Longitude)
		if distance <= radiusKm {
			nearby = append(nearby, LocationWithDistance{
				Location: loc,
				Distance: math.Round(distance*10) / 10, // Round to 1 decimal
			})
		}
	}

	// Sort by distance (bubble sort for simplicity)
	for i := 0; i < len(nearby); i++ {
		for j := i + 1; j < len(nearby); j++ {
			if nearby[i].Distance > nearby[j].Distance {
				nearby[i], nearby[j] = nearby[j], nearby[i]
			}
		}
	}

	return nearby
}
