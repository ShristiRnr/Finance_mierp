package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
)

// ========================
// Mock Queries
// ========================
type mockExpenseQueries struct {
	*db.Queries
}

func (m *mockExpenseQueries) CreateExpense(ctx context.Context, arg db.CreateExpenseParams) (db.Expense, error) {
	return db.Expense{
		ID:           uuid.New(),
		Category:     arg.Category,
		Amount:       arg.Amount,
		ExpenseDate:  arg.ExpenseDate,
		CostCenterID: arg.CostCenterID,
		CreatedBy:    arg.CreatedBy,
		UpdatedBy:    arg.UpdatedBy,
		Revision:     sql.NullInt32{Int32: 1, Valid: true},
	}, nil
}

func (m *mockExpenseQueries) GetExpense(ctx context.Context, id uuid.UUID) (db.Expense, error) {
	return db.Expense{ID: id, Category: "Travel"}, nil
}

func (m *mockExpenseQueries) ListExpenses(ctx context.Context, arg db.ListExpensesParams) ([]db.Expense, error) {
	return []db.Expense{
		{ID: uuid.New(), Category: "Travel"},
		{ID: uuid.New(), Category: "Food"},
	}, nil
}

func (m *mockExpenseQueries) UpdateExpense(ctx context.Context, arg db.UpdateExpenseParams) (db.Expense, error) {
	return db.Expense{
		ID:          arg.ID,
		Category:    arg.Category,
		Amount:      arg.Amount,
		ExpenseDate: arg.ExpenseDate,
		UpdatedBy:   arg.UpdatedBy,
	}, nil
}

func (m *mockExpenseQueries) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	return nil
}

// ========================
// Tests: Expense Repo
// ========================
func TestExpenseRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	mockQ := &mockExpenseQueries{}
	repo := repository.NewExpenseRepo(mockQ)

	exp := db.Expense{
		Category:    "Travel",
		Amount:      "1000",
		ExpenseDate: time.Now(),
	}

	// --- Create ---
	created, err := repo.Create(ctx, exp)
	assert.NoError(t, err)
	assert.Equal(t, "Travel", created.Category)

	// --- Get ---
	fetched, err := repo.Get(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)

	// --- List ---
	list, err := repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	// --- Update ---
	exp.ID = created.ID
	exp.Category = "Updated"
	updated, err := repo.Update(ctx, exp)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.Category)

	// --- Delete ---
	err = repo.Delete(ctx, created.ID)
	assert.NoError(t, err)
}

// ========================
// CostCenter Mock
// ========================
type mockCostCenterQueries struct{}

func (m *mockCostCenterQueries) CreateCostCenter(ctx context.Context, arg db.CreateCostCenterParams) (db.CostCenter, error) {
	return db.CostCenter{ID: uuid.New(), Name: arg.Name}, nil
}
func (m *mockCostCenterQueries) GetCostCenter(ctx context.Context, id uuid.UUID) (db.CostCenter, error) {
	return db.CostCenter{ID: id, Name: "CC1"}, nil
}
func (m *mockCostCenterQueries) ListCostCenters(ctx context.Context, arg db.ListCostCentersParams) ([]db.CostCenter, error) {
	return []db.CostCenter{{ID: uuid.New(), Name: "CC1"}}, nil
}
func (m *mockCostCenterQueries) UpdateCostCenter(ctx context.Context, arg db.UpdateCostCenterParams) (db.CostCenter, error) {
	return db.CostCenter{ID: arg.ID, Name: arg.Name}, nil
}
func (m *mockCostCenterQueries) DeleteCostCenter(ctx context.Context, id uuid.UUID) error {
	return nil
}

// ========================
// Tests: CostCenter Repo
// ========================
func TestCostCenterRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	mockQ := &mockCostCenterQueries{}
	repo := repository.NewCostCenterRepo(mockQ)

	cc := db.CostCenter{Name: "MainCC"}

	// Create
	created, err := repo.Create(ctx, cc)
	assert.NoError(t, err)
	assert.Equal(t, "MainCC", created.Name)

	// Get
	fetched, err := repo.Get(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)

	// List
	list, err := repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// Update
	cc.ID = created.ID
	cc.Name = "UpdatedCC"
	updated, err := repo.Update(ctx, cc)
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedCC", updated.Name)

	// Delete
	err = repo.Delete(ctx, created.ID)
	assert.NoError(t, err)
}

// ========================
// CostAllocation Mock
// ========================
type mockCostAllocationQueries struct{}

func (m *mockCostAllocationQueries) AllocateCost(ctx context.Context, arg db.AllocateCostParams) (db.CostAllocation, error) {
	return db.CostAllocation{ID: uuid.New(), CostCenterID: arg.CostCenterID, Amount: arg.Amount}, nil
}
func (m *mockCostAllocationQueries) ListCostAllocations(ctx context.Context, arg db.ListCostAllocationsParams) ([]db.CostAllocation, error) {
	return []db.CostAllocation{{ID: uuid.New()}}, nil
}

// ========================
// Tests: CostAllocation Repo
// ========================
func TestCostAllocationRepo(t *testing.T) {
	ctx := context.Background()
	mockQ := &mockCostAllocationQueries{}
	repo := repository.NewCostAllocationRepo(mockQ)

	ca := db.CostAllocation{
		CostCenterID: uuid.New(),
		Amount:       "500",
		ReferenceType:"EXPENSE",
		ReferenceID:  "EXP123",
	}

	// Allocate
	allocated, err := repo.Allocate(ctx, ca)
	assert.NoError(t, err)
	assert.Equal(t, "500", allocated.Amount)

	// List
	list, err := repo.List(ctx, 5, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}
