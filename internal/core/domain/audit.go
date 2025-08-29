package domain

import (
	"time"
)

// AuditEvent represents a single action performed by a user.
type AuditEvent struct {
	ID           string
	UserID       string
	Action       string
	Timestamp    time.Time
	Details      string
	ResourceType string
	ResourceID   string
}

// FilterParams defines the criteria for filtering audit events.
type FilterParams struct {
	UserID       *string
	Action       *string
	ResourceType *string
	ResourceID   *string
	FromDate     *time.Time
	ToDate       *time.Time
}

// Pagination defines the limit and offset for retrieving data.
type Pagination struct {
	Limit  int32
	Offset int32
}