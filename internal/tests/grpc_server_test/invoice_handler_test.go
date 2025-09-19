package grpc_server_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports/grpc_server"
)


type InvoiceServiceInterface interface {
    CreateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error)
    GetInvoice(ctx context.Context, id uuid.UUID) (db.Invoice, error)
    UpdateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error)
    DeleteInvoice(ctx context.Context, id uuid.UUID) error
    ListInvoices(ctx context.Context, limit, offset int32) ([]db.Invoice, error)
    SearchInvoices(ctx context.Context, query string, limit, offset int32) ([]db.Invoice, error)

    CreateInvoiceItem(ctx context.Context, item db.InvoiceItem) (db.InvoiceItem, error)
    ListInvoiceItems(ctx context.Context, invoiceID uuid.UUID) ([]db.InvoiceItem, error)
    AddInvoiceTax(ctx context.Context, tax db.InvoiceTax) (db.InvoiceTax, error)
    AddInvoiceDiscount(ctx context.Context, discount db.InvoiceDiscount) (db.InvoiceDiscount, error)
}


type mockInvoiceService struct {
	CreateFn         func(ctx context.Context, inv db.Invoice) (db.Invoice, error)
	GetFn            func(ctx context.Context, id uuid.UUID) (db.Invoice, error)
	UpdateFn         func(ctx context.Context, inv db.Invoice) (db.Invoice, error)
	DeleteFn         func(ctx context.Context, id uuid.UUID) error
	ListFn           func(ctx context.Context, limit, offset int32) ([]db.Invoice, error)
	SearchFn         func(ctx context.Context, query string, limit, offset int32) ([]db.Invoice, error)
	CreateItemFn     func(ctx context.Context, item db.InvoiceItem) (db.InvoiceItem, error)
	ListItemsFn      func(ctx context.Context, invoiceID uuid.UUID) ([]db.InvoiceItem, error)
	AddTaxFn         func(ctx context.Context, tax db.InvoiceTax) (db.InvoiceTax, error)
	AddDiscountFn    func(ctx context.Context, discount db.InvoiceDiscount) (db.InvoiceDiscount, error)
}

// Safe mock methods
func (m *mockInvoiceService) CreateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error) {
	if m.CreateFn == nil {
		return db.Invoice{}, nil
	}
	return m.CreateFn(ctx, inv)
}

func (m *mockInvoiceService) GetInvoice(ctx context.Context, id uuid.UUID) (db.Invoice, error) {
	if m.GetFn == nil {
		return db.Invoice{}, nil
	}
	return m.GetFn(ctx, id)
}

func (m *mockInvoiceService) UpdateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error) {
	if m.UpdateFn == nil {
		return db.Invoice{}, nil
	}
	return m.UpdateFn(ctx, inv)
}

func (m *mockInvoiceService) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFn == nil {
		return nil
	}
	return m.DeleteFn(ctx, id)
}

func (m *mockInvoiceService) ListInvoices(ctx context.Context, limit, offset int32) ([]db.Invoice, error) {
	if m.ListFn == nil {
		return nil, nil
	}
	return m.ListFn(ctx, limit, offset)
}

func (m *mockInvoiceService) SearchInvoices(ctx context.Context, query string, limit, offset int32) ([]db.Invoice, error) {
	if m.SearchFn == nil {
		return nil, nil
	}
	return m.SearchFn(ctx, query, limit, offset)
}

func (m *mockInvoiceService) CreateInvoiceItem(ctx context.Context, item db.InvoiceItem) (db.InvoiceItem, error) {
	if m.CreateItemFn == nil {
		return db.InvoiceItem{}, nil
	}
	return m.CreateItemFn(ctx, item)
}

func (m *mockInvoiceService) ListInvoiceItems(ctx context.Context, invoiceID uuid.UUID) ([]db.InvoiceItem, error) {
	if m.ListItemsFn == nil {
		return nil, nil
	}
	return m.ListItemsFn(ctx, invoiceID)
}

func (m *mockInvoiceService) AddInvoiceTax(ctx context.Context, tax db.InvoiceTax) (db.InvoiceTax, error) {
	if m.AddTaxFn == nil {
		return db.InvoiceTax{}, nil
	}
	return m.AddTaxFn(ctx, tax)
}

func (m *mockInvoiceService) AddInvoiceDiscount(ctx context.Context, discount db.InvoiceDiscount) (db.InvoiceDiscount, error) {
	if m.AddDiscountFn == nil {
		return db.InvoiceDiscount{}, nil
	}
	return m.AddDiscountFn(ctx, discount)
}

