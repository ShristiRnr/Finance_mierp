package repository_test

import (
	"context"
	"testing"
	"time"
	"database/sql"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
)

// ---------------- MockQueries ----------------
type MockQueries struct {
	mock.Mock
}

func (m *MockQueries) AddGstBreakup(ctx context.Context, arg db.AddGstBreakupParams) (db.GstBreakup, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.GstBreakup), args.Error(1)
}
func (m *MockQueries) GetGstBreakup(ctx context.Context, invoiceID uuid.UUID) (db.GstBreakup, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).(db.GstBreakup), args.Error(1)
}

func (m *MockQueries) AddGstRegime(ctx context.Context, arg db.AddGstRegimeParams) (db.GstRegime, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.GstRegime), args.Error(1)
}
func (m *MockQueries) GetGstRegime(ctx context.Context, invoiceID uuid.UUID) (db.GstRegime, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).(db.GstRegime), args.Error(1)
}

func (m *MockQueries) AddGstDocStatus(ctx context.Context, arg db.AddGstDocStatusParams) (db.GstDocStatus, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.GstDocStatus), args.Error(1)
}
func (m *MockQueries) GetGstDocStatus(ctx context.Context, invoiceID uuid.UUID) (db.GstDocStatus, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).(db.GstDocStatus), args.Error(1)
}

// ---------------- Test Helpers ----------------
func floatPtr(f float64) *float64   { return &f }
func strPtr(s string) *string       { return &s }
func boolPtr(b bool) *bool          { return &b }
func timePtr(t time.Time) *time.Time { return &t }

func makeBreakup() db.GstBreakup {
	return db.GstBreakup{
		ID:            uuid.New(),
		InvoiceID:     uuid.New(),
		TaxableAmount: "1000",
		Cgst:          sqlNullString("50"),
		Sgst:          sqlNullString("50"),
		Igst:          sqlNullString("0"),
		TotalGst:      sqlNullString("100"),
	}
}

func makeRegime() db.GstRegime {
	return db.GstRegime{
		ID:            uuid.New(),
		InvoiceID:     uuid.New(),
		Gstin:         "29ABCDE1234F2Z5",
		PlaceOfSupply: "KA",
		ReverseCharge: sqlNullBool(false),
	}
}

func makeDocStatus() db.GstDocStatus {
	now := time.Now()
	return db.GstDocStatus{
		ID:             uuid.New(),
		InvoiceID:      uuid.New(),
		EinvoiceStatus: sqlNullString("SUCCESS"),
		Irn:            sqlNullString("IRN123"),
		AckNo:          sqlNullString("ACK456"),
		AckDate:        sqlNullTime(now),
		EwayStatus:     sqlNullString("VALID"),
		EwayBillNo:     sqlNullString("EB123"),
		EwayValidUpto:  sqlNullTime(now.Add(24 * time.Hour)),
		LastError:      sqlNullString(""),
		LastSyncedAt:   sqlNullTime(now),
	}
}

func sqlNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func sqlNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: true}
}

func sqlNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

// ---------------- Tests ----------------
func TestGstRepo_AddAndGetBreakup(t *testing.T) {
	ctx := context.Background()
	mockQ := new(MockQueries)
	repo := repository.NewGstRepo(mockQ)

	exp := makeBreakup()
	params := db.AddGstBreakupParams{
		InvoiceID:     exp.InvoiceID,
		TaxableAmount: exp.TaxableAmount,
		Cgst:          exp.Cgst,
		Sgst:          exp.Sgst,
		Igst:          exp.Igst,
		TotalGst:      exp.TotalGst,
	}

	mockQ.On("AddGstBreakup", ctx, params).Return(exp, nil)
	got, err := repo.AddGstBreakup(ctx, exp.InvoiceID, 1000, floatPtr(50), floatPtr(50), floatPtr(0), floatPtr(100))
	assert.NoError(t, err)
	assert.Equal(t, exp.InvoiceID, got.InvoiceID)

	mockQ.On("GetGstBreakup", ctx, exp.InvoiceID).Return(exp, nil)
	got2, err := repo.GetGstBreakup(ctx, exp.InvoiceID)
	assert.NoError(t, err)
	assert.Equal(t, exp.ID, got2.ID)
}

func TestGstRepo_AddAndGetRegime(t *testing.T) {
	ctx := context.Background()
	mockQ := new(MockQueries)
	repo := repository.NewGstRepo(mockQ)

	exp := makeRegime()
	params := db.AddGstRegimeParams{
		InvoiceID:     exp.InvoiceID,
		Gstin:         exp.Gstin,
		PlaceOfSupply: exp.PlaceOfSupply,
		ReverseCharge: exp.ReverseCharge,
	}

	mockQ.On("AddGstRegime", ctx, params).Return(exp, nil)
	got, err := repo.AddGstRegime(ctx, exp.InvoiceID, exp.Gstin, exp.PlaceOfSupply, boolPtr(false))
	assert.NoError(t, err)
	assert.Equal(t, exp.Gstin, got.Gstin)

	mockQ.On("GetGstRegime", ctx, exp.InvoiceID).Return(exp, nil)
	got2, err := repo.GetGstRegime(ctx, exp.InvoiceID)
	assert.NoError(t, err)
	assert.Equal(t, exp.ID, got2.ID)
}

func TestGstRepo_AddAndGetDocStatus(t *testing.T) {
	ctx := context.Background()
	mockQ := new(MockQueries)
	repo := repository.NewGstRepo(mockQ)

	exp := makeDocStatus()
	params := db.AddGstDocStatusParams{
		InvoiceID:      exp.InvoiceID,
		EinvoiceStatus: exp.EinvoiceStatus,
		Irn:            exp.Irn,
		AckNo:          exp.AckNo,
		AckDate:        exp.AckDate,
		EwayStatus:     exp.EwayStatus,
		EwayBillNo:     exp.EwayBillNo,
		EwayValidUpto:  exp.EwayValidUpto,
		LastError:      exp.LastError,
		LastSyncedAt:   exp.LastSyncedAt,
	}

	mockQ.On("AddGstDocStatus", ctx, params).Return(exp, nil)
	got, err := repo.AddGstDocStatus(ctx, exp.InvoiceID,
		strPtr("SUCCESS"), strPtr("IRN123"), strPtr("ACK456"),
		timePtr(exp.AckDate.Time),
		strPtr("VALID"), strPtr("EB123"),
		timePtr(exp.EwayValidUpto.Time),
		nil, timePtr(exp.LastSyncedAt.Time),
	)
	assert.NoError(t, err)
	assert.Equal(t, exp.InvoiceID, got.InvoiceID)

	mockQ.On("GetGstDocStatus", ctx, exp.InvoiceID).Return(exp, nil)
	got2, err := repo.GetGstDocStatus(ctx, exp.InvoiceID)
	assert.NoError(t, err)
	assert.Equal(t, exp.ID, got2.ID)
}
