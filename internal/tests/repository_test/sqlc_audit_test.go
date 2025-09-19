package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
)

func TestAuditRepository(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	repo := repository.NewAuditRepository(dbConn, nil)
	ctx := context.Background()

	eventID := uuid.New()
	userID := uuid.New()
	timestamp := time.Now()
	createdAt := time.Now()
	userIDStr := userID.String()
	userIDPtr := &userIDStr

	// ---------- RecordAuditEvent ----------
	mock.ExpectQuery(`INSERT INTO audit_events`).
		WithArgs(userID, "CREATE", sqlmock.AnyArg(), "Details", "Order", "ORD123").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "action", "timestamp", "details", "resource_type", "resource_id", "created_at",
		}).AddRow(
			eventID, userID, "CREATE", timestamp, "Details", "Order", "ORD123", createdAt,
		))

	event := &db.AuditEvent{
		UserID:       userID.String(),
		Action:       "CREATE",
		Timestamp:    timestamp,
		Details:      sql.NullString{String: "Details", Valid: true},
		ResourceType: sql.NullString{String: "Order", Valid: true},
		ResourceID:   sql.NullString{String: "ORD123", Valid: true},
		CreatedAt:    sql.NullTime{Time: createdAt, Valid: true},
	}

	recorded, err := repo.RecordAuditEvent(ctx, event)
	require.NoError(t, err)
	require.Equal(t, eventID, recorded.ID)
	require.Equal(t, "CREATE", recorded.Action)

	// ---------- GetAuditEventByID ----------
	mock.ExpectQuery(`SELECT .* FROM audit_events WHERE id = \$1`).
		WithArgs(eventID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "action", "timestamp", "details", "resource_type", "resource_id", "created_at",
		}).AddRow(
			eventID, userID, "CREATE", timestamp, "Details", "Order", "ORD123", createdAt,
		))

	got, err := repo.GetAuditEventByID(ctx, eventID)
	require.NoError(t, err)
	require.Equal(t, "CREATE", got.Action)

	// ---------- ListAuditEvents ----------
	mock.ExpectQuery(`SELECT .* FROM audit_events ORDER BY timestamp DESC LIMIT \$1 OFFSET \$2`).
		WithArgs(int32(10), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "action", "timestamp", "details", "resource_type", "resource_id", "created_at",
		}).AddRow(
			eventID, userID.String(), "CREATE", timestamp, "Details", "Order", "ORD123", createdAt,
		))


	list, err := repo.ListAuditEvents(ctx, db.Pagination{Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "CREATE", list[0].Action)

	// ---------- FilterAuditEvents ----------
	filter := db.FilterParams{
		UserID:       userIDPtr, // Use pointer only if you want actual filter
		Action:       nil,
		ResourceType: nil,
		ResourceID:   nil,
		FromDate:     nil,
		ToDate:       nil,
	}

	// sqlmock must receive exactly the arguments SQLC will pass
	mock.ExpectQuery(`SELECT .* FROM audit_events.*LIMIT \$7 OFFSET \$8`).
		WithArgs(
			userID.String(),  // Column1
			"",               // Column2
			"",               // Column3
			"",               // Column4
			time.Time{},      // Column5
			time.Time{},      // Column6
			int32(10),        // Limit
			int32(0),         // Offset
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "action", "timestamp", "details", "resource_type", "resource_id", "created_at",
		}).AddRow(
			eventID, userID.String(), "CREATE", timestamp, "Details", "Order", "ORD123", createdAt,
		))



	filtered, err := repo.FilterAuditEvents(ctx, filter, db.Pagination{Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.Len(t, filtered, 1)
	require.Equal(t, "CREATE", filtered[0].Action)

	// ---------- Ensure all expectations met ----------
	require.NoError(t, mock.ExpectationsWereMet())
}