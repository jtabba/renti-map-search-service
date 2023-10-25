package main

import (
	propertiesControllerV1 "back-end/mapSearchService/src/api/v1/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cors := cors.Default()
	router := gin.Default()
	router.Use(cors)

	mapSearchApiV1 := router.Group("/api/v1")
	{
		mapSearchApiV1.GET("/find-properties", propertiesControllerV1.GetScheduledListings)
	}

	router.Run(":3001")
}