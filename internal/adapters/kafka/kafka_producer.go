package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/segmentio/kafka-go"
)

// ---------------- Topics ----------------
const (
	TopicAccounts            = "accounts"
	TopicJournals            = "journals"
	TopicAccruals            = "accruals"
	TopicAllocationRules     = "allocation_rules"
	TopicAuditEvents         = "audit_events"
	TopicBudgets             = "budgets"
	TopicBudgetAlloc         = "budget_allocations"
	TopicCashFlowForecasts   = "cash_flow_forecasts"
	TopicBankTransactions    = "bank_transactions"
	TopicInvoices            = "invoices"
	TopicExpenses            = "expenses"
	TopicCostCenters         = "cost_centers"
	TopicFinanceEvents       = "finance_events"
	TopicGST                 = "gst_events"
	TopicExchangeRates       = "exchange_rates"
	TopicCostAllocations     = "cost_allocations"
	TopicConsolidationEvents = "consolidations"
	TopicCreditDebitEvents   = "credit_debit"
	TopicFinanceReports      = "finance_reports"
)

// ---------------- Event Names ----------------
const (
	// Account
	EventAccountCreated = "account.created"
	EventAccountUpdated = "account.updated"
	EventAccountDeleted = "account.deleted"

	// Journal
	EventJournalCreated = "journal.created"
	EventJournalUpdated = "journal.updated"
	EventJournalDeleted = "journal.deleted"

	// Accrual
	EventAccrualCreated = "accrual.created"
	EventAccrualUpdated = "accrual.updated"
	EventAccrualDeleted = "accrual.deleted"

	// AllocationRule
	EventAllocationRuleCreated = "allocation_rule.created"
	EventAllocationRuleUpdated = "allocation_rule.updated"
	EventAllocationRuleDeleted = "allocation_rule.deleted"

	// Audit
	EventAuditRecorded = "audit.recorded"

	// Budget
	EventBudgetCreated           = "budget.created"
	EventBudgetUpdated           = "budget.updated"
	EventBudgetDeleted           = "budget.deleted"
	EventBudgetAllocated         = "budget.allocated"
	EventBudgetAllocationUpdated = "budget.allocation.updated"
	EventBudgetAllocationDeleted = "budget.allocation.deleted"

	// CashFlowForecast
	EventCashFlowForecastGenerated = "forecast.generated"
	EventCashFlowForecastFetched   = "forecast.fetched"
	EventCashFlowForecastListed    = "forecast.listed"

	// Bank Transactions
	EventBankTransactionImported   = "bank.transaction.imported"
	EventBankTransactionReconciled = "bank.transaction.reconciled"

	// Invoices
	RecordInvoiceCreated       = "invoice.created"
	EventInvoiceUpdated        = "invoice.updated"
	EventInvoiceDeleted        = "invoice.deleted"
	EventInvoiceItemCreated    = "invoice.item.created"
	EventInvoiceTaxAdded       = "invoice.tax.added"
	EventInvoiceDiscountAdded  = "invoice.discount.added"
	EventFinanceInvoiceCreated = "finance.invoice.created"

	// Expenses
	EventExpenseCreated = "expense.created"
	EventExpenseUpdated = "expense.updated"
	EventExpenseDeleted = "expense.deleted"

	// Cost Centers
	EventCostCenterCreated = "costcenter.created"
	EventCostCenterUpdated = "costcenter.updated"
	EventCostCenterDeleted = "costcenter.deleted"

	// GST
	EventGstBreakupAdded   = "gst.breakup.added"
	EventGstRegimeAdded    = "gst.regime.added"
	EventGstDocStatusAdded = "gst.docstatus.added"

	// Exchange rates
	EventExchangeRateCreated = "exchange_rate.created"
	EventExchangeRateUpdated = "exchange_rate.updated"
	EventExchangeRateDeleted = "exchange_rate.deleted"

	// CostAllocation
	EventCostAllocationAllocated = "cost_allocation.allocated"
	EventCostAllocationListed    = "cost_allocation.listed"

	//Consolidation
	EventConsolidationCreated = "consolidation.created"
	EventConsolidationDeleted = "consolidation.deleted"

	//Credit/Debit Note
	EventCreditDebitNoteCreated = "credit_debit_note.created"
	EventCreditDebitNoteUpdated = "credit_debit_note.updated"
	EventCreditDebitNoteDeleted = "credit_debit_note.deleted"

	//Finanace Event
	EventFinancePaymentReceived = "payment.received"
	EventInventoryCostPosted    = "inventiry.cost.posted"
	EventPayrollPosted          = "payroll.posted"
	EventVendorBillApproved     = "vendor.bill.approved"

	//Finance Reports
	EventProfitLossGenerated             = "profit_loss.generated"
	EventBalanceSheetGenerated           = "balance_sheet.generated"
	EventTrialBalanceCreatedGenerated    = "trail_balance.created"
	EventTrailBalanceEntryAddedGenerated = "trail_balance.entry_added"
	EventComplianceGenerated             = "compliance.generated"

	EventTrialBalanceGenerated  = "trial_balance.generated"
	EventTrialBalanceEntryAdded = "trial_balance.entry_added"
)

