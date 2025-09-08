package grpc_server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
)

// InvoiceHandler handles HTTP requests for invoices and related entities
type InvoiceHandler struct {
	svc *services.InvoiceService
}

// NewInvoiceHandler creates a new handler
func NewInvoiceHandler(svc *services.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{svc: svc}
}

// RegisterRoutes registers all invoice routes
// func (h *InvoiceHandler) RegisterRoutes(r chi.Router) {
// 	// Invoice CRUD
// 	r.Post("/invoices", h.CreateInvoice)
// 	r.Get("/invoices/{id}", h.GetInvoice)
// 	r.Put("/invoices/{id}", h.UpdateInvoice)
// 	r.Delete("/invoices/{id}", h.DeleteInvoice)
// 	r.Get("/invoices", h.ListInvoices)
// 	r.Get("/invoices/search", h.SearchInvoices)

// 	// Invoice items
// 	r.Post("/invoices/{id}/items", h.CreateInvoiceItem)
// 	r.Get("/invoices/{id}/items", h.ListInvoiceItems)

// 	// Invoice taxes and discounts
// 	r.Post("/invoices/{id}/taxes", h.AddInvoiceTax)
// 	r.Post("/invoices/{id}/discounts", h.AddInvoiceDiscount)
// }

// ---------- Invoice Handlers ----------

func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var inv db.Invoice
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := h.svc.CreateInvoice(r.Context(), inv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

func (h *InvoiceHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	inv, err := h.svc.GetInvoice(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(inv)
}

func (h *InvoiceHandler) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var inv db.Invoice
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	inv.ID = id

	updated, err := h.svc.UpdateInvoice(r.Context(), inv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (h *InvoiceHandler) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteInvoice(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *InvoiceHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagingParams(r)
	invoices, err := h.svc.ListInvoices(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(invoices)
}

func (h *InvoiceHandler) SearchInvoices(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limit, offset := parsePagingParams(r)
	invoices, err := h.svc.SearchInvoices(r.Context(), query, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(invoices)
}

// ---------- InvoiceItem Handlers ----------

func (h *InvoiceHandler) CreateInvoiceItem(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := parseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item db.InvoiceItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item.InvoiceID = invoiceID

	created, err := h.svc.CreateInvoiceItem(r.Context(), item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

func (h *InvoiceHandler) ListInvoiceItems(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := parseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items, err := h.svc.ListInvoiceItems(r.Context(), invoiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(items)
}

// ---------- InvoiceTax & Discount Handlers ----------

func (h *InvoiceHandler) AddInvoiceTax(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := parseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tax db.InvoiceTax
	if err := json.NewDecoder(r.Body).Decode(&tax); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tax.InvoiceID = invoiceID

	created, err := h.svc.AddInvoiceTax(r.Context(), tax)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

func (h *InvoiceHandler) AddInvoiceDiscount(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := parseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var disc db.InvoiceDiscount
	if err := json.NewDecoder(r.Body).Decode(&disc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	disc.InvoiceID = invoiceID

	created, err := h.svc.AddInvoiceDiscount(r.Context(), disc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

// ---------- Helpers ----------

func parseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, param)
	return uuid.Parse(idStr)
}

func parsePagingParams(r *http.Request) (limit, offset int32) {
	limit = 100
	offset = 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := parseInt32(l); err == nil {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := parseInt32(o); err == nil {
			offset = parsed
		}
	}
	return
}

func parseInt32(s string) (int32, error) {
	var i int32
	_, err := fmt.Sscan(s, &i)
	return i, err
}
