package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sod-lol/earth-server/internal/app/authentication"
)

func main() {

	earth := createEarthRouter()

	earth.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authentication.HandleAuthenticationApp(earth.Group("/auth"))
	earth.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