// ---------------- Kafka Publisher ----------------
type KafkaPublisher struct {
	writers map[string]*kafka.Writer
}

// create a new Kafka writer for a topic
func newWriter(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
}

func NewKafkaPublisher(brokers []string) ports.EventPublisher {
	topics := []string{
		TopicAccounts,
		TopicJournals,
		TopicAccruals,
		TopicAllocationRules,
		TopicAuditEvents,
		TopicBudgets,
		TopicBudgetAlloc,
		TopicCashFlowForecasts,
		TopicBankTransactions,
		TopicInvoices,
		TopicExpenses,
		TopicCostCenters,
		TopicFinanceEvents,
		TopicGST,
		TopicExchangeRates,
		TopicCostAllocations,
		TopicConsolidationEvents,
		TopicCreditDebitEvents,
		TopicFinanceReports,
	}

	writers := make(map[string]*kafka.Writer)
	for _, t := range topics {
		writers[t] = newWriter(brokers, t)
	}
	return &KafkaPublisher{writers: writers}
}

// generic publish function
func (p *KafkaPublisher) publish(ctx context.Context, topic, key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return p.Publish(ctx, topic, key, data)
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic, key string, payload []byte) error {
	writer, ok := p.writers[topic]
	if !ok {
		return errors.New("kafka writer not found for topic: " + topic)
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: payload,
		Time:  time.Now(),
	}

	for i := 0; i < 3; i++ {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := writer.WriteMessages(ctxTimeout, msg)
		cancel()
		if err == nil {
			return nil
		}
		log.Printf("kafka publish error attempt=%d topic=%s: %v", i+1, topic, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return errors.New("failed to publish message after retries")
}

// Close all writers gracefully
func (p *KafkaPublisher) Close() error {
	for topic, w := range p.writers {
		if err := w.Close(); err != nil {
			log.Printf("failed to close writer for topic=%s: %v", topic, err)
			return err
		}
	}
	return nil
}

// ---------------- Kafka Consumer ----------------
type KafkaConsumer struct {
	readers map[string]*kafka.Reader
}

// create a new reader for a topic
func newReader(brokers []string, groupID, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
	})
}

// initialize consumer with all topics
func NewKafkaConsumer(brokers []string, groupID string) *KafkaConsumer {
	topics := []string{
		TopicAccounts,
		TopicJournals,
		TopicAccruals,
		TopicAllocationRules,
		TopicAuditEvents,
		TopicBudgets,
		TopicBudgetAlloc,
		TopicCashFlowForecasts,
		TopicBankTransactions,
		TopicInvoices,
		TopicExpenses,
		TopicCostCenters,
		TopicFinanceEvents,
		TopicGST,
		TopicExchangeRates,
		TopicCostAllocations,
		TopicConsolidationEvents,
		TopicCreditDebitEvents,
		TopicFinanceReports,
	}

	readers := make(map[string]*kafka.Reader)
	for _, t := range topics {
		readers[t] = newReader(brokers, groupID, t)
	}
	return &KafkaConsumer{readers: readers}
}

// consume messages for a topic
func (c *KafkaConsumer) Consume(ctx context.Context, topic string, handler func(key, value []byte) error) error {
	reader, ok := c.readers[topic]
	if !ok {
		return errors.New("kafka reader not found for topic: " + topic)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("stopping consumer for topic=%s", topic)
			return ctx.Err()
		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				log.Printf("kafka consume error topic=%s: %v", topic, err)
				continue
			}
			if handler != nil {
				if err := handler(m.Key, m.Value); err != nil {
					log.Printf("handler error: %v", err)
				}
			}
		}
	}
}

