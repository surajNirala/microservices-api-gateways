package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func reverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the target microservice URL
		remote, err := url.Parse(target)
		if err != nil {
			log.Fatalf("Could not parse service URL: %v", err)
		}
		// Create a reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(remote)
		// originalPath := c.Param("proxyPath")
		fullPath := c.Request.URL.Path
		// fmt.Println("before fullPath : ", fullPath)
		fullPath = strings.TrimSuffix(fullPath, "/")
		// fmt.Println("After fullPath : ", fullPath)
		c.Request.URL.Path = fullPath
		c.Request.Host = remote.Host
		c.Request.URL.Scheme = remote.Scheme
		c.Request.Header.Set("X-Forwarded-Host", c.Request.Host)
		// Log the captured and trimmed path
		log.Printf("Forwarding request to %s%s\n", remote, fullPath)
		// Serve the request via the proxy
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	/* err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Set Gin mode based on environment
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	// Get APP_PORT and APP_URL from environment variables
	app_port := os.Getenv("APP_PORT")
	app_url := os.Getenv("APP_URL")
	if app_port == "" {
		app_port = "8000"
	}
	if app_url == "" {
		app_url = "localhost"
	}
	user_service := os.Getenv("USER_SERVICE")
	rating_service := os.Getenv("RATING_SERVICE")
	hotel_service := os.Getenv("HOTEL_SERVICE") */

	user_service := "http://34.131.139.0:9091"
	rating_service := "http://34.131.139.0:9092"
	hotel_service := "http://34.131.139.0:9093"

	r := gin.Default()

	// Define microservice URLs
	userServiceURL := user_service
	ratingServiceURL := rating_service
	hotelServiceURL := hotel_service

	// Route for User Service
	r.Any("/api/users/*proxyPath", reverseProxy(userServiceURL))

	// Route for Rating Service
	r.Any("/api/ratings/*proxyPath", reverseProxy(ratingServiceURL))

	// Route for Hotel Service
	r.Any("/api/hotels/*proxyPath", reverseProxy(hotelServiceURL))

	// Start the API Gateway on port 9094
	if err := r.Run(":9094"); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
