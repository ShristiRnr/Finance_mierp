package domain

import (
	"time"

	"github.com/google/uuid"
)

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