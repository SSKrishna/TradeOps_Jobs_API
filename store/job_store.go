package store

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"tradeops-jobs-api/models"
)

var (
	ErrJobNotFound = errors.New("job not found")
	ErrJobClosed   = errors.New("job is already closed")
)

type JobStore struct {
	mu   sync.RWMutex
	jobs map[string]*models.Job
}

func NewJobStore() *JobStore {
	s := &JobStore{jobs: make(map[string]*models.Job)}
	s.seed()
	return s
}

func (s *JobStore) seed() {
	now := time.Now()
	samples := []*models.Job{
		{
			ID:           uuid.NewString(),
			Title:        "Burst pipe under kitchen sink",
			CustomerName: "Amanda Torres",
			Address:      "1420 Maple St, Omaha NE 68102",
			TradeType:    "plumbing",
			Status:       models.StatusDispatched,
			AssignedTo:   "Mike Jimenez",
			CreatedAt:    now.Add(-2 * time.Hour),
			UpdatedAt:    now.Add(-90 * time.Minute),
		},
		{
			ID:           uuid.NewString(),
			Title:        "HVAC unit not cooling",
			CustomerName: "Derek Schultz",
			Address:      "880 Pine Ave, Omaha NE 68106",
			TradeType:    "hvac",
			Status:       models.StatusInProgress,
			AssignedTo:   "Sara Kim",
			CreatedAt:    now.Add(-5 * time.Hour),
			UpdatedAt:    now.Add(-30 * time.Minute),
		},
		{
			ID:           uuid.NewString(),
			Title:        "Breaker tripping on circuit 4",
			CustomerName: "Patricia Webb",
			Address:      "305 Oak Blvd, Omaha NE 68132",
			TradeType:    "electrical",
			Status:       models.StatusPending,
			CreatedAt:    now.Add(-20 * time.Minute),
			UpdatedAt:    now.Add(-20 * time.Minute),
		},
	}
	for _, j := range samples {
		s.jobs[j.ID] = j
	}
}

func (s *JobStore) Create(req models.CreateJobRequest) (*models.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	job := &models.Job{
		ID:           uuid.NewString(),
		Title:        req.Title,
		Description:  req.Description,
		CustomerName: req.CustomerName,
		Address:      req.Address,
		TradeType:    req.TradeType,
		Status:       models.StatusPending,
		AssignedTo:   req.AssignedTo,
		Notes:        req.Notes,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	s.jobs[job.ID] = job
	return job, nil
}

func (s *JobStore) List(statusFilter string) []*models.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*models.Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		if statusFilter == "" || string(j.Status) == statusFilter {
			result = append(result, j)
		}
	}
	return result
}

func (s *JobStore) UpdateStatus(id string, req models.UpdateStatusRequest) (*models.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[id]
	if !ok {
		return nil, ErrJobNotFound
	}
	if job.Status == models.StatusClosed {
		return nil, ErrJobClosed
	}
	job.Status = req.Status
	job.UpdatedAt = time.Now()
	if req.Notes != "" {
		job.Notes = req.Notes
	}
	return job, nil
}

func (s *JobStore) Close(id string, req models.CloseJobRequest) (*models.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[id]
	if !ok {
		return nil, ErrJobNotFound
	}
	if job.Status == models.StatusClosed {
		return nil, ErrJobClosed
	}
	now := time.Now()
	job.Status = models.StatusClosed
	job.ClosedAt = &now
	job.UpdatedAt = now
	if req.Notes != "" {
		job.Notes = req.Notes
	}
	return job, nil
}