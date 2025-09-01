package domain

import (
    "time"
    "github.com/google/uuid"
)

type Invoice struct {
    ID                   uuid.UUID
    InvoiceNumber        string
    Type                 string
    InvoiceDate          time.Time
    DueDate              *time.Time
    DeliveryDate         *time.Time
    PartyRefID           *uuid.UUID
    OrganizationID       string
    PoNumber             *string
    EwayNumberLegacy     *string
    StatusNote           *string
    Status               string
    PaymentReference     *string
    ChallanNumber        *string
    ChallanDate          *time.Time
    LrNumber             *string
    TransporterName      *string
    TransporterID        *string
    VehicleNumber        *string
    AgainstInvoiceNumber *string
    AgainstInvoiceDate   *time.Time
    Subtotal             string
    GrandTotal           string
    GstRate              *string
    GstCgst              *string
    GstSgst              *string
    GstIgst              *string
    CreatedAt            time.Time
    CreatedBy            *string
    UpdatedAt            *time.Time
    UpdatedBy            *string
    Revision             int32
}

type InvoiceItem struct {
    ID           uuid.UUID
    InvoiceID    uuid.UUID
    Name         string
    Description  *string
    Hsn          *string
    Quantity     int32
    UnitPrice    string
    LineSubtotal string
    LineTotal    string
    CostCenterID *string
    CreatedAt    time.Time
    CreatedBy    *string
    UpdatedAt    *time.Time
    UpdatedBy    *string
    Revision     int32
}

type InvoiceTax struct {
	ID        uuid.UUID
	InvoiceID uuid.UUID
	Name      string
	Rate      string
	Amount    string
	CreatedAt time.Time
	CreatedBy *string
	Revision  *int32
}

type InvoiceDiscount struct {
	ID          uuid.UUID
	InvoiceID   uuid.UUID
	Description *string
	Amount      string
	CreatedAt   time.Time
	CreatedBy   *string
	Revision    *int32
}
