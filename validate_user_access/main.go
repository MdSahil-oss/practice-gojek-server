package validate_user_access

import (
	"time"

	"github.com/gin-gonic/gin"
)

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
func (ri RequestInfos) ValidateUserAccess(c *gin.Context) bool {
	// * Recognize the user, Then request-count of the user within the time-limit.
	userId := c.Request.Header.Get("userId")

	// Check if user already exist otherise store user info
	userInfo, ok := ri.UserRequestMap[userId]
	if !ok {
		ri.AddNewUser(userId)
		userInfo = ri.UserRequestMap[userId]
	}

	// Put a condition to refresh userInfo.SessionStartTime on over TimeWindow
	if time.Since(userInfo.SessionStartTime) > ri.TimeWindow {
		userInfo.SessionStartTime = time.Now()
		userInfo.RequestCount = 0
	}

	// Store UserID -> RequestCount
	if userInfo.RequestCount < ri.RateLimit {
		userInfo.RequestCount++
		ri.UserRequestMap[userId] = userInfo
		return true
	}

	return false
}
