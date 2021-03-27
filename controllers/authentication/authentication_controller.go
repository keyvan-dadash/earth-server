package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sod-lol/earth-server/core/models/user"
	"github.com/sod-lol/earth-server/middlewares/token"
	"github.com/sod-lol/earth-server/services/redis"
)

type loginJsonExpect struct {
	Username string `form:"username" json:"username" xml:"username"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func HandleLogin(redisDB *redis.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authJson loginJsonExpect

		if err := c.ShouldBindJSON(&authJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		retrivedUser := user.User{
			Username: authJson.Username,
		}

		if err := user.UserRepository.RetrieveUser(&retrivedUser); err != nil {
			logrus.Errorf("Cannot retrive user. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		if !retrivedUser.VerifyPassword(authJson.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "username or password is incorrect",
			})

			return
		}

		genratedToken, err := token.CreateToken(retrivedUser.Username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = token.SaveTokenDetail(redisDB, genratedToken, retrivedUser.Username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access":  genratedToken.GetAccessToken(),
			"refresh": genratedToken.GetRefreshToken(),
		})

	}
}
