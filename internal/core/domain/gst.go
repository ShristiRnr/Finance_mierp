package domain

import (
	"time"

	"github.com/google/uuid"
)

//Gst Breakup
type GstBreakup struct {
	ID            uuid.UUID
	InvoiceID     uuid.UUID
	TaxableAmount float64
	CGST          *float64
	SGST          *float64
	IGST          *float64
	TotalGST      *float64
	CreatedAt     time.Time
	CreatedBy     string
	Revision      int32
}

//Gst Regime
type GstRegime struct {
	ID            uuid.UUID
	InvoiceID     uuid.UUID
	GSTIN         string
	PlaceOfSupply string
	ReverseCharge *bool
	CreatedAt     time.Time
	CreatedBy     string
	Revision      int32
}

//Gst Document Status
type GstDocStatus struct {
	ID            uuid.UUID
	InvoiceID     uuid.UUID
	EinvoiceStatus *string
	IRN            *string
	AckNo          *string
	AckDate        *time.Time
	EwayStatus     *string
	EwayBillNo     *string
	EwayValidUpto  *time.Time
	LastError      *string
	LastSyncedAt   *time.Time
	CreatedAt      time.Time
	CreatedBy      string
	Revision       int32
}