//  Tests
func newRouter(h *grpc_server.InvoiceHandler) *chi.Mux {
	r := chi.NewRouter()
	// Invoice routes
	r.Post("/invoices", h.CreateInvoice)
	r.Get("/invoices/{id}", h.GetInvoice)
	r.Put("/invoices/{id}", h.UpdateInvoice)
	r.Delete("/invoices/{id}", h.DeleteInvoice)
	r.Get("/invoices", h.ListInvoices)

	// Search route
	r.Get("/invoices/search", h.SearchInvoices)

	// Invoice item routes
	r.Post("/invoices/{id}/items", h.CreateInvoiceItem)
	r.Get("/invoices/{id}/items", h.ListInvoiceItems)

	// Tax & Discount routes
	r.Post("/invoices/{id}/taxes", h.AddInvoiceTax)
	r.Post("/invoices/{id}/discounts", h.AddInvoiceDiscount)

	return r
}

func TestInvoiceHandler(t *testing.T) {
	now := time.Now()
	invID := uuid.New()
	inv := db.Invoice{
		ID:            invID,
		InvoiceNumber: "INV-001",
		Type:          "SALE",
		InvoiceDate:   now,
		DueDate:       sql.NullTime{Time: now.AddDate(0, 0, 30), Valid: true},
		GrandTotal:    "1180.0",
		Subtotal:      "1000.0",
		GstRate:       sql.NullString{String: "18.0", Valid: true},
		GstCgst:       sql.NullString{String: "90.0", Valid: true},
		GstSgst:       sql.NullString{String: "90.0", Valid: true},
		GstIgst:       sql.NullString{String: "0.0", Valid: true},
	}

	t.Run("CreateInvoice", func(t *testing.T) {
		mockSvc := &mockInvoiceService{
			CreateFn: func(ctx context.Context, i db.Invoice) (db.Invoice, error) { return inv, nil },
		}
		handler := grpc_server.NewInvoiceHandler(mockSvc)
		router := newRouter(handler)

		body, _ := json.Marshal(inv)
		req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var got db.Invoice
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		require.Equal(t, inv.ID, got.ID)
	})

	t.Run("SearchInvoices", func(t *testing.T) {
		mockSvc := &mockInvoiceService{
			SearchFn: func(ctx context.Context, query string, limit, offset int32) ([]db.Invoice, error) {
				return []db.Invoice{inv}, nil
			},
		}
		handler := grpc_server.NewInvoiceHandler(mockSvc)
		router := newRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/invoices/search?q=INV", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var got []db.Invoice
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		require.Len(t, got, 1)
		require.Equal(t, inv.ID, got[0].ID)
	})

	t.Run("CreateInvoiceItem", func(t *testing.T) {
		item := db.InvoiceItem{ID: uuid.New(), Name: "Item1", Quantity: 2}
		mockSvc := &mockInvoiceService{
			CreateItemFn: func(ctx context.Context, i db.InvoiceItem) (db.InvoiceItem, error) { return i, nil },
		}
		handler := grpc_server.NewInvoiceHandler(mockSvc)
		router := newRouter(handler)

		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPost, "/invoices/"+invID.String()+"/items", bytes.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var got db.InvoiceItem
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		require.Equal(t, item.ID, got.ID)
	})

	t.Run("ListInvoiceItems", func(t *testing.T) {
		mockSvc := &mockInvoiceService{
			ListItemsFn: func(ctx context.Context, invoiceID uuid.UUID) ([]db.InvoiceItem, error) {
				return []db.InvoiceItem{{ID: uuid.New(), Name: "Item1"}}, nil
			},
		}
		handler := grpc_server.NewInvoiceHandler(mockSvc)
		router := newRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/invoices/"+invID.String()+"/items", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var got []db.InvoiceItem
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		require.Len(t, got, 1)
	})

	t.Run("AddInvoiceTax", func(t *testing.T) {
		tax := db.InvoiceTax{ID: uuid.New(), Rate: "18.0"}
		mockSvc := &mockInvoiceService{
			AddTaxFn: func(ctx context.Context, t db.InvoiceTax) (db.InvoiceTax, error) { return t, nil },
		}
		handler := grpc_server.NewInvoiceHandler(mockSvc)
		router := newRouter(handler)

		body, _ := json.Marshal(tax)
		req := httptest.NewRequest(http.MethodPost, "/invoices/"+invID.String()+"/taxes", bytes.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var got db.InvoiceTax
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		require.Equal(t, tax.ID, got.ID)
	})

	t.Run("AddInvoiceDiscount", func(t *testing.T) {
		discount := db.InvoiceDiscount{ID: uuid.New(), Amount: "100"}
		mockSvc := &mockInvoiceService{
			AddDiscountFn: func(ctx context.Context, d db.InvoiceDiscount) (db.InvoiceDiscount, error) { return d, nil },
		}
		handler := grpc_server.NewInvoiceHandler(mockSvc)
		router := newRouter(handler)

		body, _ := json.Marshal(discount)
		req := httptest.NewRequest(http.MethodPost, "/invoices/"+invID.String()+"/discounts", bytes.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var got db.InvoiceDiscount
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		require.Equal(t, discount.ID, got.ID)
	})
}
