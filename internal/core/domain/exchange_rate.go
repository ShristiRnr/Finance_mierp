// domain/exchange_rate.go
package domain

import (
    "time"

    "github.com/google/uuid"
)

//Exchnage Rate
type ExchangeRate struct {
    ID            uuid.UUID
    BaseCurrency  string
    QuoteCurrency string
    Rate          string // keep as string if you want decimal precision
    AsOf          time.Time
    CreatedAt     time.Time
    CreatedBy     *string
    UpdatedAt     time.Time
    UpdatedBy     *string
    Revision      *string
}