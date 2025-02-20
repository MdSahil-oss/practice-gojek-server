package validate_user_access

import (
	"fmt"
	"net/http"
	"time"
)

var ErrTooManyRequest error = fmt.Errorf("%d - user exceeded request limit", http.StatusTooManyRequests)
var ErrUserNotExist error = fmt.Errorf("%d - user doesn't exist", http.StatusNotFound)

type RequestInfos struct {
	RateLimit  uint
	TimeWindow time.Duration

	UserRequestMap map[string]userRequest
}

type userRequest struct {
	RequestCount     uint      `json:"RequestCount"`
	SessionStartTime time.Time `json:"SessionStartTime"`
}

func New(RateLimit uint, tempTimeWindow time.Duration) RequestInfos {
	return RequestInfos{
		RateLimit:      RateLimit,
		TimeWindow:     tempTimeWindow,
		UserRequestMap: make(map[string]userRequest),
	}
}

func (ri RequestInfos) AddNewUser(userId string) {
	ri.UserRequestMap[userId] = userRequest{
		RequestCount:     0,
		SessionStartTime: time.Now(),
	}
}

// Validates the user's RateLimit
func (ri RequestInfos) ValidateUserAccess(userId string) error {
	// Check if user already exist otherise store user info
	userInfo, ok := ri.UserRequestMap[userId]
	if !ok {
		return ErrUserNotExist
	}

	// Put a condition to refresh userInfo.SessionStartTime on over TimeWindow
	if time.Since(userInfo.SessionStartTime) > ri.TimeWindow {
		userInfo.SessionStartTime = time.Now()
		userInfo.RequestCount = 0
	}

	// Store UserID -> RequestCount
	if userInfo.RequestCount > ri.RateLimit {
		return ErrTooManyRequest
	}

	userInfo.RequestCount++
	ri.UserRequestMap[userId] = userInfo
	return nil
}

func (ri RequestInfos) IsUserExist(userId string) bool {
	if _, ok := ri.UserRequestMap[userId]; !ok {
		return false
	}
	return true
}
