package main

// Rate Limiter Problem Statement
// Design and implement a rate-limiting service for an API. The service should:
// 1. Limit the number of API requests a user can make in a fixed time window (e.g., 100 requests per minute).
// 2. Return an appropriate response (e.g., HTTP 429 - Too Many Requests) when the limit is exceeded.
// 3. Support multiple users, each with their own independent rate limit tracking.

// Key Considerations:
// Use in-memory data structures to track requests.
// Ensure efficiency and scalability.
// Handle edge cases like `burst traffic` and `time window` overlaps.

// Extensions:
// Make the time window and request limit configurable.
// Discuss how this could scale in a distributed system.

/*
Future plan:
* Make validate_user_access package initialigible to accept values UserID, RequestCount, TimeWindow.
* Create multiple endpoints for diffrent type of requests (e.g. GET, POST)
*/

import (
	"fmt"
	validateuseraccess "gojek-first/validate_user_access"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var requestInfos validateuseraccess.RequestInfos

func main() {
	fmt.Println("Hello Gojek")

	// Create an API Endpoint.
	// * Recognize the user, Then request-count of the user in the time-limit.
	// * If the user didn't exceed their request-limit then return 200 response, Otherwise return HTTP 429.
	router := gin.Default()
	requestInfos = validateuseraccess.New(100, time.Minute)

	router.GET("/", getRootRequest)
	router.Run("localhost:8080")

	// On that endpoint enable go-routine. So that multiple users can iteract with API endpoint at once.
}

func getRootRequest(c *gin.Context) {
	// validateuseraccess.CanUserAccess(c)

	// Store UserID -> RequestCount
	if requestInfos.CanUserAccess(c) {
		c.String(http.StatusAccepted, "Successful! root request")
		return
	}

	c.String(http.StatusTooManyRequests, "HTTP 429 - Too Many Requests")
}
