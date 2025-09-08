package services

import (
	"context"
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type FinancialReportsService struct {
	repo ports.FinancialReportsRepository
}

func NewFinancialReportsService(repo ports.FinancialReportsRepository) *FinancialReportsService {
	return &FinancialReportsService{repo: repo}
}

//
// ==========================
// Profit & Loss
// ==========================
func (s *FinancialReportsService) GenerateProfitLoss(
	ctx context.Context,
	orgID string,
	start, end string,
	revenue, expenses float64,
) (db.ProfitLossReport, error) {
	report := db.ProfitLossReport{
		ID:            uuid.New(),
		OrganizationID: orgID,
		PeriodStart:   parseDate(start),
		PeriodEnd:     parseDate(end),
		TotalRevenue:  fmt.Sprintf("%.2f", revenue),
		TotalExpenses: fmt.Sprintf("%.2f", expenses),
		NetProfit:     fmt.Sprintf("%.2f", revenue-expenses),
	}
	return s.repo.GenerateProfitLossReport(ctx, report)
}

func (s *FinancialReportsService) GetProfitLoss(ctx context.Context, id string) (db.ProfitLossReport, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return db.ProfitLossReport{}, fmt.Errorf("invalid UUID: %w", err)
	}

	report, err := s.repo.GetProfitLossReport(ctx, uid)
	if err != nil {
		return db.ProfitLossReport{}, err
	}

	return report, nil
}


func (s *FinancialReportsService) ListProfitLossReports(ctx context.Context, orgID string, limit, offset int32) ([]db.ProfitLossReport, error) {
	return s.repo.ListProfitLossReports(ctx, orgID, limit, offset)
}

//
// ==========================
// Balance Sheet
// ==========================
func (s *FinancialReportsService) GenerateBalanceSheet(
	ctx context.Context,
	orgID string,
	start, end string,
	totalAssets, totalLiabilities float64,
) (db.BalanceSheetReport, error) {
	report := db.BalanceSheetReport{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		PeriodStart:     parseDate(start),
		PeriodEnd:       parseDate(end),
		TotalAssets:     fmt.Sprintf("%.2f", totalAssets),
		TotalLiabilities: fmt.Sprintf("%.2f", totalLiabilities),
		NetWorth:        fmt.Sprintf("%.2f", totalAssets-totalLiabilities),
	}
	return s.repo.GenerateBalanceSheetReport(ctx, report)
}

func (s *FinancialReportsService) GetBalanceSheet(ctx context.Context, id string) (db.BalanceSheetReport, error) {
    uid, err := uuid.Parse(id)
    if err != nil {
        return db.BalanceSheetReport{}, fmt.Errorf("invalid UUID: %w", err)
    }
    return s.repo.GetBalanceSheetReport(ctx, uid)
}


func (s *FinancialReportsService) ListBalanceSheetReports(ctx context.Context, orgID string, limit, offset int32) ([]db.BalanceSheetReport, error) {
	return s.repo.ListBalanceSheetReports(ctx, orgID, limit, offset)
}

//
// ==========================
// Trial Balance
// ==========================
func (s *FinancialReportsService) CreateTrialBalance(
	ctx context.Context,
	orgID, start, end string,
) (db.TrialBalanceReport, error) {
	report := db.TrialBalanceReport{
		ID:            uuid.New(),
		OrganizationID: orgID,
		PeriodStart:   parseDate(start),
		PeriodEnd:     parseDate(end),
	}
	return s.repo.CreateTrialBalanceReport(ctx, report)
}

func (s *FinancialReportsService) AddTrialBalanceEntry(
	ctx context.Context,
	reportID, ledgerAccount string,
	debit, credit float64,
) (db.TrialBalanceEntry, error) {
	entry := db.TrialBalanceEntry{
		ID:           uuid.New(),
		ReportID:     uuid.MustParse(reportID),
		LedgerAccount: ledgerAccount,
		Debit:         fmt.Sprintf("%.2f", debit),
    	Credit:        fmt.Sprintf("%.2f", credit),	
	}
	return s.repo.AddTrialBalanceEntry(ctx, entry)
}

func (s *FinancialReportsService) GetTrialBalance(ctx context.Context, id string) (db.TrialBalanceReport, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return db.TrialBalanceReport{}, fmt.Errorf("invalid UUID: %w", err)
	}
	return s.repo.GetTrialBalanceReport(ctx, uid)
}

func (s *FinancialReportsService) ListTrialBalanceEntries(ctx context.Context, reportID string) ([]db.TrialBalanceEntry, error) {
	uid, err := uuid.Parse(reportID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for reportID: %w", err)
	}

	return s.repo.ListTrialBalanceEntries(ctx, uid)
}


func (s *FinancialReportsService) ListTrialBalanceReports(ctx context.Context, orgID string, limit, offset int32) ([]db.TrialBalanceReport, error) {
	return s.repo.ListTrialBalanceReports(ctx, orgID, limit, offset)
}

//
// ==========================
// Compliance
// ==========================
func (s *FinancialReportsService) GenerateCompliance(
	ctx context.Context,
	orgID, start, end, jurisdiction, details string,
) (db.ComplianceReport, error) {
	report := db.ComplianceReport{
		ID:            uuid.New(),
		OrganizationID: orgID,
		PeriodStart:   parseDate(start),
		PeriodEnd:     parseDate(end),
		Jurisdiction:  jurisdiction,
		Details:       details,
	}
	return s.repo.GenerateComplianceReport(ctx, report)
}

func (s *FinancialReportsService) GetCompliance(ctx context.Context, id string) (db.ComplianceReport, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return db.ComplianceReport{}, fmt.Errorf("invalid UUID: %w", err)
	}
	return s.repo.GetComplianceReport(ctx, uid)
}

func (s *FinancialReportsService) ListComplianceReports(ctx context.Context, orgID, jurisdiction string, limit, offset int32) ([]db.ComplianceReport, error) {
	return s.repo.ListComplianceReports(ctx, orgID, jurisdiction, limit, offset)
}

//
// ==========================
// Helpers
// ==========================
func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}
