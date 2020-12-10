package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sod-lol/earth-server/internal/app/authentication"
	"github.com/sod-lol/earth-server/pkg/redis"
)

func main() {

	redisAuth := redis.CreateRedisClient("redis-auth:6379", "", 0)

	earth := createEarthRouter()

	earth.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authentication.HandleAuthenticationApp(earth.Group("/auth"), redisAuth)
	earth.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
