package grpc_server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
)

type FinancialReportsHandler struct {
	svc *services.FinancialReportsService
}

func NewFinancialReportsHandler(svc *services.FinancialReportsService) *FinancialReportsHandler {
	return &FinancialReportsHandler{svc: svc}
}

func (h *FinancialReportsHandler) RegisterRoutes(r chi.Router) {
	// Profit & Loss
	r.Post("/reports/profit-loss", h.GenerateProfitLoss)
	r.Get("/reports/profit-loss/{id}", h.GetProfitLoss)
	r.Get("/reports/profit-loss", h.ListProfitLossReports)

	// Balance Sheet
	r.Post("/reports/balance-sheet", h.GenerateBalanceSheet)
	r.Get("/reports/balance-sheet/{id}", h.GetBalanceSheet)
	r.Get("/reports/balance-sheet", h.ListBalanceSheetReports)

	// Trial Balance
	r.Post("/reports/trial-balance", h.CreateTrialBalance)
	r.Post("/reports/trial-balance/{id}/entries", h.AddTrialBalanceEntry)
	r.Get("/reports/trial-balance/{id}", h.GetTrialBalance)
	r.Get("/reports/trial-balance/{id}/entries", h.ListTrialBalanceEntries)
	r.Get("/reports/trial-balance", h.ListTrialBalanceReports)

	// Compliance
	r.Post("/reports/compliance", h.GenerateCompliance)
	r.Get("/reports/compliance/{id}", h.GetCompliance)
	r.Get("/reports/compliance", h.ListComplianceReports)
}

//
// ==========================
// Profit & Loss
// ==========================
func (h *FinancialReportsHandler) GenerateProfitLoss(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrganizationID string  `json:"organization_id"`
		PeriodStart    string  `json:"period_start"`
		PeriodEnd      string  `json:"period_end"`
		Revenue        float64 `json:"total_revenue"`
		Expenses       float64 `json:"total_expenses"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	report, err := h.svc.GenerateProfitLoss(r.Context(), req.OrganizationID, req.PeriodStart, req.PeriodEnd, req.Revenue, req.Expenses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(report)
}

func (h *FinancialReportsHandler) GetProfitLoss(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	report, err := h.svc.GetProfitLoss(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(report)
}

func (h *FinancialReportsHandler) ListProfitLossReports(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("org_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	reports, err := h.svc.ListProfitLossReports(r.Context(), orgID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(reports)
}

//
// ==========================
// Balance Sheet
// ==========================
func (h *FinancialReportsHandler) GenerateBalanceSheet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrganizationID   string  `json:"organization_id"`
		PeriodStart      string  `json:"period_start"`
		PeriodEnd        string  `json:"period_end"`
		TotalAssets      float64 `json:"total_assets"`
		TotalLiabilities float64 `json:"total_liabilities"`
		NetWorth         float64 `json:"net_worth"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	report, err := h.svc.GenerateBalanceSheet(r.Context(), req.OrganizationID, req.PeriodStart, req.PeriodEnd, req.TotalAssets, req.TotalLiabilities)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(report)
}

func (h *FinancialReportsHandler) GetBalanceSheet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	report, err := h.svc.GetBalanceSheet(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(report)
}

func (h *FinancialReportsHandler) ListBalanceSheetReports(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("org_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	reports, err := h.svc.ListBalanceSheetReports(r.Context(), orgID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(reports)
}

//
// ==========================
// Trial Balance
// ==========================
func (h *FinancialReportsHandler) CreateTrialBalance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrganizationID string `json:"organization_id"`
		PeriodStart    string `json:"period_start"`
		PeriodEnd      string `json:"period_end"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	report, err := h.svc.CreateTrialBalance(r.Context(), req.OrganizationID, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(report)
}

func (h *FinancialReportsHandler) AddTrialBalanceEntry(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ReportID      string  `json:"report_id"`
		LedgerAccount string  `json:"ledger_account"`
		Debit         float64 `json:"debit"`
		Credit        float64 `json:"credit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	entry, err := h.svc.AddTrialBalanceEntry(r.Context(), req.ReportID, req.LedgerAccount, req.Debit, req.Credit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(entry)
}

func (h *FinancialReportsHandler) GetTrialBalance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	report, err := h.svc.GetTrialBalance(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(report)
}

func (h *FinancialReportsHandler) ListTrialBalanceEntries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	entries, err := h.svc.ListTrialBalanceEntries(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(entries)
}

func (h *FinancialReportsHandler) ListTrialBalanceReports(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("org_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	reports, err := h.svc.ListTrialBalanceReports(r.Context(), orgID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(reports)
}

//
// ==========================
// Compliance
// ==========================
func (h *FinancialReportsHandler) GenerateCompliance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrganizationID string `json:"organization_id"`
		PeriodStart    string `json:"period_start"`
		PeriodEnd      string `json:"period_end"`
		Jurisdiction   string `json:"jurisdiction"`
		Details        string `json:"details"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	report, err := h.svc.GenerateCompliance(r.Context(), req.OrganizationID, req.PeriodStart, req.PeriodEnd, req.Jurisdiction, req.Details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(report)
}

func (h *FinancialReportsHandler) GetCompliance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	report, err := h.svc.GetCompliance(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(report)
}

func (h *FinancialReportsHandler) ListComplianceReports(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("org_id")
	jurisdiction := r.URL.Query().Get("jurisdiction")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	reports, err := h.svc.ListComplianceReports(r.Context(), orgID, jurisdiction, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(reports)
}
