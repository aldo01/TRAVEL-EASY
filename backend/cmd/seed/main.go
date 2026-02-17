package main

import (
	"log"
	"time"

	"locker-storage-api/config"
	"locker-storage-api/internal/models"
	"locker-storage-api/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := config.GetDB()

	log.Println("Seeding database...")

	// Create test user
	hashedPassword, _ := utils.HashPassword("password123")
	testUser := models.User{
		Email:     "test@example.com",
		Password:  hashedPassword,
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "+4512345678",
	}
	db.FirstOrCreate(&testUser, models.User{Email: testUser.Email})
	log.Println("✓ Created test user:", testUser.Email)

	// Create locations
	locations := []models.Location{
		{
			Name:        "Copenhagen Central Station",
			Address:     "Bernstorffsgade 16, 1577 Copenhagen",
			City:        "Copenhagen",
			Country:     "Denmark",
			Type:        models.LocationKiosk,
			Latitude:    55.6730,
			Longitude:   12.5650,
			HourlyRate:  50.0,
			DailyRate:   300.0,
			OpeningTime: "06:00",
			ClosingTime: "23:00",
			Rating:      4.5,
			Description: "Main train station with 24/7 locker access",
			Images:      []string{"https://images.unsplash.com/photo-1581092918056-0c4c3acd3789"},
			IsActive:    true,
		},
		{
			Name:        "Fitness World Nørreport",
			Address:     "Frederiksborggade 15, 1360 Copenhagen",
			City:        "Copenhagen",
			Country:     "Denmark",
			Type:        models.LocationGym,
			Latitude:    55.6833,
			Longitude:   12.5725,
			HourlyRate:  40.0,
			DailyRate:   250.0,
			OpeningTime: "05:00",
			ClosingTime: "23:00",
			Rating:      4.3,
			Description: "Modern gym with secure locker storage",
			Images:      []string{"https://images.unsplash.com/photo-1534438327276-14e5300c3a48"},
			IsActive:    true,
		},
		{
			Name:        "Copenhagen University Club",
			Address:     "Nørregade 10, 1165 Copenhagen",
			City:        "Copenhagen",
			Country:     "Denmark",
			Type:        models.LocationClub,
			Latitude:    55.6804,
			Longitude:   12.5728,
			HourlyRate:  35.0,
			DailyRate:   200.0,
			OpeningTime: "08:00",
			ClosingTime: "20:00",
			Rating:      4.7,
			Description: "Student club with affordable storage",
			Images:      []string{"https://images.unsplash.com/photo-1523050854058-8df90110c9f1"},
			IsActive:    true,
		},
		{
			Name:        "Nyhavn Tourist Center",
			Address:     "Nyhavn 15, 1051 Copenhagen",
			City:        "Copenhagen",
			Country:     "Denmark",
			Type:        models.LocationKiosk,
			Latitude:    55.6795,
			Longitude:   12.5915,
			HourlyRate:  60.0,
			DailyRate:   350.0,
			OpeningTime: "07:00",
			ClosingTime: "22:00",
			Rating:      4.6,
			Description: "Prime tourist location near colorful harbor",
			Images:      []string{"https://images.unsplash.com/photo-1513622470522-26c3c8a854bc"},
			IsActive:    true,
		},
		{
			Name:        "Copenhagen Airport Storage",
			Address:     "Lufthavnsboulevarden 6, 2770 Kastrup",
			City:        "Copenhagen",
			Country:     "Denmark",
			Type:        models.LocationKiosk,
			Latitude:    55.6181,
			Longitude:   12.6508,
			HourlyRate:  70.0,
			DailyRate:   400.0,
			OpeningTime: "00:00",
			ClosingTime: "23:59",
			Rating:      4.8,
			Description: "24/7 airport locker facility",
			Images:      []string{"https://images.unsplash.com/photo-1436491865332-7a61a109cc05"},
			IsActive:    true,
		},
	}

	for _, loc := range locations {
		var existing models.Location
		if err := db.Where("name = ?", loc.Name).First(&existing).Error; err != nil {
			db.Create(&loc)
			log.Printf("✓ Created location: %s", loc.Name)

			// Create lockers for each location
			lockerSizes := []models.LockerSize{
				models.LockerSmall,
				models.LockerSmall,
				models.LockerSmall,
				models.LockerMedium,
				models.LockerMedium,
				models.LockerLarge,
			}

			for i, size := range lockerSizes {
				locker := models.Locker{
					LocationID:   loc.ID,
					LockerNumber: i + 1,
					Size:         size,
					Status:       models.LockerAvailable,
					IsActive:     true,
				}
				db.Create(&locker)
			}
			log.Printf("  ✓ Created %d lockers for %s", len(lockerSizes), loc.Name)
		}
	}

	// Create a sample booking
	var location models.Location
	db.First(&location)

	var locker models.Locker
	db.Where("location_id = ?", location.ID).First(&locker)

	startTime := time.Now().Add(2 * time.Hour)
	endTime := startTime.Add(24 * time.Hour)
	checkInTime := startTime.Add(30 * time.Minute)

	booking := models.Booking{
		UserID:        testUser.ID,
		LocationID:    location.ID,
		LockerID:      locker.ID,
		BookingNumber: "BK20240101001",
		StartTime:     startTime,
		EndTime:       &endTime,
		RateType:      models.RateDaily,
		BasePrice:     location.DailyRate,
		TotalPrice:    location.DailyRate,
		Status:        models.BookingActive,
		QRCode:        "QR-" + "test-booking-001",
		PaymentStatus: models.PaymentCompleted,
		CheckInTime:   &checkInTime,
	}

	var existingBooking models.Booking
	if err := db.Where("booking_number = ?", booking.BookingNumber).First(&existingBooking).Error; err != nil {
		db.Create(&booking)
		log.Println("✓ Created sample booking")

		// Update locker status
		db.Model(&locker).Update("status", models.LockerOccupied)
	}

	log.Println("✓ Database seeding completed successfully!")
	log.Println("\nTest credentials:")
	log.Println("Email: test@example.com")
	log.Println("Password: password123")
}
