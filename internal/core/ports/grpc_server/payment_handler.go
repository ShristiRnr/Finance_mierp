package grpc_server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
)

// Example for BankAccount handler
type BankAccountHandler struct {
	svc *services.BankService
}

func NewBankAccountHandler(svc *services.BankService) *BankAccountHandler {
	return &BankAccountHandler{svc: svc}
}

func (h *BankAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var ba db.BankAccount
	if err := json.NewDecoder(r.Body).Decode(&ba); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := h.svc.CreateBankAccount(r.Context(), ba)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

func (h *BankAccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}
	ba, err := h.svc.GetBankAccount(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(ba)
}

func (h *BankAccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}
	var ba db.BankAccount
	if err := json.NewDecoder(r.Body).Decode(&ba); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ba.ID = id
	updated, err := h.svc.UpdateBankAccount(r.Context(), ba)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (h *BankAccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	if err := h.svc.DeleteBankAccount(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BankAccountHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagingParams(r)
	items, err := h.svc.ListBankAccounts(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	json.NewEncoder(w).Encode(items)
}

type PaymentDueHandler struct {
	svc *services.PaymentDueService
}

func NewPaymentDueHandler(svc *services.PaymentDueService) *PaymentDueHandler {
	return &PaymentDueHandler{svc: svc}
}

func (h *PaymentDueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var pd db.PaymentDue
	if err := json.NewDecoder(r.Body).Decode(&pd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := h.svc.CreatePaymentDue(r.Context(), pd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

func (h *PaymentDueHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}
	pd, err := h.svc.GetPaymentDue(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pd)
}

func (h *PaymentDueHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}
	var pd db.PaymentDue
	if err := json.NewDecoder(r.Body).Decode(&pd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pd.ID = id
	updated, err := h.svc.UpdatePaymentDue(r.Context(), pd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (h *PaymentDueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	// Get user ID from request context or headers
	userID := r.Header.Get("X-User-ID") // adjust to your auth setup

	if err := h.svc.DeletePaymentDue(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (h *PaymentDueHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagingParams(r)
	items, err := h.svc.ListPaymentDues(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func (h *PaymentDueHandler) MarkAsPaid(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}
	var updatedBy string
	if err := json.NewDecoder(r.Body).Decode(&updatedBy); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paid, err := h.svc.MarkPaymentAsPaid(r.Context(), id, updatedBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(paid)
}

type BankTransactionHandler struct {
	svc *services.BankTransactionService
}

func NewBankTransactionHandler(svc *services.BankTransactionService) *BankTransactionHandler {
	return &BankTransactionHandler{svc: svc}
}

func (h *BankTransactionHandler) Import(w http.ResponseWriter, r *http.Request) {
	bankAccountID, err := uuid.Parse(chi.URLParam(r, "bank_account_id"))
	if err != nil {
		http.Error(w, "invalid bank account UUID", http.StatusBadRequest)
		return
	}

	var tx db.BankTransaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tx.BankAccountID = bankAccountID

	userID := r.Header.Get("X-User-ID") // pass the user ID

	created, err := h.svc.ImportBankTransaction(r.Context(), tx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}


func (h *BankTransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	bankAccountID, err := uuid.Parse(chi.URLParam(r, "bank_account_id"))
	if err != nil {
		http.Error(w, "invalid bank account UUID", http.StatusBadRequest)
		return
	}

	limit, offset := parsePagingParams(r)
	items, err := h.svc.ListBankTransactions(r.Context(), bankAccountID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func (h *BankTransactionHandler) Reconcile(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid transaction UUID", http.StatusBadRequest)
		return
	}

	var tx db.BankTransaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tx.ID = id

	userID := r.Header.Get("X-User-ID") // pass the user ID

	reconciled, err := h.svc.ReconcileTransaction(r.Context(), tx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reconciled)
}
