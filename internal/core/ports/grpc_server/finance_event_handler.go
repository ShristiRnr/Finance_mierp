package grpc_server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	"github.com/go-chi/chi/v5"
)

type FinanceEventHandler struct {
	svc *services.FinanceEventService
}

func NewFinanceEventHandler(svc *services.FinanceEventService) *FinanceEventHandler {
	return &FinanceEventHandler{svc: svc}
}

// ============================
// Invoice Created
// ============================

// POST /org/{orgID}/invoice-events
func (h *FinanceEventHandler) InsertInvoiceCreated(w http.ResponseWriter, r *http.Request) {
	var event domain.FinanceInvoiceCreatedEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	event.OrganizationID = chi.URLParam(r, "orgID")

	created, err := h.svc.RecordInvoiceCreated(r.Context(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

// GET /org/{orgID}/invoice-events?limit=10&offset=0
func (h *FinanceEventHandler) ListInvoiceCreated(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	events, err := h.svc.GetInvoiceCreatedEvents(r.Context(), orgID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

// ============================
// Payment Received
// ============================

// POST /org/{orgID}/payment-received-events
func (h *FinanceEventHandler) InsertPaymentReceived(w http.ResponseWriter, r *http.Request) {
	var event domain.FinancePaymentReceivedEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	event.OrganizationID = chi.URLParam(r, "orgID")

	created, err := h.svc.RecordPaymentReceived(r.Context(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

// GET /org/{orgID}/payment-received-events
func (h *FinanceEventHandler) ListPaymentReceived(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	events, err := h.svc.GetPaymentReceivedEvents(r.Context(), orgID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

// ============================
// Inventory Cost Posted
// ============================

// POST /org/{orgID}/inventory-cost-posted-events
func (h *FinanceEventHandler) InsertInventoryCostPosted(w http.ResponseWriter, r *http.Request) {
	var event domain.InventoryCostPostedEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	event.OrganizationID = chi.URLParam(r, "orgID")

	created, err := h.svc.RecordInventoryCostPosted(r.Context(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

// GET /org/{orgID}/inventory-cost-posted-events
func (h *FinanceEventHandler) ListInventoryCostPosted(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	events, err := h.svc.GetInventoryCostPostedEvents(r.Context(), orgID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

// ============================
// Payroll Posted
// ============================

// POST /org/{orgID}/payroll-posted-events
func (h *FinanceEventHandler) InsertPayrollPosted(w http.ResponseWriter, r *http.Request) {
	var event domain.PayrollPostedEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	event.OrganizationID = chi.URLParam(r, "orgID")

	created, err := h.svc.RecordPayrollPosted(r.Context(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

// GET /org/{orgID}/payroll-posted-events
func (h *FinanceEventHandler) ListPayrollPosted(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	events, err := h.svc.GetPayrollPostedEvents(r.Context(), orgID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

// ============================
// Vendor Bill Approved
// ============================

// POST /org/{orgID}/vendor-bill-approved-events
func (h *FinanceEventHandler) InsertVendorBillApproved(w http.ResponseWriter, r *http.Request) {
	var event domain.VendorBillApprovedEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	event.OrganizationID = chi.URLParam(r, "orgID")

	created, err := h.svc.RecordVendorBillApproved(r.Context(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

// GET /org/{orgID}/vendor-bill-approved-events
func (h *FinanceEventHandler) ListVendorBillApproved(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	events, err := h.svc.GetVendorBillApprovedEvents(r.Context(), orgID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}
