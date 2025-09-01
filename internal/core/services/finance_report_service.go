package services

import (
	"context"
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
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
) (domain.ProfitLossReport, error) {
	report := domain.ProfitLossReport{
		ID:            uuid.New(),
		OrganizationID: orgID,
		PeriodStart:   parseDate(start),
		PeriodEnd:     parseDate(end),
		TotalRevenue:  revenue,
		TotalExpenses: expenses,
		NetProfit:     revenue - expenses,
	}
	return s.repo.GenerateProfitLossReport(ctx, report)
}

func (s *FinancialReportsService) GetProfitLoss(ctx context.Context, id string) (domain.ProfitLossReport, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.ProfitLossReport{}, fmt.Errorf("invalid UUID: %w", err)
	}

	report, err := s.repo.GetProfitLossReport(ctx, uid)
	if err != nil {
		return domain.ProfitLossReport{}, err
	}

	return report, nil
}


func (s *FinancialReportsService) ListProfitLossReports(ctx context.Context, orgID string, limit, offset int32) ([]domain.ProfitLossReport, error) {
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
) (domain.BalanceSheetReport, error) {
	report := domain.BalanceSheetReport{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		PeriodStart:     parseDate(start),
		PeriodEnd:       parseDate(end),
		TotalAssets:     totalAssets,
		TotalLiabilities: totalLiabilities,
		NetWorth:        totalAssets - totalLiabilities,
	}
	return s.repo.GenerateBalanceSheetReport(ctx, report)
}

func (s *FinancialReportsService) GetBalanceSheet(ctx context.Context, id string) (domain.BalanceSheetReport, error) {
    uid, err := uuid.Parse(id)
    if err != nil {
        return domain.BalanceSheetReport{}, fmt.Errorf("invalid UUID: %w", err)
    }
    return s.repo.GetBalanceSheetReport(ctx, uid)
}


func (s *FinancialReportsService) ListBalanceSheetReports(ctx context.Context, orgID string, limit, offset int32) ([]domain.BalanceSheetReport, error) {
	return s.repo.ListBalanceSheetReports(ctx, orgID, limit, offset)
}

//
// ==========================
// Trial Balance
// ==========================
func (s *FinancialReportsService) CreateTrialBalance(
	ctx context.Context,
	orgID, start, end string,
) (domain.TrialBalanceReport, error) {
	report := domain.TrialBalanceReport{
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
) (domain.TrialBalanceEntry, error) {
	entry := domain.TrialBalanceEntry{
		ID:           uuid.New(),
		ReportID:     uuid.MustParse(reportID),
		LedgerAccount: ledgerAccount,
		Debit:        debit,
		Credit:       credit,
	}
	return s.repo.AddTrialBalanceEntry(ctx, entry)
}

func (s *FinancialReportsService) GetTrialBalance(ctx context.Context, id string) (domain.TrialBalanceReport, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.TrialBalanceReport{}, fmt.Errorf("invalid UUID: %w", err)
	}
	return s.repo.GetTrialBalanceReport(ctx, uid)
}

func (s *FinancialReportsService) ListTrialBalanceEntries(ctx context.Context, reportID string) ([]domain.TrialBalanceEntry, error) {
	uid, err := uuid.Parse(reportID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for reportID: %w", err)
	}

	return s.repo.ListTrialBalanceEntries(ctx, uid)
}


func (s *FinancialReportsService) ListTrialBalanceReports(ctx context.Context, orgID string, limit, offset int32) ([]domain.TrialBalanceReport, error) {
	return s.repo.ListTrialBalanceReports(ctx, orgID, limit, offset)
}

//
// ==========================
// Compliance
// ==========================
func (s *FinancialReportsService) GenerateCompliance(
	ctx context.Context,
	orgID, start, end, jurisdiction, details string,
) (domain.ComplianceReport, error) {
	report := domain.ComplianceReport{
		ID:            uuid.New(),
		OrganizationID: orgID,
		PeriodStart:   parseDate(start),
		PeriodEnd:     parseDate(end),
		Jurisdiction:  jurisdiction,
		Details:       details,
	}
	return s.repo.GenerateComplianceReport(ctx, report)
}

func (s *FinancialReportsService) GetCompliance(ctx context.Context, id string) (domain.ComplianceReport, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.ComplianceReport{}, fmt.Errorf("invalid UUID: %w", err)
	}
	return s.repo.GetComplianceReport(ctx, uid)
}

func (s *FinancialReportsService) ListComplianceReports(ctx context.Context, orgID, jurisdiction string, limit, offset int32) ([]domain.ComplianceReport, error) {
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
