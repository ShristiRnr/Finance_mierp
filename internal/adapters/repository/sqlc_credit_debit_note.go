package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type CreditDebitNoteRepo struct {
	queries *db.Queries
}

func NewCreditDebitNoteRepo(q *db.Queries) ports.CreditDebitNoteRepository {
	return &CreditDebitNoteRepo{queries: q}
}

func (r *CreditDebitNoteRepo) Create(ctx context.Context, note domain.CreditDebitNote) (domain.CreditDebitNote, error) {
	params := db.CreateCreditDebitNoteParams{
		InvoiceID: note.InvoiceID,
		Type:      string(note.Type),
		Amount:    note.Amount,
		Reason:    sql.NullString{String: note.Reason, Valid: note.Reason != ""},
		CreatedBy: sql.NullString{String: note.CreatedBy, Valid: note.CreatedBy != ""},
		UpdatedBy: sql.NullString{String: note.UpdatedBy, Valid: note.UpdatedBy != ""},
	}
	n, err := r.queries.CreateCreditDebitNote(ctx, params)
	if err != nil {
		return domain.CreditDebitNote{}, err
	}
	return mapSQLCToDomains(n), nil
}

func (r *CreditDebitNoteRepo) Get(ctx context.Context, id uuid.UUID) (domain.CreditDebitNote, error) {
	n, err := r.queries.GetCreditDebitNote(ctx, id)
	if err != nil {
		return domain.CreditDebitNote{}, err
	}
	return mapSQLCToDomains(n), nil
}

func (r *CreditDebitNoteRepo) List(ctx context.Context, limit, offset int32) ([]domain.CreditDebitNote, error) {
	rows, err := r.queries.ListCreditDebitNotes(ctx, db.ListCreditDebitNotesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	items := make([]domain.CreditDebitNote, len(rows))
	for i, n := range rows {
		items[i] = mapSQLCToDomains(n)
	}
	return items, nil
}

func (r *CreditDebitNoteRepo) Update(ctx context.Context, note domain.CreditDebitNote) (domain.CreditDebitNote, error) {
	params := db.UpdateCreditDebitNoteParams{
		ID:        note.ID,
		InvoiceID: note.InvoiceID,
		Type:      string(note.Type),
		Amount:    note.Amount,
		Reason:    sql.NullString{String: note.Reason, Valid: note.Reason != ""},
		UpdatedBy: sql.NullString{String: note.UpdatedBy, Valid: note.UpdatedBy != ""},
	}
	n, err := r.queries.UpdateCreditDebitNote(ctx, params)
	if err != nil {
		return domain.CreditDebitNote{}, err
	}
	return mapSQLCToDomains(n), nil
}

func (r *CreditDebitNoteRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteCreditDebitNote(ctx, id)
}

func (r *CreditDebitNoteRepo) AddExternalRef(ctx context.Context, ref domain.ExternalRef) (domain.ExternalRef, error) {
	params := db.AddCreditDebitNoteExternalRefParams{
		NoteID: ref.ID,
		System: ref.System,
		RefID:  ref.RefID,
	}
	rf, err := r.queries.AddCreditDebitNoteExternalRef(ctx, params)
	if err != nil {
		return domain.ExternalRef{}, err
	}
	return domain.ExternalRef{
		ID:        rf.ID,
		System:    rf.System,
		RefID:     rf.RefID,
		CreatedAt: rf.CreatedAt.Time,
	}, nil
}

func (r *CreditDebitNoteRepo) ListExternalRefs(ctx context.Context, noteID uuid.UUID) ([]domain.ExternalRef, error) {
	refs, err := r.queries.ListCreditDebitNoteExternalRefs(ctx, noteID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.ExternalRef, len(refs))
	for i, rf := range refs {
		out[i] = domain.ExternalRef{
			ID:        rf.ID,
			System:    rf.System,
			RefID:     rf.RefID,
			CreatedAt: rf.CreatedAt.Time,
		}
	}
	return out, nil
}

func mapSQLCToDomains(n db.CreditDebitNote) domain.CreditDebitNote {
	var reason, createdBy, updatedBy string
	if n.Reason.Valid {
		reason = n.Reason.String
	}
	if n.CreatedBy.Valid {
		createdBy = n.CreatedBy.String
	}
	if n.UpdatedBy.Valid {
		updatedBy = n.UpdatedBy.String
	}

	return domain.CreditDebitNote{
		ID:        n.ID,
		InvoiceID: n.InvoiceID,
		Type:      domain.NoteType(n.Type),
		Amount:    n.Amount,
		Reason:    reason,
		CreatedAt: n.CreatedAt.Time,
		CreatedBy: createdBy,
		UpdatedAt: n.UpdatedAt.Time,
		UpdatedBy: updatedBy,
	}
}