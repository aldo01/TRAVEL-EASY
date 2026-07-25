package main

import (
	"fmt"
	"log"
	"time"

	"locker-storage-api/config"
	"locker-storage-api/internal/models"
	"locker-storage-api/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	strPtr := func(s string) *string { return &s }

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

	// Create test/demo users
	createUser := func(email, name, password string, phone *string) {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			log.Println("✗ Failed to hash password for", email, ":", err)
			return
		}

		user := models.User{
			Email:       email,
			Password:    hashedPassword,
			Name:        name,
			PhoneNumber: phone,
			IsActive:    true,
		}

		// FirstOrCreate keeps seeding idempotent.
		if err := db.FirstOrCreate(&user, models.User{Email: user.Email}).Error; err != nil {
			log.Println("✗ Failed to create user:", email, ":", err)
			return
		}
		log.Println("✓ Ensured demo user:", email)
	}

	createUser("test@example.com", "John Doe", "password123", strPtr("+4512345678"))
	createUser("alice@example.com", "Alice Johnson", "password123", strPtr("+14155550101"))
	createUser("bob@example.com", "Bob Singh", "password123", strPtr("+14155550102"))
	createUser("neha@demo.com", "Neha Sharma", "password123", strPtr("+919900000001"))
	createUser("rahul@demo.com", "Rahul Verma", "password123", strPtr("+919900000002"))
	createUser("host@demo.com", "Demo Host", "password123", strPtr("+919900000003"))

	// Load the main test user for sample booking creation
	var testUser models.User
	if err := db.Where("email = ?", "test@example.com").First(&testUser).Error; err != nil {
		log.Fatal("Failed to load test user after creation:", err)
	}

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
			Description: strPtr("Main train station with 24/7 locker access"),
			ImageURL:    strPtr("https://images.unsplash.com/photo-1581092918056-0c4c3acd3789"),
			IsOpen24Hours: false,
			Amenities:     "[\"CCTV\",\"Staffed\",\"Near transit\"]",
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
			Description: strPtr("Modern gym with secure locker storage"),
			ImageURL:    strPtr("https://images.unsplash.com/photo-1534438327276-14e5300c3a48"),
			IsOpen24Hours: false,
			Amenities:     "[\"CCTV\",\"Indoor\"]",
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
			Description: strPtr("Student club with affordable storage"),
			ImageURL:    strPtr("https://images.unsplash.com/photo-1523050854058-8df90110c9f1"),
			IsOpen24Hours: false,
			Amenities:     "[\"Budget\",\"Staffed\"]",
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
			Description: strPtr("Prime tourist location near colorful harbor"),
			ImageURL:    strPtr("https://images.unsplash.com/photo-1513622470522-26c3c8a854bc"),
			IsOpen24Hours: false,
			Amenities:     "[\"Near attraction\",\"CCTV\"]",
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
			Description: strPtr("24/7 airport locker facility"),
			ImageURL:    strPtr("https://images.unsplash.com/photo-1436491865332-7a61a109cc05"),
			IsOpen24Hours: true,
			Amenities:     "[\"24/7\",\"CCTV\",\"Near airport\"]",
			IsActive:    true,
		},
		// India (Bengaluru)
		{
			Name:          "MG Road Coffee Storage",
			Address:       "MG Road, Bengaluru",
			City:          "Bengaluru",
			Country:       "India",
			Type:          models.LocationCafe,
			Latitude:      12.9756,
			Longitude:     77.6068,
			HourlyRate:    60.0,
			DailyRate:     250.0,
			OpeningTime:   "08:00",
			ClosingTime:   "22:00",
			IsOpen24Hours: false,
			Rating:        4.4,
			Description:   strPtr("Cafe partner with secure counter storage"),
			Amenities:     "[\"Staffed\",\"Indoor\"]",
			IsActive:      true,
		},
		{
			Name:          "Indiranagar Boutique Hotel Lockers",
			Address:       "Indiranagar, Bengaluru",
			City:          "Bengaluru",
			Country:       "India",
			Type:          models.LocationHotel,
			Latitude:      12.9719,
			Longitude:     77.6412,
			HourlyRate:    80.0,
			DailyRate:     300.0,
			OpeningTime:   "00:00",
			ClosingTime:   "23:59",
			IsOpen24Hours: true,
			Rating:        4.6,
			Description:   strPtr("Hotel lobby lockers with 24/7 access"),
			Amenities:     "[\"24/7\",\"CCTV\"]",
			IsActive:      true,
		},
		{
			Name:          "Majestic Transit Storage",
			Address:       "Kempegowda Bus Station (Majestic), Bengaluru",
			City:          "Bengaluru",
			Country:       "India",
			Type:          models.LocationKiosk,
			Latitude:      12.9767,
			Longitude:     77.5727,
			HourlyRate:    55.0,
			DailyRate:     220.0,
			OpeningTime:   "06:00",
			ClosingTime:   "23:00",
			IsOpen24Hours: false,
			Rating:        4.3,
			Description:   strPtr("Kiosk storage near the main transit hub"),
			Amenities:     "[\"Near transit\",\"Staffed\",\"CCTV\"]",
			IsActive:      true,
		},
		{
			Name:          "Koramangala Store Lockers",
			Address:       "Koramangala 5th Block, Bengaluru",
			City:          "Bengaluru",
			Country:       "India",
			Type:          models.LocationStore,
			Latitude:      12.9346,
			Longitude:     77.6163,
			HourlyRate:    65.0,
			DailyRate:     240.0,
			OpeningTime:   "09:00",
			ClosingTime:   "22:00",
			IsOpen24Hours: false,
			Rating:        4.5,
			Description:   strPtr("Convenience store partner with secure backroom storage"),
			Amenities:     "[\"Staffed\",\"Indoor\"]",
			IsActive:      true,
		},
		{
			Name:          "Bengaluru Airport 24/7 Storage",
			Address:       "Kempegowda International Airport, Bengaluru",
			City:          "Bengaluru",
			Country:       "India",
			Type:          models.LocationKiosk,
			Latitude:      13.1986,
			Longitude:     77.7066,
			HourlyRate:    90.0,
			DailyRate:     350.0,
			OpeningTime:   "00:00",
			ClosingTime:   "23:59",
			IsOpen24Hours: true,
			Rating:        4.6,
			Description:   strPtr("Airport storage facility with 24/7 access"),
			Amenities:     "[\"24/7\",\"Near airport\",\"CCTV\"]",
			IsActive:      true,
		},
		// India (Delhi)
		{
			Name:          "Connaught Place Storage Kiosk",
			Address:       "Connaught Place, New Delhi",
			City:          "New Delhi",
			Country:       "India",
			Type:          models.LocationKiosk,
			Latitude:      28.6315,
			Longitude:     77.2167,
			HourlyRate:    70.0,
			DailyRate:     280.0,
			OpeningTime:   "09:00",
			ClosingTime:   "21:00",
			IsOpen24Hours: false,
			Rating:        4.2,
			Description:   strPtr("Central kiosk near metro"),
			Amenities:     "[\"Near transit\",\"Staffed\"]",
			IsActive:      true,
		},
		{
			Name:          "New Delhi Railway Station Lockers",
			Address:       "New Delhi Railway Station, New Delhi",
			City:          "New Delhi",
			Country:       "India",
			Type:          models.LocationKiosk,
			Latitude:      28.6448,
			Longitude:     77.2167,
			HourlyRate:    65.0,
			DailyRate:     260.0,
			OpeningTime:   "06:00",
			ClosingTime:   "23:00",
			IsOpen24Hours: false,
			Rating:        4.1,
			Description:   strPtr("Station partner storage close to platforms"),
			Amenities:     "[\"Near transit\",\"CCTV\",\"Staffed\"]",
			IsActive:      true,
		},
		{
			Name:          "Karol Bagh Market Storage",
			Address:       "Karol Bagh, New Delhi",
			City:          "New Delhi",
			Country:       "India",
			Type:          models.LocationStore,
			Latitude:      28.6513,
			Longitude:     77.1906,
			HourlyRate:    55.0,
			DailyRate:     220.0,
			OpeningTime:   "10:00",
			ClosingTime:   "22:00",
			IsOpen24Hours: false,
			Rating:        4.3,
			Description:   strPtr("Shop partner near the market area"),
			Amenities:     "[\"Indoor\",\"Staffed\"]",
			IsActive:      true,
		},
		{
			Name:          "Hauz Khas Village Storage",
			Address:       "Hauz Khas, New Delhi",
			City:          "New Delhi",
			Country:       "India",
			Type:          models.LocationCafe,
			Latitude:      28.5535,
			Longitude:     77.1940,
			HourlyRate:    60.0,
			DailyRate:     240.0,
			OpeningTime:   "09:00",
			ClosingTime:   "23:00",
			IsOpen24Hours: false,
			Rating:        4.4,
			Description:   strPtr("Cafe partner in the Hauz Khas area"),
			Amenities:     "[\"Staffed\",\"Indoor\"]",
			IsActive:      true,
		},
		{
			Name:          "IGI Airport Terminal Storage",
			Address:       "Indira Gandhi International Airport, New Delhi",
			City:          "New Delhi",
			Country:       "India",
			Type:          models.LocationKiosk,
			Latitude:      28.5562,
			Longitude:     77.1000,
			HourlyRate:    95.0,
			DailyRate:     380.0,
			OpeningTime:   "00:00",
			ClosingTime:   "23:59",
			IsOpen24Hours: true,
			Rating:        4.6,
			Description:   strPtr("24/7 airport luggage storage"),
			Amenities:     "[\"24/7\",\"Near airport\",\"CCTV\"]",
			IsActive:      true,
		},
		// India (Mumbai)
		{
			Name:          "CST Station Luggage Storage",
			Address:       "Chhatrapati Shivaji Maharaj Terminus, Mumbai",
			City:          "Mumbai",
			Country:       "India",
			Type:          models.LocationKiosk,
			Latitude:      18.9402,
			Longitude:     72.8356,
			HourlyRate:    70.0,
			DailyRate:     280.0,
			OpeningTime:   "06:00",
			ClosingTime:   "23:00",
			IsOpen24Hours: false,
			Rating:        4.2,
			Description:   strPtr("Station partner storage near CST"),
			Amenities:     "[\"Near transit\",\"Staffed\",\"CCTV\"]",
			IsActive:      true,
		},
		{
			Name:          "Bandra Cafe Storage",
			Address:       "Bandra West, Mumbai",
			City:          "Mumbai",
			Country:       "India",
			Type:          models.LocationCafe,
			Latitude:      19.0600,
			Longitude:     72.8296,
			HourlyRate:    65.0,
			DailyRate:     250.0,
			OpeningTime:   "08:00",
			ClosingTime:   "22:30",
			IsOpen24Hours: false,
			Rating:        4.5,
			Description:   strPtr("Cafe partner in Bandra"),
			Amenities:     "[\"Staffed\",\"Indoor\"]",
			IsActive:      true,
		},
		{
			Name:          "Andheri Metro Storage Kiosk",
			Address:       "Andheri, Mumbai",
			City:          "Mumbai",
			Country:       "India",
			Type:          models.LocationKiosk,
			Latitude:      19.1197,
			Longitude:     72.8468,
			HourlyRate:    60.0,
			DailyRate:     240.0,
			OpeningTime:   "07:00",
			ClosingTime:   "23:00",
			IsOpen24Hours: false,
			Rating:        4.3,
			Description:   strPtr("Kiosk storage close to metro access"),
			Amenities:     "[\"Near transit\",\"Staffed\"]",
			IsActive:      true,
		},
		{
			Name:          "Colaba Causeway Store Lockers",
			Address:       "Colaba Causeway, Mumbai",
			City:          "Mumbai",
			Country:       "India",
			Type:          models.LocationStore,
			Latitude:      18.9217,
			Longitude:     72.8331,
			HourlyRate:    75.0,
			DailyRate:     300.0,
			OpeningTime:   "10:00",
			ClosingTime:   "22:00",
			IsOpen24Hours: false,
			Rating:        4.4,
			Description:   strPtr("Tourist area partner with secure storage"),
			Amenities:     "[\"Near attraction\",\"Staffed\",\"Indoor\"]",
			IsActive:      true,
		},
		{
			Name:          "Mumbai Airport 24/7 Storage",
			Address:       "Chhatrapati Shivaji Maharaj International Airport, Mumbai",
			City:          "Mumbai",
			Country:       "India",
			Type:          models.LocationKiosk,
			Latitude:      19.0896,
			Longitude:     72.8656,
			HourlyRate:    95.0,
			DailyRate:     380.0,
			OpeningTime:   "00:00",
			ClosingTime:   "23:59",
			IsOpen24Hours: true,
			Rating:        4.6,
			Description:   strPtr("24/7 airport storage facility"),
			Amenities:     "[\"24/7\",\"Near airport\",\"CCTV\"]",
			IsActive:      true,
		},
		// UK (London)
		{
			Name:          "King's Cross Luggage Storage",
			Address:       "King's Cross, London",
			City:          "London",
			Country:       "United Kingdom",
			Type:          models.LocationStore,
			Latitude:      51.5308,
			Longitude:     -0.1238,
			HourlyRate:    5.0,
			DailyRate:     20.0,
			OpeningTime:   "07:00",
			ClosingTime:   "23:00",
			IsOpen24Hours: false,
			Rating:        4.6,
			Description:   strPtr("Convenience store partner near station"),
			Amenities:     "[\"CCTV\",\"Near transit\"]",
			IsActive:      true,
		},
		// USA (NYC)
		{
			Name:          "Times Square Storage Shop",
			Address:       "Times Square, New York, NY",
			City:          "New York",
			Country:       "USA",
			Type:          models.LocationStore,
			Latitude:      40.7580,
			Longitude:     -73.9855,
			HourlyRate:    6.0,
			DailyRate:     22.0,
			OpeningTime:   "08:00",
			ClosingTime:   "22:00",
			IsOpen24Hours: false,
			Rating:        4.5,
			Description:   strPtr("Retail partner with secure backroom storage"),
			Amenities:     "[\"Staffed\",\"Near attraction\"]",
			IsActive:      true,
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
					LockerNumber: fmt.Sprintf("%d", i+1),
					Size:         size,
					Status:       models.LockerAvailable,
					IsOperational: true,
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
		PaymentStatus: models.PaymentPaid,
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
