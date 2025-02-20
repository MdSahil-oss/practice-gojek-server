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
var ErrNullUser error = fmt.Errorf("%d - null user found", http.StatusNotFound)

type postCreateUserContent struct {
	Name string `json:"name" binding:"required"`
}

func main() {
	fmt.Println("Hello Gojek")

	// Create an API Endpoint.
	// * Recognize the user, Then request-count of the user in the time-limit.
	// * If the user didn't exceed their request-limit then return 200 response, Otherwise return HTTP 429.
	router := gin.Default()
	requestInfos = validateuseraccess.New(100, time.Minute)

	router.GET("/", getRootRequest)
	router.GET("/get-users", getUsersRequest)
	router.POST("/create-user", createUser)

	router.Run("localhost:8080")

	// On that endpoint enable go-routine. So that multiple users can iteract with API endpoint at once.
}

func getRootRequest(ctx *gin.Context) {
	userId := ctx.GetHeader("userId")
	if userId == "" {
		ctx.String(http.StatusNotFound, ErrNullUser.Error())
		return
	}

	// Store UserID -> RequestCount
	if err := requestInfos.ValidateUserAccess(userId); err != nil {
		ctx.Error(err)
		if err == validateuseraccess.ErrUserNotExist {
			ctx.String(http.StatusNotFound, err.Error())
		} else {
			ctx.String(http.StatusTooManyRequests, err.Error())
		}
		return
	}

	ctx.String(http.StatusAccepted, "Successful! root request\n")
}

// Implement below one for Post method: https://spdeepak.hashnode.dev/golang-gin-tutorial-3
func createUser(ctx *gin.Context) {
	var content postCreateUserContent
	ctx.ShouldBindJSON(&content)

	if !requestInfos.IsUserExist(content.Name) {
		requestInfos.AddNewUser(content.Name)
		ctx.String(http.StatusCreated, "Successful! User created\n")
		return
	}

	err := fmt.Errorf("%d - Failed! User already exist", http.StatusAlreadyReported)
	ctx.Error(err)
	ctx.String(http.StatusAlreadyReported, err.Error())
}

func getUsersRequest(ctx *gin.Context) {
	userId := ctx.GetHeader("userId")
	if userId == "" {
		ctx.String(http.StatusNotFound, ErrNullUser.Error())
		return
	}

	// Store UserID -> RequestCount
	if err := requestInfos.ValidateUserAccess(userId); err != nil {
		ctx.Error(err)
	}

	byteData, err := json.Marshal(requestInfos.UserRequestMap)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.String(http.StatusAccepted, string(byteData)+"\n")
}