// close all readers
func (c *KafkaConsumer) Close() error {
	for topic, r := range c.readers {
		if err := r.Close(); err != nil {
			log.Printf("failed to close reader for topic=%s: %v", topic, err)
			return err
		}
	}
	return nil
}

// -------------------- Typed Kafka Publish Methods --------------------

// Account events
func (p *KafkaPublisher) PublishAccountCreated(ctx context.Context, a *db.Account) error {
	return p.publish(ctx, TopicAccounts, EventAccountCreated, a)
}
func (p *KafkaPublisher) PublishAccountUpdated(ctx context.Context, a *db.Account) error {
	return p.publish(ctx, TopicAccounts, EventAccountUpdated, a)
}
func (p *KafkaPublisher) PublishAccountDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicAccounts, EventAccountDeleted, map[string]string{"id": id})
}

// Journal events
func (p *KafkaPublisher) PublishJournalCreated(ctx context.Context, j *db.JournalEntry) error {
	return p.publish(ctx, TopicJournals, EventJournalCreated, j)
}
func (p *KafkaPublisher) PublishJournalUpdated(ctx context.Context, j *db.JournalEntry) error {
	return p.publish(ctx, TopicJournals, EventJournalUpdated, j)
}
func (p *KafkaPublisher) PublishJournalDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicJournals, EventJournalDeleted, map[string]string{"id": id})
}

// Accrual events
func (p *KafkaPublisher) PublishAccrualCreated(ctx context.Context, a *db.Accrual) error {
	return p.publish(ctx, TopicAccruals, EventAccrualCreated, a)
}
func (p *KafkaPublisher) PublishAccrualUpdated(ctx context.Context, a *db.Accrual) error {
	return p.publish(ctx, TopicAccruals, EventAccrualUpdated, a)
}
func (p *KafkaPublisher) PublishAccrualDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicAccruals, EventAccrualDeleted, map[string]string{"id": id})
}

// AllocationRule events
func (p *KafkaPublisher) PublishAllocationRuleCreated(ctx context.Context, r *db.AllocationRule) error {
	return p.publish(ctx, TopicAllocationRules, EventAllocationRuleCreated, r)
}
func (p *KafkaPublisher) PublishAllocationRuleUpdated(ctx context.Context, r *db.AllocationRule) error {
	return p.publish(ctx, TopicAllocationRules, EventAllocationRuleUpdated, r)
}
func (p *KafkaPublisher) PublishAllocationRuleDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicAllocationRules, EventAllocationRuleDeleted, map[string]string{"id": id})
}

// Audit events
func (p *KafkaPublisher) PublishAuditRecorded(ctx context.Context, a *db.AuditEvent) error {
	return p.publish(ctx, TopicAuditEvents, EventAuditRecorded, a)
}

// Budget events
func (p *KafkaPublisher) PublishBudgetCreated(ctx context.Context, b *db.Budget) error {
	return p.publish(ctx, TopicBudgets, EventBudgetCreated, b)
}
func (p *KafkaPublisher) PublishBudgetUpdated(ctx context.Context, b *db.Budget) error {
	return p.publish(ctx, TopicBudgets, EventBudgetUpdated, b)
}
func (p *KafkaPublisher) PublishBudgetDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicBudgets, EventBudgetDeleted, map[string]string{"id": id})
}

// BudgetAllocation events
func (p *KafkaPublisher) PublishBudgetAllocated(ctx context.Context, ba *db.BudgetAllocation) error {
	return p.publish(ctx, TopicBudgetAlloc, EventBudgetAllocated, ba)
}
func (p *KafkaPublisher) PublishBudgetAllocationUpdated(ctx context.Context, ba *db.BudgetAllocation) error {
	return p.publish(ctx, TopicBudgetAlloc, EventBudgetAllocationUpdated, ba)
}
func (p *KafkaPublisher) PublishBudgetAllocationDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicBudgetAlloc, EventBudgetAllocationDeleted, map[string]string{"id": id})
}

// CashFlowForecast events
// CashFlowForecast events
func (p *KafkaPublisher) PublishCashFlowForecastGenerated(ctx context.Context, ev *db.CashFlowForecast) error {
	return p.publish(ctx, TopicCashFlowForecasts, EventCashFlowForecastGenerated, ev)
}

