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
	param := db.AddCreditDebitNoteExternalRefParams{
		System:  ref.System,
		RefID:   ref.RefID,
	}
	e, err := r.queries.AddCreditDebitNoteExternalRef(ctx, param)
	if err != nil {
		return domain.ExternalRef{}, err
	}
	return mapSQLCExternalRefToDomain(e), nil
}

func (r *CreditDebitNoteRepo) ListExternalRefs(ctx context.Context, noteID uuid.UUID) ([]domain.ExternalRef, error) {
	rows, err := r.queries.ListCreditDebitNoteExternalRefs(ctx, noteID)
	if err != nil {
		return nil, err
	}
	refs := make([]domain.ExternalRef, len(rows))
	for i, e := range rows {
		refs[i] = mapSQLCExternalRefToDomain(e)
	}
	return refs, nil
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

func mapSQLCExternalRefToDomain(r db.CreditDebitNoteExternalRef) domain.ExternalRef {
	return domain.ExternalRef{
		ID:        r.ID,
		System:    r.System,
		RefID:     r.RefID,
		CreatedAt: r.CreatedAt.Time,
	}
}
