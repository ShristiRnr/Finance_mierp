package grpc_server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/go-chi/chi/v5"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports/grpc_server"
)

// -------------------- Mock Services --------------------

type MockBankService struct {
	mock.Mock
}

func (m *MockBankService) CreateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error) {
	args := m.Called(ctx, ba)
	return args.Get(0).(db.BankAccount), args.Error(1)
}

func (m *MockBankService) GetBankAccount(ctx context.Context, id uuid.UUID) (db.BankAccount, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.BankAccount), args.Error(1)
}

func (m *MockBankService) UpdateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error) {
	args := m.Called(ctx, ba)
	return args.Get(0).(db.BankAccount), args.Error(1)
}

func (m *MockBankService) DeleteBankAccount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBankService) ListBankAccounts(ctx context.Context, limit, offset int32) ([]db.BankAccount, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.BankAccount), args.Error(1)
}

// -------------------- Helper: Attach chi Route Params --------------------

func attachChiParams(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

// -------------------- BankAccountHandler Tests --------------------

func TestBankAccountHandler(t *testing.T) {
	mockSvc := new(MockBankService)
	handler := grpc_server.NewBankAccountHandler(mockSvc)

	baID := uuid.New()
	ba := db.BankAccount{
		ID:            baID,
		Name:          "Test Bank",
		AccountNumber: "1234567890",
		IfscOrSwift:   "IFSC0001",
	}

	t.Run("Create Bank Account", func(t *testing.T) {
		mockSvc.On("CreateBankAccount", mock.Anything, mock.Anything).Return(ba, nil)

		body, _ := json.Marshal(ba)
		req := httptest.NewRequest(http.MethodPost, "/bank_accounts", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var resp db.BankAccount
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, ba.ID, resp.ID)
	})

	t.Run("Get Bank Account", func(t *testing.T) {
		mockSvc.On("GetBankAccount", mock.Anything, baID).Return(ba, nil)

		req := httptest.NewRequest(http.MethodGet, "/bank_accounts/"+baID.String(), nil)
		req = attachChiParams(req, "id", baID.String())
		w := httptest.NewRecorder()

		handler.Get(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("Update Bank Account", func(t *testing.T) {
		mockSvc.On("UpdateBankAccount", mock.Anything, mock.Anything).Return(ba, nil)

		body, _ := json.Marshal(ba)
		req := httptest.NewRequest(http.MethodPut, "/bank_accounts/"+baID.String(), bytes.NewReader(body))
		req = attachChiParams(req, "id", baID.String())
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("Delete Bank Account", func(t *testing.T) {
		mockSvc.On("DeleteBankAccount", mock.Anything, baID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/bank_accounts/"+baID.String(), nil)
		req = attachChiParams(req, "id", baID.String())
		w := httptest.NewRecorder()

		handler.Delete(w, req)
		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})

	t.Run("List Bank Accounts", func(t *testing.T) {
		mockSvc.On("ListBankAccounts", mock.Anything, int32(100), int32(0)).Return([]db.BankAccount{ba}, nil)

		req := httptest.NewRequest(http.MethodGet, "/bank_accounts", nil)
		w := httptest.NewRecorder()

		handler.List(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}
