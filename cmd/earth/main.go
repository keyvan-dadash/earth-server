package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/sod-lol/earth-server/internal/app/authentication"
	"github.com/sod-lol/earth-server/pkg/redis"
)

func checkDB() {
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// insert a tweet
	if err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").Exec(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	checkDB()

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
