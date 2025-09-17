package services_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
)

// ------------------- Helper -------------------
func makeAuditEvent() *db.AuditEvent {
	return &db.AuditEvent{
		ID:           uuid.New(),
		UserID:       "user1",
		Action:       "CREATE",
		Timestamp:    time.Now(),
		Details:      sql.NullString{String: "details", Valid: true},
		ResourceType: sql.NullString{String: "Order", Valid: true},
		ResourceID:   sql.NullString{String: "ORD123", Valid: true},
	}
}

// ------------------- Test Record/Insert -------------------
func TestAuditService_Record(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbConn.Close()

	repo := repository.NewAuditRepository(dbConn, nil)
	service := services.NewAuditService(repo, nil)

	event := makeAuditEvent()

	// Mock SQL insert
	mock.ExpectQuery(`INSERT INTO audit_events .* RETURNING .*`).
		WithArgs(event.UserID, event.Action, event.Timestamp, event.Details, event.ResourceType, event.ResourceID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "action", "timestamp", "details", "resource_type", "resource_id", "created_at",
		}).AddRow(
			event.ID, event.UserID, event.Action, event.Timestamp,
			event.Details, event.ResourceType, event.ResourceID, time.Now(),
		))

	got, err := service.Record(context.Background(), event)
	assert.NoError(t, err)
	assert.Equal(t, event.UserID, got.UserID)
	assert.Equal(t, event.Action, got.Action)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ------------------- Test GetByID -------------------
func TestAuditService_GetByID(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbConn.Close()

	repo := repository.NewAuditRepository(dbConn, nil)
	service := services.NewAuditService(repo, nil)

	event := makeAuditEvent()

	mock.ExpectQuery(`SELECT .* FROM audit_events WHERE id = \$1`).
		WithArgs(event.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "action", "timestamp", "details", "resource_type", "resource_id", "created_at",
		}).AddRow(
			event.ID, event.UserID, event.Action, event.Timestamp,
			event.Details, event.ResourceType, event.ResourceID, time.Now(),
		))

	got, err := service.GetByID(context.Background(), event.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, event.ID, got.ID)
	assert.Equal(t, event.UserID, got.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ------------------- Test List -------------------
func TestAuditService_List(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbConn.Close()

	repo := repository.NewAuditRepository(dbConn, nil)
	service := services.NewAuditService(repo, nil)

	event1 := makeAuditEvent()
	event2 := makeAuditEvent()

	mock.ExpectQuery(`SELECT .* FROM audit_events ORDER BY timestamp DESC LIMIT .* OFFSET .*`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "action", "timestamp", "details", "resource_type", "resource_id", "created_at",
		}).AddRow(event1.ID, event1.UserID, event1.Action, event1.Timestamp, event1.Details, event1.ResourceType, event1.ResourceID, time.Now()).
			AddRow(event2.ID, event2.UserID, event2.Action, event2.Timestamp, event2.Details, event2.ResourceType, event2.ResourceID, time.Now()))

	got, err := service.List(context.Background(), db.Pagination{Limit: 10, Offset: 0})
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ------------------- Test Filter -------------------
func TestAuditService_Filter(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbConn.Close()

	repo := repository.NewAuditRepository(dbConn, nil)
	service := services.NewAuditService(repo, nil)

	event := makeAuditEvent()

	filter := db.FilterParams{
		UserID:       &event.UserID,
		Action:       nil,
		ResourceType: nil,
		ResourceID:   nil,
		FromDate:     nil,
		ToDate:       nil,
	}

	mock.ExpectQuery(`SELECT .* FROM audit_events WHERE .*`).
		WithArgs(event.UserID, "", "", "", sqlmock.AnyArg(), sqlmock.AnyArg(), int32(5), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "action", "timestamp", "details", "resource_type", "resource_id", "created_at",
		}).AddRow(event.ID, event.UserID, event.Action, event.Timestamp, event.Details, event.ResourceType, event.ResourceID, time.Now()))

	got, err := service.Filter(context.Background(), filter, db.Pagination{Limit: 5, Offset: 0})
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, event.UserID, got[0].UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
