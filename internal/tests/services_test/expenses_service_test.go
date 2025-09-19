package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================
// Mock Repositories
// ============================
type MockExpenseRepo struct{ mock.Mock }

func (m *MockExpenseRepo) Create(ctx context.Context, exp db.Expense) (db.Expense, error) {
	args := m.Called(ctx, exp)
	return args.Get(0).(db.Expense), args.Error(1)
}
func (m *MockExpenseRepo) Get(ctx context.Context, id uuid.UUID) (db.Expense, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Expense), args.Error(1)
}
func (m *MockExpenseRepo) List(ctx context.Context, limit, offset int32) ([]db.Expense, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.Expense), args.Error(1)
}
func (m *MockExpenseRepo) Update(ctx context.Context, exp db.Expense) (db.Expense, error) {
	args := m.Called(ctx, exp)
	return args.Get(0).(db.Expense), args.Error(1)
}
func (m *MockExpenseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Mock CostCenter Repo ---
type MockCostCenterRepo struct{ mock.Mock }

func (m *MockCostCenterRepo) Create(ctx context.Context, cc db.CostCenter) (db.CostCenter, error) {
	args := m.Called(ctx, cc)
	return args.Get(0).(db.CostCenter), args.Error(1)
}
func (m *MockCostCenterRepo) Get(ctx context.Context, id uuid.UUID) (db.CostCenter, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.CostCenter), args.Error(1)
}
func (m *MockCostCenterRepo) List(ctx context.Context, limit, offset int32) ([]db.CostCenter, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.CostCenter), args.Error(1)
}
func (m *MockCostCenterRepo) Update(ctx context.Context, cc db.CostCenter) (db.CostCenter, error) {
	args := m.Called(ctx, cc)
	return args.Get(0).(db.CostCenter), args.Error(1)
}
func (m *MockCostCenterRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Mock CostAllocation Repo ---
type MockCostAllocationRepo struct{ mock.Mock }

func (m *MockCostAllocationRepo) Allocate(ctx context.Context, ca db.CostAllocation) (db.CostAllocation, error) {
	args := m.Called(ctx, ca)
	return args.Get(0).(db.CostAllocation), args.Error(1)
}
func (m *MockCostAllocationRepo) List(ctx context.Context, limit, offset int32) ([]db.CostAllocation, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.CostAllocation), args.Error(1)
}

// ============================
// Mock Publisher
// ============================
type MockoPublisher struct{ mock.Mock }

func (m *MockoPublisher) PublishExpenseCreated(ctx context.Context, e *db.Expense) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockoPublisher) PublishExpenseUpdated(ctx context.Context, e *db.Expense) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockoPublisher) PublishExpenseDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockoPublisher) PublishCostCenterCreated(ctx context.Context, cc *db.CostCenter) error {
	return m.Called(ctx, cc).Error(0)
}
func (m *MockoPublisher) PublishCostCenterUpdated(ctx context.Context, cc *db.CostCenter) error {
	return m.Called(ctx, cc).Error(0)
}
func (m *MockoPublisher) PublishCostCenterDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockoPublisher) PublishCostAllocationAllocated(ctx context.Context, ca *db.CostAllocation) error {
	return m.Called(ctx, ca).Error(0)
}
func (m *MockoPublisher) PublishCostAllocationListed(ctx context.Context, ca []db.CostAllocation) error {
	return m.Called(ctx, ca).Error(0)
}

