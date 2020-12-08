package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

//EarthRouter is router for whole earch project
type EarthRouter struct {
	*gin.Engine
}

var (
	once     sync.Once
	instance *EarthRouter
)

func createEarthRouter() *EarthRouter {

	once.Do(func() {
		earth := gin.Default()
		instance = &EarthRouter{
			earth,
		}
	})

	return instance
}

//GetRouter is function that return gin router of EarthRouter
func GetRouter() *EarthRouter {
	return instance
}
