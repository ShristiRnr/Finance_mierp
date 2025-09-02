package domain

import (
	"time"

	"github.com/google/uuid"
)	

// CashFlowForecast represents a forecast for cash inflows and outflows
type CashFlowForecast struct {
	ID             uuid.UUID
	OrganizationID string
	PeriodStart    time.Time
	PeriodEnd      time.Time
	ForecastDetails string
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedAt      time.Time
	UpdatedBy      string
	Revision       int32
}