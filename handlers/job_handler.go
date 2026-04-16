package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"tradeops-jobs-api/models"
	"tradeops-jobs-api/store"
)

type JobHandler struct {
	store *store.JobStore
}

func NewJobHandler(s *store.JobStore) *JobHandler {
	return &JobHandler{store: s}
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	var req models.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	job, _ := h.store.Create(req)
	c.JSON(http.StatusCreated, gin.H{"message": "Job created", "job": job})
}

func (h *JobHandler) ListJobs(c *gin.Context) {
	jobs := h.store.List(c.Query("status"))
	c.JSON(http.StatusOK, gin.H{"count": len(jobs), "jobs": jobs})
}

func (h *JobHandler) UpdateJobStatus(c *gin.Context) {
	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	job, err := h.store.UpdateStatus(c.Param("id"), req)
	if err != nil {
		if errors.Is(err, store.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status updated", "job": job})
}

func (h *JobHandler) CloseJob(c *gin.Context) {
	var req models.CloseJobRequest
	_ = c.ShouldBindJSON(&req)
	job, err := h.store.Close(c.Param("id"), req)
	if err != nil {
		if errors.Is(err, store.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job closed", "job": job})
}