package repository

import (
	"context"
	"strconv"
	"fmt"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
)

type FinancialReportsRepo struct {
	q *db.Queries
}

func NewFinancialReportsRepo(q *db.Queries) ports.FinancialReportsRepository {
	return &FinancialReportsRepo{q: q}
}


// =============================================== Mapping functions ===============================================

func mapProfitLossToDomain(r db.ProfitLossReport) domain.ProfitLossReport {
	revenue, err := strconv.ParseFloat(r.TotalRevenue, 64)
	if err != nil {
		revenue = 0
	}

	expense, err := strconv.ParseFloat(r.TotalExpenses, 64)
	if err != nil {
		expense = 0
	}

	profit, err := strconv.ParseFloat(r.NetProfit, 64)
	if err != nil {
		profit = 0
	}

	return domain.ProfitLossReport{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		PeriodStart:    r.PeriodStart,
		PeriodEnd:      r.PeriodEnd,
		TotalRevenue:   revenue,
		TotalExpenses:  expense,
		NetProfit:      profit,
		CreatedAt:      r.CreatedAt.Time,
		CreatedBy:      r.CreatedBy.String,
		UpdatedAt:      r.UpdatedAt.Time,
		UpdatedBy:      r.UpdatedBy.String,
		Revision:       r.Revision.Int32,
	}
}


func mapBalanceSheetToDomain(r db.BalanceSheetReport) domain.BalanceSheetReport {
	asset, err := strconv.ParseFloat(r.TotalAssets, 64)
	if err != nil {
		asset = 0 // or log the error
	}

	liability, err := strconv.ParseFloat(r.TotalLiabilities, 64)
	if err != nil {
		liability = 0 // or log the error
	}

	networth, err := strconv.ParseFloat(r.NetWorth, 64)
	if err != nil {
		networth = 0 // or log the error
	}

	return domain.BalanceSheetReport{
		ID:              r.ID,
		OrganizationID:  r.OrganizationID,
		PeriodStart:     r.PeriodStart,
		PeriodEnd:       r.PeriodEnd,
		TotalAssets:     asset,
		TotalLiabilities: liability,
		NetWorth:        networth,
		CreatedAt:       r.CreatedAt.Time,
		CreatedBy:       r.CreatedBy.String,
		UpdatedAt:       r.UpdatedAt.Time,
		UpdatedBy:       r.UpdatedBy.String,
		Revision:        r.Revision.Int32,
	}
}


func mapTrialBalanceToDomain(r db.TrialBalanceReport) domain.TrialBalanceReport {
	return domain.TrialBalanceReport{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		PeriodStart:    r.PeriodStart,
		PeriodEnd:      r.PeriodEnd,
		CreatedAt:      r.CreatedAt.Time,
		CreatedBy:      r.CreatedBy.String,
		UpdatedAt:      r.UpdatedAt.Time,
		UpdatedBy:      r.UpdatedBy.String,
		Revision:       r.Revision.Int32,
	}
}

func mapTrialBalanceEntryToDomain(r db.TrialBalanceEntry) domain.TrialBalanceEntry {
	debit, err := strconv.ParseFloat(r.Debit, 64)
	if err != nil {
		debit = 0 // or handle/log the error
	}

	credit, err := strconv.ParseFloat(r.Credit, 64)
	if err != nil {
		credit = 0 // or handle/log the error
	}

	return domain.TrialBalanceEntry{
		ID:            r.ID,
		ReportID:      r.ReportID,
		LedgerAccount: r.LedgerAccount,
		Debit:         debit,
		Credit:        credit,
		CreatedAt:     r.CreatedAt.Time,
		CreatedBy:     r.CreatedBy.String,
	}
}


func mapComplianceToDomain(r db.ComplianceReport) domain.ComplianceReport {
	return domain.ComplianceReport{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		PeriodStart:    r.PeriodStart,
		PeriodEnd:      r.PeriodEnd,
		Jurisdiction:   r.Jurisdiction,
		Details:        r.Details,
		CreatedAt:      r.CreatedAt.Time,
		CreatedBy:      r.CreatedBy.String,
		UpdatedAt:      r.UpdatedAt.Time,
		UpdatedBy:      r.UpdatedBy.String,
		Revision:       r.Revision.Int32,
	}
}


// =============================================== Profit & Loss ================================================

func (r *FinancialReportsRepo) GenerateProfitLossReport(ctx context.Context, report domain.ProfitLossReport) (domain.ProfitLossReport, error) {
	dbReport, err := r.q.GenerateProfitLossReport(ctx, db.GenerateProfitLossReportParams{
		OrganizationID: report.OrganizationID,
		PeriodStart:    report.PeriodStart,
		PeriodEnd:      report.PeriodEnd,
		TotalRevenue:   fmt.Sprintf("%f", report.TotalRevenue),  // convert float64 → string
		TotalExpenses:  fmt.Sprintf("%f", report.TotalExpenses),
		NetProfit:      fmt.Sprintf("%f", report.NetProfit),
	})
	if err != nil {
		return domain.ProfitLossReport{}, err
	}
	return mapProfitLossToDomain(dbReport), nil
}


