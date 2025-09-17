package grpc_server_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	money "google.golang.org/genproto/googleapis/type/money"
	financepb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	grpcserver "github.com/ShristiRnr/Finance_mierp/internal/core/ports/grpc_server"
	"google.golang.org/protobuf/types/known/emptypb"
)

// --- Mock Service ---
type MockBudgetService struct {
	mock.Mock
}

func (m *MockBudgetService) CreateBudget(ctx context.Context, b *db.Budget) (*db.Budget, error) {
	args := m.Called(ctx, b)
	return args.Get(0).(*db.Budget), args.Error(1)
}

func (m *MockBudgetService) GetBudget(ctx context.Context, id uuid.UUID) (*db.Budget, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*db.Budget), args.Error(1)
}

func (m *MockBudgetService) ListBudgets(ctx context.Context, limit, offset int32) ([]db.Budget, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.Budget), args.Error(1)
}

func (m *MockBudgetService) UpdateBudget(ctx context.Context, b *db.Budget) (*db.Budget, error) {
	args := m.Called(ctx, b)
	return args.Get(0).(*db.Budget), args.Error(1)
}

func (m *MockBudgetService) DeleteBudget(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBudgetService) AllocateBudget(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error) {
	args := m.Called(ctx, ba)
	return args.Get(0).(*db.BudgetAllocation), args.Error(1)
}

func (m *MockBudgetService) GetBudgetAllocation(ctx context.Context, id uuid.UUID) (*db.BudgetAllocation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*db.BudgetAllocation), args.Error(1)
}

func (m *MockBudgetService) ListBudgetAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]db.BudgetAllocation, error) {
	args := m.Called(ctx, budgetID, limit, offset)
	return args.Get(0).([]db.BudgetAllocation), args.Error(1)
}

func (m *MockBudgetService) UpdateBudgetAllocation(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error) {
	args := m.Called(ctx, ba)
	return args.Get(0).(*db.BudgetAllocation), args.Error(1)
}

func (m *MockBudgetService) DeleteBudgetAllocation(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBudgetService) GetBudgetComparisonReport(ctx context.Context, id uuid.UUID) (*db.GetBudgetComparisonReportRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*db.GetBudgetComparisonReportRow), args.Error(1)
}

// --- Helper ---
func makeUUID() uuid.UUID {
	return uuid.New()
}

// --- Tests ---
func TestCreateBudget(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	expected := &db.Budget{
		ID:          makeUUID(),
		Name:        "Test Budget",
		TotalAmount: "1000", // should match moneyToString output
	}

	mockSvc.On("CreateBudget", mock.Anything, mock.MatchedBy(func(b *db.Budget) bool {
		return b.Name == "Test Budget" && b.TotalAmount == "1000"
	})).Return(expected, nil)

	req := &financepb.CreateBudgetRequest{
		Budget: &financepb.Budget{
			Name:        "Test Budget",
			TotalAmount: &money.Money{Units: 1000}, // Units=1000 produces "1000" in moneyToString
		},
	}

	resp, err := handler.CreateBudget(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID.String(), resp.Id)
	assert.Equal(t, expected.Name, resp.Name)
}


func TestGetBudget(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	id := makeUUID()
	expected := &db.Budget{ID: id, Name: "Budget1", TotalAmount: "500"}

	mockSvc.On("GetBudget", mock.Anything, id).Return(expected, nil)

	req := &financepb.GetBudgetRequest{Id: id.String()}
	resp, err := handler.GetBudget(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID.String(), resp.Id)
	assert.Equal(t, expected.Name, resp.Name)
}

func TestListBudgets(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	budgets := []db.Budget{
		{ID: makeUUID(), Name: "B1", TotalAmount: "100"},
		{ID: makeUUID(), Name: "B2", TotalAmount: "200"},
	}

	mockSvc.On("ListBudgets", mock.Anything, int32(50), int32(0)).Return(budgets, nil)

	req := &financepb.ListBudgetsRequest{
		Page: &financepb.PageRequest{PageSize: 50},
	}

	resp, err := handler.ListBudgets(context.Background(), req)
	assert.NoError(t, err)
	assert.Len(t, resp.Budgets, 2)
	assert.Equal(t, budgets[0].Name, resp.Budgets[0].Name)
}