func (p *KafkaPublisher) PublishCashFlowForecastFetched(ctx context.Context, ev *db.CashFlowForecast) error {
	return p.publish(ctx, TopicCashFlowForecasts, EventCashFlowForecastFetched, ev)
}

func (p *KafkaPublisher) PublishCashFlowForecastListed(ctx context.Context, ev []db.CashFlowForecast) error {
	return p.publish(ctx, TopicCashFlowForecasts, EventCashFlowForecastListed, ev)
}


// Consolidation events
func (p *KafkaPublisher) PublishConsolidationCreated(ctx context.Context, c *db.Consolidation) error {
	return p.publish(ctx, TopicConsolidationEvents, EventConsolidationCreated, c)
}
func (p *KafkaPublisher) PublishConsolidationDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicConsolidationEvents, EventConsolidationDeleted, map[string]string{"id": id})
}

// Credit/Debit Note events
func (p *KafkaPublisher) PublishCreditDebitNoteCreated(ctx context.Context, n *db.CreditDebitNote) error {
	return p.publish(ctx, TopicCreditDebitEvents, EventCreditDebitNoteCreated, n)
}
func (p *KafkaPublisher) PublishCreditDebitNoteUpdated(ctx context.Context, n *db.CreditDebitNote) error {
	return p.publish(ctx, TopicCreditDebitEvents, EventCreditDebitNoteUpdated, n)
}
func (p *KafkaPublisher) PublishCreditDebitNoteDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicCreditDebitEvents, EventCreditDebitNoteDeleted, map[string]string{"id": id})
}

// ExchangeRate events
func (p *KafkaPublisher) PublishExchangeRateCreated(ctx context.Context, r *db.ExchangeRate) error {
	return p.publish(ctx, TopicExchangeRates, EventExchangeRateCreated, r)
}
func (p *KafkaPublisher) PublishExchangeRateUpdated(ctx context.Context, r *db.ExchangeRate) error {
	return p.publish(ctx, TopicExchangeRates, EventExchangeRateUpdated, r)
}
func (p *KafkaPublisher) PublishExchangeRateDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicExchangeRates, EventExchangeRateDeleted, map[string]string{"id": id})
}

// Expense events
func (p *KafkaPublisher) PublishExpenseCreated(ctx context.Context, e *db.Expense) error {
	return p.publish(ctx, TopicExpenses, EventExpenseCreated, e)
}
func (p *KafkaPublisher) PublishExpenseUpdated(ctx context.Context, e *db.Expense) error {
	return p.publish(ctx, TopicExpenses, EventExpenseUpdated, e)
}
func (p *KafkaPublisher) PublishExpenseDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicExpenses, EventExpenseDeleted, map[string]string{"id": id})
}

// CostCenter events
func (p *KafkaPublisher) PublishCostCenterCreated(ctx context.Context, c *db.CostCenter) error {
	return p.publish(ctx, TopicCostCenters, EventCostCenterCreated, c)
}
func (p *KafkaPublisher) PublishCostCenterUpdated(ctx context.Context, c *db.CostCenter) error {
	return p.publish(ctx, TopicCostCenters, EventCostCenterUpdated, c)
}
func (p *KafkaPublisher) PublishCostCenterDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicCostCenters, EventCostCenterDeleted, map[string]string{"id": id})
}

// CostAllocation events
func (p *KafkaPublisher) PublishCostAllocationAllocated(ctx context.Context, ca *db.CostAllocation) error {
	return p.publish(ctx, TopicCostAllocations, EventCostAllocationAllocated, ca)
}
func (p *KafkaPublisher) PublishCostAllocationListed(ctx context.Context, ca []db.CostAllocation) error {
	return p.publish(ctx, TopicCostAllocations, EventCostAllocationListed, ca)
}

// Finance domain events
func (p *KafkaPublisher) PublishFinanceInvoiceCreated(ctx context.Context, ev *db.FinanceInvoiceCreatedEvent) error {
	return p.publish(ctx, TopicFinanceEvents, EventFinanceInvoiceCreated, ev)
}

func (p *KafkaPublisher) PublishFinancePaymentReceived(ctx context.Context, ev *db.FinancePaymentReceivedEvent) error {
	return p.publish(ctx, TopicFinanceEvents, EventFinancePaymentReceived, ev)
}

