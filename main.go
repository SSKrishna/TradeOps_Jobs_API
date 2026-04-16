package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"tradeops-jobs-api/handlers"
	"tradeops-jobs-api/store"
)

func main() {
	r := gin.Default()
	jobStore := store.NewJobStore()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "TradeOps Jobs API"})
	})

	v1 := r.Group("/api/v1")
	h := handlers.NewJobHandler(jobStore)
	v1.POST("/jobs", h.CreateJob)
	v1.GET("/jobs", h.ListJobs)
	v1.PATCH("/jobs/:id/status", h.UpdateJobStatus)
	v1.POST("/jobs/:id/close", h.CloseJob)

	log.Println("Server running on :8080")
	r.Run(":8080")
}