func TestUpdateBudget(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	id := makeUUID()
	expected := &db.Budget{
		ID:          id,
		Name:        "Updated",
		TotalAmount: "300",
	}

	mockSvc.On("UpdateBudget", mock.Anything, mock.MatchedBy(func(b *db.Budget) bool {
		return b.ID == id && b.Name == "Updated" && b.TotalAmount == "300"
	})).Return(expected, nil)

	req := &financepb.UpdateBudgetRequest{
		Budget: &financepb.Budget{
			Id:          id.String(),
			Name:        "Updated",
			TotalAmount: &money.Money{Units: 300},
		},
	}

	resp, err := handler.UpdateBudget(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID.String(), resp.Id)
	assert.Equal(t, expected.Name, resp.Name)
	assert.Equal(t, expected.TotalAmount, "300")
}


func TestDeleteBudget(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	id := makeUUID()
	mockSvc.On("DeleteBudget", mock.Anything, id).Return(nil)

	req := &financepb.DeleteBudgetRequest{Id: id.String()}
	resp, err := handler.DeleteBudget(context.Background(), req)
	assert.NoError(t, err)
	assert.IsType(t, &emptypb.Empty{}, resp)
}

func TestAllocateBudget(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	expected := &db.BudgetAllocation{
		ID:             makeUUID(),
		BudgetID:       makeUUID(),
		AllocatedAmount: "500",
	}

	mockSvc.On("AllocateBudget", mock.Anything, mock.MatchedBy(func(ba *db.BudgetAllocation) bool {
		return ba.BudgetID == expected.BudgetID && ba.AllocatedAmount == "500"
	})).Return(expected, nil)

	req := &db.BudgetAllocation{
		BudgetID:       expected.BudgetID,
		AllocatedAmount: "500",
	}

	resp, err := handler.AllocateBudget(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.AllocatedAmount, resp.AllocatedAmount)
}

func TestGetBudgetAllocation(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	id := makeUUID()
	expected := &db.BudgetAllocation{ID: id, BudgetID: makeUUID(), AllocatedAmount: "200"}

	mockSvc.On("GetBudgetAllocation", mock.Anything, id).Return(expected, nil)

	resp, err := handler.GetBudgetAllocation(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.AllocatedAmount, resp.AllocatedAmount)
}

func TestUpdateBudgetAllocation(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	id := makeUUID()
	expected := &db.BudgetAllocation{ID: id, BudgetID: makeUUID(), AllocatedAmount: "800"}

	mockSvc.On("UpdateBudgetAllocation", mock.Anything, mock.MatchedBy(func(ba *db.BudgetAllocation) bool {
		return ba.ID == id && ba.AllocatedAmount == "800"
	})).Return(expected, nil)

	req := &db.BudgetAllocation{ID: id, AllocatedAmount: "800"}

	resp, err := handler.UpdateBudgetAllocation(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.AllocatedAmount, resp.AllocatedAmount)
}

func TestListBudgetAllocations(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	budgetID := makeUUID()
	allocs := []db.BudgetAllocation{
		{ID: makeUUID(), BudgetID: budgetID, AllocatedAmount: "100"},
		{ID: makeUUID(), BudgetID: budgetID, AllocatedAmount: "200"},
	}

	mockSvc.On("ListBudgetAllocations", mock.Anything, budgetID, int32(50), int32(0)).Return(allocs, nil)

	resp, err := handler.ListBudgetAllocations(context.Background(), budgetID, 50, 0)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, allocs[0].AllocatedAmount, resp[0].AllocatedAmount)
	assert.Equal(t, allocs[1].AllocatedAmount, resp[1].AllocatedAmount)
}

func TestDeleteBudgetAllocation(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	id := makeUUID()
	mockSvc.On("DeleteBudgetAllocation", mock.Anything, id).Return(nil)

	err := handler.DeleteBudgetAllocation(context.Background(), id)
	assert.NoError(t, err)
}


func TestGetBudgetComparison(t *testing.T) {
	mockSvc := new(MockBudgetService)
	handler := grpcserver.NewBudgetHandler(mockSvc)

	id := makeUUID()
	expected := &db.GetBudgetComparisonReportRow{BudgetID: id, TotalSpent: "500"}

	mockSvc.On("GetBudgetComparisonReport", mock.Anything, id).Return(expected, nil)

	resp, err := handler.GetBudgetComparison(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, expected.BudgetID, resp.BudgetID)
}
