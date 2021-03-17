package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/sirupsen/logrus"
	"github.com/sod-lol/earth-server/config"
	"github.com/sod-lol/earth-server/middlewares/token"
	"github.com/sod-lol/earth-server/routers"
	"github.com/sod-lol/earth-server/services/redis"
)

func connectToDB() *gocqlx.Session {
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "earth"
	cluster.Consistency = gocql.One

	session, err := gocqlx.WrapSession(cluster.CreateSession())

	if err != nil {
		logrus.Fatalf("Cannot create session and connect to database. err: %v", err)
		return nil
	}

	return &session
}

func configAndSetupDB(session *gocqlx.Session) (bool, error) {
	if _, err := config.InitializeTables(session); err != nil {
		logrus.Fatalf("[Fatal](configAndSetupDB) Something went wrong during initialize tables. error: %v", err)
		return false, err
	}

	return true, nil
}

func main() {

	session := connectToDB()
	defer session.Close()

	if _, err := configAndSetupDB(session); err != nil {
		logrus.Fatal("[Fatal](main) terminate program due to hot error during initialize tables. error: %v", err)
		return
	}

	redisAuth := redis.CreateRedisClient("redis-auth:6379", "", 0)

	earth := createEarthRouter()

	earth.GET("/ping", token.TokenMiddleWareAuth(redisAuth), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routers.HandleAuthenticationApp(earth.Group("/auth"), redisAuth)
	earth.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