func (r *FinancialReportsRepo) GetProfitLossReport(ctx context.Context, id uuid.UUID) (domain.ProfitLossReport, error) {
	dbReport, err := r.q.GetProfitLossReport(ctx, id)
	if err != nil {
		return domain.ProfitLossReport{}, err
	}
	return mapProfitLossToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) ListProfitLossReports(ctx context.Context, orgID string, limit, offset int32) ([]domain.ProfitLossReport, error) {
	dbReports, err := r.q.ListProfitLossReports(ctx, db.ListProfitLossReportsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	var reports []domain.ProfitLossReport
	for _, r := range dbReports {
		reports = append(reports, mapProfitLossToDomain(r))
	}
	return reports, nil
}


// ================================================== Balance Sheet ================================================

func (r *FinancialReportsRepo) GenerateBalanceSheetReport(ctx context.Context, report domain.BalanceSheetReport) (domain.BalanceSheetReport, error) {
	dbReport, err := r.q.GenerateBalanceSheetReport(ctx, db.GenerateBalanceSheetReportParams{
		OrganizationID:   report.OrganizationID,
		PeriodStart:      report.PeriodStart,
		PeriodEnd:        report.PeriodEnd,
		TotalAssets:      fmt.Sprintf("%f", report.TotalAssets),       // float64 → string
		TotalLiabilities: fmt.Sprintf("%f", report.TotalLiabilities),  // float64 → string
		NetWorth:         fmt.Sprintf("%f", report.NetWorth),          // float64 → string
	})
	if err != nil {
		return domain.BalanceSheetReport{}, err
	}
	return mapBalanceSheetToDomain(dbReport), nil
}


func (r *FinancialReportsRepo) GetBalanceSheetReport(ctx context.Context, id uuid.UUID) (domain.BalanceSheetReport, error) {
	dbReport, err := r.q.GetBalanceSheetReport(ctx, id)
	if err != nil {
		return domain.BalanceSheetReport{}, err
	}
	return mapBalanceSheetToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) ListBalanceSheetReports(ctx context.Context, orgID string, limit, offset int32) ([]domain.BalanceSheetReport, error) {
	dbReports, err := r.q.ListBalanceSheetReports(ctx, db.ListBalanceSheetReportsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	var reports []domain.BalanceSheetReport
	for _, r := range dbReports {
		reports = append(reports, mapBalanceSheetToDomain(r))
	}
	return reports, nil
}


// ================================================== Trial Balance ===================================================

func (r *FinancialReportsRepo) CreateTrialBalanceReport(ctx context.Context, report domain.TrialBalanceReport) (domain.TrialBalanceReport, error) {
	dbReport, err := r.q.CreateTrialBalanceReport(ctx, db.CreateTrialBalanceReportParams{
		OrganizationID: report.OrganizationID,
		PeriodStart:    report.PeriodStart,
		PeriodEnd:      report.PeriodEnd,
	})
	if err != nil {
		return domain.TrialBalanceReport{}, err
	}
	return mapTrialBalanceToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) AddTrialBalanceEntry(ctx context.Context, entry domain.TrialBalanceEntry) (domain.TrialBalanceEntry, error) {
	dbEntry, err := r.q.AddTrialBalanceEntry(ctx, db.AddTrialBalanceEntryParams{
		ReportID:      entry.ReportID,
		LedgerAccount: entry.LedgerAccount,
		Debit:         fmt.Sprintf("%f", entry.Debit),  // float64 → string
		Credit:        fmt.Sprintf("%f", entry.Credit), // float64 → string
	})
	if err != nil {
		return domain.TrialBalanceEntry{}, err
	}
	return mapTrialBalanceEntryToDomain(dbEntry), nil
}


func (r *FinancialReportsRepo) GetTrialBalanceReport(ctx context.Context, id uuid.UUID) (domain.TrialBalanceReport, error) {
	dbReport, err := r.q.GetTrialBalanceReport(ctx, id)
	if err != nil {
		return domain.TrialBalanceReport{}, err
	}
	return mapTrialBalanceToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) ListTrialBalanceReports(ctx context.Context, orgID string, limit, offset int32) ([]domain.TrialBalanceReport, error) {
	dbReports, err := r.q.ListTrialBalanceReports(ctx, db.ListTrialBalanceReportsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	var reports []domain.TrialBalanceReport
	for _, r := range dbReports {
		reports = append(reports, mapTrialBalanceToDomain(r))
	}
	return reports, nil
}

func (r *FinancialReportsRepo) ListTrialBalanceEntries(ctx context.Context, reportID uuid.UUID) ([]domain.TrialBalanceEntry, error) {
	dbEntries, err := r.q.ListTrialBalanceEntries(ctx, reportID)
	if err != nil {
		return nil, err
	}
	var entries []domain.TrialBalanceEntry
	for _, e := range dbEntries {
		entries = append(entries, mapTrialBalanceEntryToDomain(e))
	}
	return entries, nil
}


// ====================================================== Compliance ======================================================

func (r *FinancialReportsRepo) GenerateComplianceReport(ctx context.Context, report domain.ComplianceReport) (domain.ComplianceReport, error) {
	dbReport, err := r.q.GenerateComplianceReport(ctx, db.GenerateComplianceReportParams{
		OrganizationID: report.OrganizationID,
		PeriodStart:    report.PeriodStart,
		PeriodEnd:      report.PeriodEnd,
		Jurisdiction:   report.Jurisdiction,
		Details:        report.Details,
	})
	if err != nil {
		return domain.ComplianceReport{}, err
	}
	return mapComplianceToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) GetComplianceReport(ctx context.Context, id uuid.UUID) (domain.ComplianceReport, error) {
	dbReport, err := r.q.GetComplianceReport(ctx, id)
	if err != nil {
		return domain.ComplianceReport{}, err
	}
	return mapComplianceToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) ListComplianceReports(ctx context.Context, orgID, jurisdiction string, limit, offset int32) ([]domain.ComplianceReport, error) {
	dbReports, err := r.q.ListComplianceReports(ctx, db.ListComplianceReportsParams{
		OrganizationID: orgID,
		Jurisdiction:   jurisdiction,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	var reports []domain.ComplianceReport
	for _, r := range dbReports {
		reports = append(reports, mapComplianceToDomain(r))
	}
	return reports, nil
}