func (p *KafkaPublisher) PublishInventoryCostPosted(ctx context.Context, ev *db.InventoryCostPostedEvent) error {
	return p.publish(ctx, TopicFinanceEvents, EventInventoryCostPosted, ev)
}

func (p *KafkaPublisher) PublishPayrollPosted(ctx context.Context, ev *db.PayrollPostedEvent) error {
	return p.publish(ctx, TopicFinanceEvents, EventPayrollPosted, ev)
}

func (p *KafkaPublisher) PublishVendorBillApproved(ctx context.Context, ev *db.VendorBillApprovedEvent) error {
	return p.publish(ctx, TopicFinanceEvents, EventVendorBillApproved, ev)
}

// Finance Report
func (p *KafkaPublisher) PublishProfitLossGenerated(ctx context.Context, result *db.ProfitLossReport) error {
	return p.publish(ctx, TopicFinanceReports, EventProfitLossGenerated, result)
}
func (p *KafkaPublisher) PublishBalanceSheetGenerated(ctx context.Context, result *db.BalanceSheetReport) error {
	return p.publish(ctx, TopicFinanceReports, EventBalanceSheetGenerated, result)
}
func (p *KafkaPublisher) PublishTrialBalanceGenerated(ctx context.Context, result *db.TrialBalanceReport) error {
	return p.publish(ctx, TopicFinanceReports, EventTrialBalanceCreatedGenerated, result)
}
func (p *KafkaPublisher) PublishTrailBalanceGeneratedAdded(ctx context.Context, result *db.TrialBalanceEntry) error {
	return p.publish(ctx, TopicFinanceReports, EventTrailBalanceEntryAddedGenerated, result)
}
func (p *KafkaPublisher) PublishComplianceReportGenerated(ctx context.Context, result *db.ComplianceReport) error {
	return p.publish(ctx, TopicFinanceReports, EventComplianceGenerated, result)
}

// GST Events
func (p *KafkaPublisher) PublishGstBreakupAdded(ctx context.Context, inv *db.GstBreakup) error {
	return p.publish(ctx, TopicGST, EventGstBreakupAdded, inv)
}
func (p *KafkaPublisher) PublishGstRegimeAdded(ctx context.Context, regime *db.GstRegime) error {
	return p.publish(ctx, TopicGST, EventGstRegimeAdded, regime)
}
func (p *KafkaPublisher) PublishGstDocStatusAdded(ctx context.Context, status *db.GstDocStatus) error {
	return p.publish(ctx, TopicGST, EventGstDocStatusAdded, status)
}

// Invoice Events
func (p *KafkaPublisher) PublishInvoiceCreated(ctx context.Context, inv *db.Invoice) error {
	return p.publish(ctx, TopicInvoices, RecordInvoiceCreated, inv)
}
func (p *KafkaPublisher) PublishInvoiceUpdated(ctx context.Context, inv *db.Invoice) error {
	return p.publish(ctx, TopicInvoices, EventInvoiceUpdated, inv)
}
func (p *KafkaPublisher) PublishInvoiceDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicInvoices, EventInvoiceDeleted, id)
}
func (p *KafkaPublisher) PublishInvoiceItemCreated(ctx context.Context, item *db.InvoiceItem) error {
	return p.publish(ctx, TopicInvoices, EventInvoiceItemCreated, item)
}
func (p *KafkaPublisher) PublishInvoiceTaxAdded(ctx context.Context, tax *db.InvoiceTax) error {
	return p.publish(ctx, TopicInvoices, EventInvoiceTaxAdded, tax)
}
func (p *KafkaPublisher) PublishInvoiceDiscountAdded(ctx context.Context, disc *db.InvoiceDiscount) error {
	return p.publish(ctx, TopicInvoices, EventInvoiceDiscountAdded, disc)
}

func (p *KafkaPublisher) PublishBankTransactionImported(ctx context.Context, tx *db.BankAccount) error {
	return p.publish(ctx, TopicBankTransactions, EventBankTransactionImported, tx)
}
func (p *KafkaPublisher) PublishBankTransactionReconciled(ctx context.Context, tx *db.BankAccount) error {
	return p.publish(ctx, TopicBankTransactions, EventBankTransactionReconciled, tx)
}

func (p *KafkaPublisher) PublishTrialBalanceEntryAdded(ctx context.Context, entry *db.TrialBalanceEntry) error {
	return p.publish(ctx, TopicFinanceReports, EventTrialBalanceEntryAdded, entry)
}
