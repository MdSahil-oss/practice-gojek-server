package main

import (
	"fmt"
	validateuseraccess "gojek-first/validate_user_access"
	"net/http"
	"time"

	"encoding/json"

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
	router.GET("/get-users", getUsersRequest)
	// router.POST("/create-user", createUser)

	router.Run("localhost:8080")

	// On that endpoint enable go-routine. So that multiple users can iteract with API endpoint at once.
}

func getRootRequest(ctx *gin.Context) {
	if ctx.GetHeader("userId") == "" {
		ctx.String(http.StatusNotFound, "HTTP 404 - User Not Found\n")
		return
	}

	// Store UserID -> RequestCount
	if requestInfos.ValidateUserAccess(ctx) {
		ctx.String(http.StatusAccepted, "Successful! root request\n")
		return
	}

	ctx.String(http.StatusTooManyRequests, "HTTP 429 - Too Many Requests\n")
}

// Implement below one for Post method: https://spdeepak.hashnode.dev/golang-gin-tutorial-3
// func createUser(ctx *gin.Context) {
// 	var content string
// 	// ctx.ShouldBindString(&content)
// 	fmt.Println("Null for now")

// 	ctx.String(http.StatusTooManyRequests, "On /data endpoint\n")
// }

func getUsersRequest(ctx *gin.Context) {
	if ctx.GetHeader("userId") == "" {
		ctx.String(http.StatusNotFound, "HTTP 404 - User Not Found\n")
		return
	}

	// Store UserID -> RequestCount
	if requestInfos.ValidateUserAccess(ctx) {
		byteData, err := json.Marshal(requestInfos.UserRequestMap)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusAccepted, string(byteData)+"\n")
		return
	}

	ctx.String(http.StatusTooManyRequests, "HTTP 429 - Too Many Requests\n")
}
