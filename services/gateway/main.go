package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// route describes the target service and whether auth is required.
type route struct {
	target      string
	requireAuth bool
}

func resolveRoute(method, path string) (route, bool) {
	authURL := envOr("AUTH_SERVICE_URL", "http://auth-service:8001")
	locationURL := envOr("LOCATION_SERVICE_URL", "http://location-service:8002")
	bookingURL := envOr("BOOKING_SERVICE_URL", "http://booking-service:8003")
	reviewURL := envOr("REVIEW_SERVICE_URL", "http://review-service:8004")

	p := strings.TrimPrefix(path, "/api/v1")

	switch {
	case strings.HasPrefix(p, "/auth/"):
		return route{authURL, false}, true
	case p == "/profile":
		return route{authURL, true}, true
	case strings.HasPrefix(p, "/bookings"):
		return route{bookingURL, true}, true
	case strings.HasPrefix(p, "/reviews"):
		// GET is public, POST requires auth.
		return route{reviewURL, method != http.MethodGet}, true
	case strings.HasPrefix(p, "/locations"):
		// Creating a location (POST) requires auth; reads are public.
		return route{locationURL, method == http.MethodPost}, true
	case strings.HasPrefix(p, "/lockers"), strings.HasPrefix(p, "/geocode"):
		return route{locationURL, false}, true
	default:
		return route{}, false
	}
}

func proxyTo(target string, c *gin.Context) {
	u, err := url.Parse(target)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "error": "bad gateway target"})
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		log.Printf("[gateway] upstream error for %s: %v", target, e)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"success":false,"error":"Service unavailable"}`))
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func gatewayHandler(c *gin.Context) {
	r, ok := resolveRoute(c.Request.Method, c.Request.URL.Path)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Not found"})
		return
	}

	// Never trust client-supplied identity headers.
	c.Request.Header.Del("X-User-Id")
	c.Request.Header.Del("X-User-Email")

	if r.requireAuth {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Authorization header required"})
			return
		}
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))
		claims, err := validateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Invalid or expired token"})
			return
		}
		c.Request.Header.Set("X-User-Id", claims.UserID)
		c.Request.Header.Set("X-User-Email", claims.Email)
	}

	proxyTo(r.target, c)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	router := gin.Default()

	corsOrigin := envOr("CORS_ORIGIN", "http://localhost:5173")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin, "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "gateway"})
	})

	// Everything under /api/v1 is proxied to the owning service.
	router.NoRoute(gatewayHandler)

	port := envOr("PORT", "8080")
	log.Printf("gateway listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start gateway:", err)
	}
}
