package validate_user_access

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RequestInfos struct {
	rateLimit  uint
	timeWindow time.Duration

	userRequestMap map[string]userRequest
}

type userRequest struct {
	requestCount uint
	lastTime     time.Time
}

func New(ratelimit uint, tempTimeWindow time.Duration) RequestInfos {
	return RequestInfos{
		rateLimit:      ratelimit,
		timeWindow:     tempTimeWindow,
		userRequestMap: make(map[string]userRequest),
	}
}

func (ri RequestInfos) AddNewUser(userId string) {
	ri.userRequestMap[userId] = userRequest{
		requestCount: 0,
		lastTime:     time.Now(),
	}
}

// Validates the user's ratelimit
func (ri RequestInfos) CanUserAccess(c *gin.Context) bool {
	// * Recognize the user, Then request-count of the user within the time-limit.
	userId := c.Request.Header.Get("userId")

	// Check if user already exist otherise store user info
	userInfo, ok := ri.userRequestMap[userId]
	if !ok {
		ri.AddNewUser(userId)
		userInfo = ri.userRequestMap[userId]
	}

	// Put a condition to refresh userInfo.lastTime on over timeWindow
	if time.Since(userInfo.lastTime) > ri.timeWindow {
		userInfo.lastTime = time.Now()
		userInfo.requestCount = 0
	}

	// Store UserID -> RequestCount
	if userInfo.requestCount < ri.rateLimit {
		userInfo.requestCount++
		ri.userRequestMap[userId] = userInfo
		return true
	}

	return false
}
