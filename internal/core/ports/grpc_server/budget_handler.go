package grpc_server

import (
	"context"
	"strconv"
    "database/sql"

	financepb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"github.com/google/uuid"
    "google.golang.org/genproto/googleapis/type/money"
)

type BudgetHandler struct {
	financepb.UnimplementedBudgetServiceServer
	service ports.BudgetService
}

func NewBudgetHandler(service ports.BudgetService) *BudgetHandler {
	return &BudgetHandler{service: service}
}

func dbBudgetToPb(b *db.Budget) *financepb.Budget {
    if b == nil {
        return nil
    }
    return &financepb.Budget{
        Id:          b.ID.String(),
        Name:        b.Name,
        TotalAmount: stringToMoney(b.TotalAmount, "USD"),
        Status:      b.Status,
    }
}

// Protobuf → DB
func pbBudgetToDb(pbB *financepb.Budget) *db.Budget {
    if pbB == nil {
        return nil
    }
    id, _ := uuid.Parse(pbB.GetId())
    return &db.Budget{
        ID:          id,
        Name:        pbB.GetName(),
        TotalAmount: moneyToString(pbB.GetTotalAmount()),
        Status:      pbB.GetStatus(),
        CreatedAt:   sql.NullTime{},   // set in service/repo
        CreatedBy:   sql.NullString{}, // set in service/repo
    }
}

// Helpers
func moneyToString(m *money.Money) string {
    if m == nil {
        return "0"
    }
    return strconv.FormatInt(m.Units, 10)
}

func stringToMoney(amount string, currency string) *money.Money {
    units, err := strconv.ParseInt(amount, 10, 64)
    if err != nil {
        units = 0
    }
    return &money.Money{CurrencyCode: currency, Units: units}
}

// ---------------- Budgets ----------------

func (h *BudgetHandler) CreateBudget(ctx context.Context, req *financepb.CreateBudgetRequest) (*financepb.Budget, error) {
    dbBudget := pbBudgetToDb(req.GetBudget())

    created, err := h.service.CreateBudget(ctx, dbBudget)
    if err != nil {
        return nil, err
    }

    return dbBudgetToPb(created), nil
}



func (h *BudgetHandler) GetBudget(ctx context.Context, req *financepb.GetBudgetRequest) (*financepb.Budget, error) {
	// Convert protobuf request ID string to UUID
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid budget id: %v", err)
	}

	// Call internal service
	b, err := h.service.GetBudget(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "budget not found: %v", err)
	}

	// Map internal db.Budget to protobuf response
	pbBudget := &financepb.Budget{
		Id:          b.ID.String(),
		Name:        b.Name,
		TotalAmount: stringToMoney(b.TotalAmount, "USD"), // adjust if using money proto
		Status:      b.Status,
	}

	return pbBudget, nil
}


func (h *BudgetHandler) ListBudgets(ctx context.Context, req *financepb.ListBudgetsRequest) (*financepb.ListBudgetsResponse, error) {
    limit := req.GetPage().GetPageSize()
    if limit == 0 {
        limit = 50
    }

    offset := int32(0)
    if req.GetPage().GetPageToken() != "" {
        o, err := strconv.Atoi(req.GetPage().GetPageToken())
        if err != nil {
            return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
        }
        offset = int32(o)
    }

    // Call internal service
    budgets, err := h.service.ListBudgets(ctx, limit, offset)
    if err != nil {
        return nil, err
    }

    // Convert internal db.Budget slice to protobuf
    pbBudgets := make([]*financepb.Budget, len(budgets))
    for i, b := range budgets {
        pbBudgets[i] = &financepb.Budget{
            Id:          b.ID.String(),
            Name:        b.Name,
            TotalAmount: stringToMoney(b.TotalAmount, "USD"), // adjust for money proto
            Status:      b.Status,
        }
    }

    // Compute next page token
    nextToken := ""
    if len(budgets) == int(limit) {
        nextToken = strconv.Itoa(int(offset + limit))
    }

    return &financepb.ListBudgetsResponse{
        Budgets: pbBudgets,
        Page: &financepb.PageResponse{
            NextPageToken: nextToken,
            TotalSize:     int64(len(budgets)), // adjust if you have total count
        },
    }, nil
}


func (h *BudgetHandler) UpdateBudget(ctx context.Context, req *financepb.UpdateBudgetRequest) (*financepb.Budget, error) {
    // Map protobuf request to internal db.Budget
    b := &db.Budget{
        ID:          uuid.MustParse(req.GetBudget().GetId()),
        Name:        req.GetBudget().GetName(),
        TotalAmount: moneyToString(req.GetBudget().GetTotalAmount()),
        Status:      req.GetBudget().GetStatus(),
    }

    updatedBudget, err := h.service.UpdateBudget(ctx, b)
    if err != nil {
        return nil, err
    }

    // Map db.Budget to protobuf response
    pbBudget := &financepb.Budget{
        Id:          updatedBudget.ID.String(),
        Name:        updatedBudget.Name,
        TotalAmount: stringToMoney(updatedBudget.TotalAmount, "USD"),
        Status:      updatedBudget.Status,
    }

    return pbBudget, nil
}


func (h *BudgetHandler) DeleteBudget(ctx context.Context, req *financepb.DeleteBudgetRequest) (*emptypb.Empty, error) {
	// Convert protobuf ID string to uuid
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid budget id: %v", err)
	}

	// Call internal service
	if err := h.service.DeleteBudget(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete budget: %v", err)
	}

	// Return empty response for gRPC
	return &emptypb.Empty{}, nil
}

// ---------------- Budget Allocations ----------------

func (h *BudgetHandler) AllocateBudget(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error) {
	return h.service.AllocateBudget(ctx, ba)
}

func (h *BudgetHandler) GetBudgetAllocation(ctx context.Context, id uuid.UUID) (*db.BudgetAllocation, error) {
	return h.service.GetBudgetAllocation(ctx, id)
}

func (h *BudgetHandler) ListBudgetAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]db.BudgetAllocation, error) {
	return h.service.ListBudgetAllocations(ctx, budgetID, limit, offset)
}


func (h *BudgetHandler) UpdateBudgetAllocation(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error) {
	return h.service.UpdateBudgetAllocation(ctx, ba)
}

func (h *BudgetHandler) DeleteBudgetAllocation(ctx context.Context, id uuid.UUID) error {
	return h.service.DeleteBudgetAllocation(ctx, id)
}

// ---------------- Budget Comparison ----------------

func (h *BudgetHandler) GetBudgetComparison(ctx context.Context, id uuid.UUID) (*db.GetBudgetComparisonReportRow, error) {
	return h.service.GetBudgetComparisonReport(ctx, id)
}
