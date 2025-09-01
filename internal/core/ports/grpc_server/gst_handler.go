package grpc_server

import (
	"encoding/json"
	"net/http"
	"time"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
)

type GstHandler struct {
	svc *services.GstService
}

func NewGstHandler(svc *services.GstService) *GstHandler {
	return &GstHandler{svc: svc}
}

func (h *GstHandler) RegisterRoutes(r chi.Router) {
	r.Post("/gst/breakup", h.AddGstBreakup)
	r.Get("/gst/breakup/{invoice_id}", h.GetGstBreakup)

	r.Post("/gst/regime", h.AddGstRegime)
	r.Get("/gst/regime/{invoice_id}", h.GetGstRegime)

	r.Post("/gst/doc-status", h.AddGstDocStatus)
	r.Get("/gst/doc-status/{invoice_id}", h.GetGstDocStatus)
}

// ---------- Breakup ----------
func (h *GstHandler) AddGstBreakup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		InvoiceID     string   `json:"invoice_id"`
		TaxableAmount float64  `json:"taxable_amount"`
		CGST          *float64 `json:"cgst"`
		SGST          *float64 `json:"sgst"`
		IGST          *float64 `json:"igst"`
		TotalGST      *float64 `json:"total_gst"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.svc.AddGstBreakup(r.Context(), req.InvoiceID, req.TaxableAmount, req.CGST, req.SGST, req.IGST, req.TotalGST)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *GstHandler) GetGstBreakup(w http.ResponseWriter, r *http.Request) {
	invoiceID := chi.URLParam(r, "invoice_id")
	resp, err := h.svc.GetGstBreakup(r.Context(), invoiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// ---------- Regime ----------
func (h *GstHandler) AddGstRegime(w http.ResponseWriter, r *http.Request) {
	var req struct {
		InvoiceID     string  `json:"invoice_id"`
		GSTIN         string  `json:"gstin"`
		PlaceOfSupply string  `json:"place_of_supply"`
		ReverseCharge *bool   `json:"reverse_charge"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.svc.AddGstRegime(r.Context(), req.InvoiceID, req.GSTIN, req.PlaceOfSupply, req.ReverseCharge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *GstHandler) GetGstRegime(w http.ResponseWriter, r *http.Request) {
	invoiceID := chi.URLParam(r, "invoice_id")
	resp, err := h.svc.GetGstRegime(r.Context(), invoiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// ---------- Doc Status ----------
func (h *GstHandler) AddGstDocStatus(w http.ResponseWriter, r *http.Request) {
    var req struct {
        InvoiceID      string     `json:"invoice_id"`
        EinvoiceStatus *string    `json:"einvoice_status"`
        IRN            *string    `json:"irn"`
        AckNo          *string    `json:"ack_no"`
        AckDate        *time.Time `json:"ack_date"`
        EwayStatus     *string    `json:"eway_status"`
        EwayBillNo     *string    `json:"eway_bill_no"`
        EwayValidUpto  *time.Time `json:"eway_valid_upto"`
        LastError      *string    `json:"last_error"`
        LastSyncedAt   *time.Time `json:"last_synced_at"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    resp, err := h.svc.AddGstDocStatus(
        r.Context(),
        req.InvoiceID,
        req.EinvoiceStatus,
        req.IRN,
        req.AckNo,
        req.AckDate,
        req.EwayStatus,
        req.EwayBillNo,
        req.EwayValidUpto,
        req.LastError,
        req.LastSyncedAt,
    )
    if err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, resp)
}


func (h *GstHandler) GetGstDocStatus(w http.ResponseWriter, r *http.Request) {
	invoiceID := chi.URLParam(r, "invoice_id")
	resp, err := h.svc.GetGstDocStatus(r.Context(), invoiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func respondJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

type JSONTime time.Time

func (jt *JSONTime) UnmarshalJSON(b []byte) error {
    str := strings.Trim(string(b), `"`)
    if str == "null" || str == "" {
        return nil
    }
    t, err := time.Parse("2006-01-02", str) // supports YYYY-MM-DD
    if err != nil {
        return err
    }
    *jt = JSONTime(t)
    return nil
}

