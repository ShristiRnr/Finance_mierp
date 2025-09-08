package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type CreditDebitNoteRepo struct {
	queries *db.Queries
}

func NewCreditDebitNoteRepo(q *db.Queries) ports.CreditDebitNoteRepository {
	return &CreditDebitNoteRepo{queries: q}
}

func (r *CreditDebitNoteRepo) Create(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error) {
	params := db.CreateCreditDebitNoteParams{
		InvoiceID: note.InvoiceID,
		Type:      note.Type,
		Amount:    note.Amount,
		Reason:    note.Reason,
		CreatedBy: note.CreatedBy,
		UpdatedBy: note.UpdatedBy,
	}
	n, err := r.queries.CreateCreditDebitNote(ctx, params)
	if err != nil {
		return db.CreditDebitNote{}, err
	}
	return mapSQLCToDomains(n), nil
}

func (r *CreditDebitNoteRepo) Get(ctx context.Context, id uuid.UUID) (db.CreditDebitNote, error) {
	n, err := r.queries.GetCreditDebitNote(ctx, id)
	if err != nil {
		return db.CreditDebitNote{}, err
	}
	return mapSQLCToDomains(n), nil
}

func (r *CreditDebitNoteRepo) List(ctx context.Context, limit, offset int32) ([]db.CreditDebitNote, error) {
	rows, err := r.queries.ListCreditDebitNotes(ctx, db.ListCreditDebitNotesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	items := make([]db.CreditDebitNote, len(rows))
	for i, n := range rows {
		items[i] = mapSQLCToDomains(n)
	}
	return items, nil
}

func (r *CreditDebitNoteRepo) Update(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error) {
	params := db.UpdateCreditDebitNoteParams{
		ID:        note.ID,
		InvoiceID: note.InvoiceID,
		Type:      note.Type,
		Amount:    note.Amount, 
		Reason:    note.Reason,
		UpdatedBy: note.UpdatedBy,
	}

	n, err := r.queries.UpdateCreditDebitNote(ctx, params)
	if err != nil {
		return db.CreditDebitNote{}, err
	}
	return mapSQLCToDomains(n), nil
}


func (r *CreditDebitNoteRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteCreditDebitNote(ctx, id)
}

func mapSQLCToDomains(n db.CreditDebitNote) db.CreditDebitNote {
	return db.CreditDebitNote{
		ID:        n.ID,
		InvoiceID: n.InvoiceID,
		Type:      n.Type,
		Amount:    n.Amount,
		Reason:    n.Reason,
		CreatedAt: n.CreatedAt,
		CreatedBy: n.CreatedBy,
		UpdatedAt: n.UpdatedAt,
		UpdatedBy: n.UpdatedBy,
	}
}
