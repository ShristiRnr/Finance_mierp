package domain

import (
	"time"

	"github.com/google/uuid"
)
//Profit and Loss Report
type ProfitLossReport struct {
	ID            uuid.UUID
	OrganizationID string
	PeriodStart    time.Time
	PeriodEnd      time.Time
	TotalRevenue   float64
	TotalExpenses  float64
	NetProfit      float64
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedAt      time.Time
	UpdatedBy      string
	Revision       int32
}

// Balance Sheet Report
type BalanceSheetReport struct {
	ID              uuid.UUID
	OrganizationID   string
	PeriodStart      time.Time
	PeriodEnd        time.Time
	TotalAssets      float64
	TotalLiabilities float64
	NetWorth         float64
	CreatedAt        time.Time
	CreatedBy        string
	UpdatedAt        time.Time
	UpdatedBy        string
	Revision         int32
}

//Trial Balance Report
type TrialBalanceReport struct {
	ID             uuid.UUID
	OrganizationID string
	PeriodStart    time.Time
	PeriodEnd      time.Time
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedAt      time.Time
	UpdatedBy      string
	Revision       int32
}

//Trail Balance Entry
type TrialBalanceEntry struct {
	ID            uuid.UUID
	ReportID      uuid.UUID
	LedgerAccount string
	Debit         float64
	Credit        float64
	CreatedAt     time.Time
	CreatedBy     string
}

//Compliance Report
type ComplianceReport struct {
	ID             uuid.UUID
	OrganizationID string
	PeriodStart    time.Time
	PeriodEnd      time.Time
	Jurisdiction   string
	Details        string
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedAt      time.Time
	UpdatedBy      string
	Revision       int32
}