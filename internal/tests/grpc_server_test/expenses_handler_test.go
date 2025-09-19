package grpc_server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports/grpc_server"
)

// =============================
// Mock Services
// =============================
type MockExpenseService struct{ mock.Mock }

func (m *MockExpenseService) CreateExpense(ctx context.Context, e db.Expense) (db.Expense, error) {
	args := m.Called(ctx, e)
	return args.Get(0).(db.Expense), args.Error(1)
}
func (m *MockExpenseService) GetExpense(ctx context.Context, id uuid.UUID) (db.Expense, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Expense), args.Error(1)
}
func (m *MockExpenseService) ListExpenses(ctx context.Context, limit, offset int32) ([]db.Expense, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.Expense), args.Error(1)
}
func (m *MockExpenseService) UpdateExpense(ctx context.Context, e db.Expense) (db.Expense, error) {
	args := m.Called(ctx, e)
	return args.Get(0).(db.Expense), args.Error(1)
}
func (m *MockExpenseService) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// =============================
// Test ExpenseHandler
// =============================
func setupRouter(h *grpc_server.ExpenseHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/expenses", h.CreateExpense)
	r.GET("/expenses/:id", h.GetExpense)
	r.GET("/expenses", h.ListExpenses)
	r.PUT("/expenses/:id", h.UpdateExpense)
	r.DELETE("/expenses/:id", h.DeleteExpense)
	return r
}

func TestCreateExpense(t *testing.T) {
	mockSvc := new(MockExpenseService)
	handler := grpc_server.NewExpenseHandler(mockSvc)
	router := setupRouter(handler)

	exp := db.Expense{
		ID:          uuid.New(),
		Category:    "Travel",
		Amount:      "1000",
		ExpenseDate: time.Now(),
	}
	mockSvc.On("CreateExpense", mock.Anything, mock.AnythingOfType("db.Expense")).Return(exp, nil)

	body, _ := json.Marshal(exp)
	req, _ := http.NewRequest("POST", "/expenses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var got db.Expense
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, exp.Category, got.Category)
	mockSvc.AssertExpectations(t)
}

func TestGetExpense(t *testing.T) {
	mockSvc := new(MockExpenseService)
	handler := grpc_server.NewExpenseHandler(mockSvc)
	router := setupRouter(handler)

	id := uuid.New()
	exp := db.Expense{ID: id, Category: "Food", Amount: "500"}
	mockSvc.On("GetExpense", mock.Anything, id).Return(exp, nil)

	req, _ := http.NewRequest("GET", "/expenses/"+id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got db.Expense
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, exp.ID, got.ID)
	mockSvc.AssertExpectations(t)
}

func TestListExpenses(t *testing.T) {
	mockSvc := new(MockExpenseService)
	handler := grpc_server.NewExpenseHandler(mockSvc)
	router := setupRouter(handler)

	expenses := []db.Expense{
		{ID: uuid.New(), Category: "Rent", Amount: "2000"},
		{ID: uuid.New(), Category: "Utilities", Amount: "800"},
	}
	mockSvc.On("ListExpenses", mock.Anything, int32(50), int32(0)).Return(expenses, nil)

	req, _ := http.NewRequest("GET", "/expenses", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []db.Expense
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mockSvc.AssertExpectations(t)
}

func TestUpdateExpense(t *testing.T) {
	mockSvc := new(MockExpenseService)
	handler := grpc_server.NewExpenseHandler(mockSvc)
	router := setupRouter(handler)

	id := uuid.New()
	exp := db.Expense{ID: id, Category: "Health", Amount: "1200"}
	mockSvc.On("UpdateExpense", mock.Anything, exp).Return(exp, nil)

	body, _ := json.Marshal(exp)
	req, _ := http.NewRequest("PUT", "/expenses/"+id.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got db.Expense
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, exp.Category, got.Category)
	mockSvc.AssertExpectations(t)
}

func TestDeleteExpense(t *testing.T) {
	mockSvc := new(MockExpenseService)
	handler := grpc_server.NewExpenseHandler(mockSvc)
	router := setupRouter(handler)

	id := uuid.New()
	mockSvc.On("DeleteExpense", mock.Anything, id).Return(nil)

	req, _ := http.NewRequest("DELETE", "/expenses/"+id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}
