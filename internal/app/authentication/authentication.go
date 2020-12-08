package authentication

import "github.com/gin-gonic/gin"

//HandleAuthenticationApp is function that handle authentication of earth
func HandleAuthenticationApp(authRouter *gin.RouterGroup) {

	authRouter.GET("/hel", handleAuth)
}

func handleAuth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
