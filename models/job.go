package models

import "time"

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusDispatched JobStatus = "dispatched"
	StatusInProgress JobStatus = "in_progress"
	StatusOnHold     JobStatus = "on_hold"
	StatusClosed     JobStatus = "closed"
)

type Job struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	CustomerName string     `json:"customer_name"`
	Address      string     `json:"address"`
	TradeType    string     `json:"trade_type"`
	Status       JobStatus  `json:"status"`
	AssignedTo   string     `json:"assigned_to,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	Notes        string     `json:"notes,omitempty"`
}

type CreateJobRequest struct {
	Title        string `json:"title"         binding:"required"`
	Description  string `json:"description"`
	CustomerName string `json:"customer_name"  binding:"required"`
	Address      string `json:"address"        binding:"required"`
	TradeType    string `json:"trade_type"     binding:"required"`
	AssignedTo   string `json:"assigned_to"`
	Notes        string `json:"notes"`
}

type UpdateStatusRequest struct {
	Status JobStatus `json:"status" binding:"required"`
	Notes  string    `json:"notes"`
}

type CloseJobRequest struct {
	Notes string `json:"notes"`
}