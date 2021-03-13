package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sod-lol/earth-server/services/redis"
)

//HandleAuthenticationApp is function that handle authentication of earth
func HandleAuthenticationApp(authRouter *gin.RouterGroup, redis *redis.Redis) {
	authRouter.GET("/hel", handleAuth)
}

func handleAuth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
