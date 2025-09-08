package repository

import (
	"context"
	"strconv"
	"fmt"

	"github.com/google/uuid"
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

func mapProfitLossToDomain(r db.ProfitLossReport) db.ProfitLossReport {
	// Parse TotalRevenue
	totalRevenue, err := strconv.ParseFloat(r.TotalRevenue, 64)
	if err != nil {
		totalRevenue = 0
	}

	// Parse TotalExpenses
	totalExpenses, err := strconv.ParseFloat(r.TotalExpenses, 64)
	if err != nil {
		totalExpenses = 0
	}

	// Calculate NetProfit
	netProfit := totalRevenue - totalExpenses

	return db.ProfitLossReport{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		PeriodStart:    r.PeriodStart,
		PeriodEnd:      r.PeriodEnd,
		TotalRevenue:   fmt.Sprintf("%.2f", totalRevenue), // store as string
		TotalExpenses:  fmt.Sprintf("%.2f", totalExpenses),
		NetProfit:      fmt.Sprintf("%.2f", netProfit),
		CreatedAt:      r.CreatedAt,
		CreatedBy:      r.CreatedBy,
		UpdatedAt:      r.UpdatedAt,
		UpdatedBy:      r.UpdatedBy,
		Revision:       r.Revision,
	}
}

func mapBalanceSheetToDomain(r db.BalanceSheetReport) db.BalanceSheetReport {

	return db.BalanceSheetReport{
		ID:              r.ID,
		OrganizationID:  r.OrganizationID,
		PeriodStart:     r.PeriodStart,
		PeriodEnd:       r.PeriodEnd,
		TotalAssets:     r.TotalAssets,
		TotalLiabilities: r.TotalLiabilities,
		NetWorth:        r.NetWorth,
		CreatedAt:       r.CreatedAt,
		CreatedBy:       r.CreatedBy,
		UpdatedAt:       r.UpdatedAt,
		UpdatedBy:       r.UpdatedBy,
		Revision:        r.Revision,
	}
}


func mapTrialBalanceToDomain(r db.TrialBalanceReport) db.TrialBalanceReport {
	return db.TrialBalanceReport{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		PeriodStart:    r.PeriodStart,
		PeriodEnd:      r.PeriodEnd,
		CreatedAt:      r.CreatedAt,
		CreatedBy:      r.CreatedBy,
		UpdatedAt:      r.UpdatedAt,
		UpdatedBy:      r.UpdatedBy,
		Revision:       r.Revision,
	}
}

func mapTrialBalanceEntryToDomain(r db.TrialBalanceEntry) db.TrialBalanceEntry {

	return db.TrialBalanceEntry{
		ID:            r.ID,
		ReportID:      r.ReportID,
		LedgerAccount: r.LedgerAccount,
		Debit:         r.Debit,
		Credit:        r.Credit,
		CreatedAt:     r.CreatedAt,
		CreatedBy:     r.CreatedBy,
	}
}


func mapComplianceToDomain(r db.ComplianceReport) db.ComplianceReport {
	return db.ComplianceReport{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		PeriodStart:    r.PeriodStart,
		PeriodEnd:      r.PeriodEnd,
		Jurisdiction:   r.Jurisdiction,
		Details:        r.Details,
		CreatedAt:      r.CreatedAt,
		CreatedBy:      r.CreatedBy,
		UpdatedAt:      r.UpdatedAt,
		UpdatedBy:      r.UpdatedBy,
		Revision:       r.Revision,
	}
}


// =============================================== Profit & Loss ================================================

func (r *FinancialReportsRepo) GenerateProfitLossReport(ctx context.Context, report db.ProfitLossReport) (db.ProfitLossReport, error) {
	dbReport, err := r.q.GenerateProfitLossReport(ctx, db.GenerateProfitLossReportParams{
		OrganizationID: report.OrganizationID,
		PeriodStart:    report.PeriodStart,
		PeriodEnd:      report.PeriodEnd,
		TotalRevenue:   report.TotalRevenue,  // convert float64 → string
		TotalExpenses:  report.TotalExpenses,
		NetProfit:      report.NetProfit,
	})
	if err != nil {
		return db.ProfitLossReport{}, err
	}
	return mapProfitLossToDomain(dbReport), nil
}


func (r *FinancialReportsRepo) GetProfitLossReport(ctx context.Context, id uuid.UUID) (db.ProfitLossReport, error) {
	dbReport, err := r.q.GetProfitLossReport(ctx, id)
	if err != nil {
		return db.ProfitLossReport{}, err
	}
	return mapProfitLossToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) ListProfitLossReports(ctx context.Context, orgID string, limit, offset int32) ([]db.ProfitLossReport, error) {
	dbReports, err := r.q.ListProfitLossReports(ctx, db.ListProfitLossReportsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	var reports []db.ProfitLossReport
	for _, r := range dbReports {
		reports = append(reports, mapProfitLossToDomain(r))
	}
	return reports, nil
}


// ================================================== Balance Sheet ================================================

func (r *FinancialReportsRepo) GenerateBalanceSheetReport(ctx context.Context, report db.BalanceSheetReport) (db.BalanceSheetReport, error) {
	dbReport, err := r.q.GenerateBalanceSheetReport(ctx, db.GenerateBalanceSheetReportParams{
		OrganizationID:   report.OrganizationID,
		PeriodStart:      report.PeriodStart,
		PeriodEnd:        report.PeriodEnd,
		TotalAssets:      report.TotalAssets,       // float64 → string
		TotalLiabilities: report.TotalLiabilities,  // float64 → string
		NetWorth:         report.NetWorth,          // float64 → string
	})
	if err != nil {
		return db.BalanceSheetReport{}, err
	}
	return mapBalanceSheetToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) ListBalanceSheetReports(ctx context.Context, orgID string, limit, offset int32) ([]db.BalanceSheetReport, error) {
	dbReports, err := r.q.ListBalanceSheetReports(ctx, db.ListBalanceSheetReportsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	var reports []db.BalanceSheetReport
	reports = append(reports, dbReports...)

	return reports, nil
}

func (r *FinancialReportsRepo) GetBalanceSheetReport(ctx context.Context, id uuid.UUID) (db.BalanceSheetReport, error) {
    // delegate to SQLC / DB query
    report, err := r.q.GetBalanceSheetReport(ctx, id)
    if err != nil {
        return db.BalanceSheetReport{}, err
    }
    return report, nil
}


// ================================================== Trial Balance ===================================================

func (r *FinancialReportsRepo) CreateTrialBalanceReport(ctx context.Context, report db.TrialBalanceReport) (db.TrialBalanceReport, error) {
	dbReport, err := r.q.CreateTrialBalanceReport(ctx, db.CreateTrialBalanceReportParams{
		OrganizationID: report.OrganizationID,
		PeriodStart:    report.PeriodStart,
		PeriodEnd:      report.PeriodEnd,
	})
	if err != nil {
		return db.TrialBalanceReport{}, err
	}
	return mapTrialBalanceToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) AddTrialBalanceEntry(ctx context.Context, entry db.TrialBalanceEntry) (db.TrialBalanceEntry, error) {
	dbEntry, err := r.q.AddTrialBalanceEntry(ctx, db.AddTrialBalanceEntryParams{
		ReportID:      entry.ReportID,
		LedgerAccount: entry.LedgerAccount,
		Debit:         entry.Debit,  // float64 → string
		Credit:        entry.Credit, // float64 → string
	})
	if err != nil {
		return db.TrialBalanceEntry{}, err
	}
	return mapTrialBalanceEntryToDomain(dbEntry), nil
}


func (r *FinancialReportsRepo) GetTrialBalanceReport(ctx context.Context, id uuid.UUID) (db.TrialBalanceReport, error) {
	dbReport, err := r.q.GetTrialBalanceReport(ctx, id)
	if err != nil {
		return db.TrialBalanceReport{}, err
	}
	return mapTrialBalanceToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) ListTrialBalanceReports(ctx context.Context, orgID string, limit, offset int32) ([]db.TrialBalanceReport, error) {
	dbReports, err := r.q.ListTrialBalanceReports(ctx, db.ListTrialBalanceReportsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	var reports []db.TrialBalanceReport
	reports = append(reports, dbReports...)
	return reports, nil
}

func (r *FinancialReportsRepo) ListTrialBalanceEntries(ctx context.Context, reportID uuid.UUID) ([]db.TrialBalanceEntry, error) {
	dbEntries, err := r.q.ListTrialBalanceEntries(ctx, reportID)
	if err != nil {
		return nil, err
	}
	var entries []db.TrialBalanceEntry
	entries = append(entries, dbEntries...)
	return entries, nil
}


// ====================================================== Compliance ======================================================

func (r *FinancialReportsRepo) GenerateComplianceReport(ctx context.Context, report db.ComplianceReport) (db.ComplianceReport, error) {
	dbReport, err := r.q.GenerateComplianceReport(ctx, db.GenerateComplianceReportParams{
		OrganizationID: report.OrganizationID,
		PeriodStart:    report.PeriodStart,
		PeriodEnd:      report.PeriodEnd,
		Jurisdiction:   report.Jurisdiction,
		Details:        report.Details,
	})
	if err != nil {
		return db.ComplianceReport{}, err
	}
	return mapComplianceToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) GetComplianceReport(ctx context.Context, id uuid.UUID) (db.ComplianceReport, error) {
	dbReport, err := r.q.GetComplianceReport(ctx, id)
	if err != nil {
		return db.ComplianceReport{}, err
	}
	return mapComplianceToDomain(dbReport), nil
}

func (r *FinancialReportsRepo) ListComplianceReports(ctx context.Context, orgID, jurisdiction string, limit, offset int32) ([]db.ComplianceReport, error) {
	dbReports, err := r.q.ListComplianceReports(ctx, db.ListComplianceReportsParams{
		OrganizationID: orgID,
		Jurisdiction:   jurisdiction,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	var reports []db.ComplianceReport
	reports = append(reports, dbReports...)
	return reports, nil
}