// ============================
// ExpenseService Tests
// ============================
func TestExpenseService_CRUD(t *testing.T) {
	ctx := context.Background()
	expRepo := new(MockExpenseRepo)
	pub := new(MockPublisher)
	svc := services.NewExpenseService(expRepo, pub)

	exp := db.Expense{ID: uuid.New(), Category: "Travel", Amount: "1000", ExpenseDate: time.Now()}

	// Create
	expRepo.On("Create", ctx, exp).Return(exp, nil)
	pub.On("PublishExpenseCreated", ctx, &exp).Return(nil)

	got, err := svc.CreateExpense(ctx, exp)
	assert.NoError(t, err)
	assert.Equal(t, exp.ID, got.ID)

	// Get
	expRepo.On("Get", ctx, exp.ID).Return(exp, nil)
	got, err = svc.GetExpense(ctx, exp.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Travel", got.Category)

	// List
	expRepo.On("List", ctx, int32(10), int32(0)).Return([]db.Expense{exp}, nil)
	list, err := svc.ListExpenses(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// Update
	updated := exp
	updated.Amount = "2000"
	expRepo.On("Update", ctx, updated).Return(updated, nil)
	pub.On("PublishExpenseUpdated", ctx, &updated).Return(nil)

	got, err = svc.UpdateExpense(ctx, updated)
	assert.NoError(t, err)
	assert.Equal(t, "2000", got.Amount)

	// Delete
	expRepo.On("Delete", ctx, exp.ID).Return(nil)
	pub.On("PublishExpenseDeleted", ctx, exp.ID.String()).Return(nil)

	err = svc.DeleteExpense(ctx, exp.ID)
	assert.NoError(t, err)
}

// ============================
// CostCenterService Tests
// ============================
func TestCostCenterService_CRUD(t *testing.T) {
	ctx := context.Background()
	ccRepo := new(MockCostCenterRepo)
	pub := new(MockPublisher)
	svc := services.NewCostCenterService(ccRepo, pub)

	cc := db.CostCenter{ID: uuid.New(), Name: "IT Dept"}

	// Create
	ccRepo.On("Create", ctx, cc).Return(cc, nil)
	pub.On("PublishCostCenterCreated", ctx, &cc).Return(nil)
	created, err := svc.CreateCostCenter(ctx, cc)
	assert.NoError(t, err)
	assert.Equal(t, cc.ID, created.ID)

	// Get
	ccRepo.On("Get", ctx, cc.ID).Return(cc, nil)
	got, err := svc.GetCostCenter(ctx, cc.ID)
	assert.NoError(t, err)
	assert.Equal(t, "IT Dept", got.Name)

	// List
	ccRepo.On("List", ctx, int32(5), int32(0)).Return([]db.CostCenter{cc}, nil)
	list, err := svc.ListCostCenters(ctx, 5, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// Update
	updated := cc
	updated.Name = "Finance Dept"
	ccRepo.On("Update", ctx, updated).Return(updated, nil)
	pub.On("PublishCostCenterUpdated", ctx, &updated).Return(nil)

	got, err = svc.UpdateCostCenter(ctx, updated)
	assert.NoError(t, err)
	assert.Equal(t, "Finance Dept", got.Name)

	// Delete
	ccRepo.On("Delete", ctx, cc.ID).Return(nil)
	pub.On("PublishCostCenterDeleted", ctx, cc.ID.String()).Return(nil)
	err = svc.DeleteCostCenter(ctx, cc.ID)
	assert.NoError(t, err)
}

// ============================
// CostAllocationService Tests
// ============================
func TestCostAllocationService(t *testing.T) {
    ctx := context.Background()
    caRepo := new(MockCostAllocationRepo)
    pub := new(MockPublisher)

    // Use constructor
    svc := services.NewCostAllocationService(caRepo, pub)

    ca := db.CostAllocation{ID: uuid.New(), ReferenceType: "Invoice", ReferenceID: "inv-1", Amount: "500"}

    // Allocate
    caRepo.On("Allocate", ctx, ca).Return(ca, nil)
    pub.On("PublishCostAllocationAllocated", ctx, &ca).Return(nil)

    got, err := svc.AllocateCost(ctx, ca)
    assert.NoError(t, err)
    assert.Equal(t, ca.ID, got.ID)

    // List
    caRepo.On("List", ctx, int32(5), int32(0)).Return([]db.CostAllocation{ca}, nil)
    pub.On("PublishCostAllocationListed", ctx, []db.CostAllocation{ca}).Return(nil)

    list, err := svc.ListAllocations(ctx, 5, 0)
    assert.NoError(t, err)
    assert.Len(t, list, 1)
}

