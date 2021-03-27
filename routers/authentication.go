package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sod-lol/earth-server/controllers/authentication"
	"github.com/sod-lol/earth-server/services/redis"
)

//HandleAuthenticationApp is function that handle authentication of earth
func HandleAuthenticationApp(ctx context.Context, authRouter *gin.RouterGroup) {

	redisDB := ctx.Value("redisDB").(*redis.Redis)

	authRouter.POST("/login", authentication.HandleLogin(redisDB))
	authRouter.POST("/signup", authentication.HandleSignUp())
}
