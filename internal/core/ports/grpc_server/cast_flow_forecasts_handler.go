package grpc_server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	"encoding/json"
)

type CashFlowHandler struct {
	service *services.CashFlowService
}

func NewCashFlowHandler(s *services.CashFlowService) *CashFlowHandler {
	return &CashFlowHandler{service: s}
}

// HTTP Example: Generate Forecast
func (h *CashFlowHandler) GenerateForecast(w http.ResponseWriter, r *http.Request) {
	var cf domain.CashFlowForecast
	if err := json.NewDecoder(r.Body).Decode(&cf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.GenerateForecast(r.Context(), &cf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// HTTP Example: Get Forecast
func (h *CashFlowHandler) GetForecast(w http.ResponseWriter, r *http.Request, id string) {
	uid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	result, err := h.service.GetForecast(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// HTTP Example: List Forecasts
func (h *CashFlowHandler) ListForecasts(w http.ResponseWriter, r *http.Request, organizationID string, limit, offset int32) {
	result, err := h.service.ListForecasts(r.Context(), organizationID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
