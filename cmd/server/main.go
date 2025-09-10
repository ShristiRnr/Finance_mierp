package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/kafka"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports/grpc_server"

	"google.golang.org/grpc"
)

func main() {
	// ---------------- DB Connection ----------------
	conn, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/dbname?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close()

	queries := db.New(conn) // sqlc generated Queries struct

	// ---------------- Kafka Publisher ----------------
	kpub := kafka.NewKafkaPublisher([]string{"localhost:9092"})
	defer kpub.Close()

	// ---------------- Repositories ----------------
	accRepo := repository.NewAccountRepository(queries, kpub)
	jrRepo := repository.NewJournalRepository(conn, queries, kpub)
	ldRepo := repository.NewLedgerRepository(queries)
	accrualRepo := repository.NewAccrualRepository(queries, kpub)
	allocationRepo :=repository.NewAllocationRuleRepository(queries, kpub)
	auditRepo := repository.NewAuditRepository(conn, kpub)
	budgetRepo := repository.NewBudgetRepository(conn)
	CashFlowRepo := repository.NewCashFlowForecastRepo(conn)
	ConsolidationRepo :=repository.NewConsolidationRepo(queries)
	CreditDebitNoteRepo := repository.NewCreditDebitNoteRepo(queries)
	ExchangeRateRepo := repository.NewExchangeRateRepo(queries)

	// ---------------- Services ----------------
	accSvc := services.NewAccountService(accRepo, kpub)
	journalSvc := services.NewJournalService(jrRepo, kpub)
	ledgerSvc := services.NewLedgerService(ldRepo)
	accrualSvc := services.NewAccrualService(accrualRepo, kpub, "accruals")
	allocationSvc := services.NewAllocationService(allocationRepo, kpub, "allocation")
	auditSvc := services.NewAuditService(auditRepo, kpub)
	budgetSvc := services.NewBudgetService(budgetRepo, kpub)
	CashFlowSvc := services.NewCashFlowService(CashFlowRepo, kpub)
	ConsolidationSvc := services.NewConsolidationService(ConsolidationRepo, kpub)
	CreditDebitNoteSvc := services.NewCreditDebitNoteService(CreditDebitNoteRepo, kpub)
	ExchangeRateSvc := services.NewExchangeRateService(ExchangeRateRepo, kpub)

	// ---------------- gRPC Handlers ----------------
	ledgerHandler := grpc_server.NewLedgerHandler(accSvc, journalSvc, ledgerSvc)
	accrualHandler := grpc_server.NewAccrualHandler(accrualSvc)
	allocationHandler := grpc_server.NewAllocationHandler(allocationSvc, kpub)
	auditHandler := grpc_server.NewAuditHandler(auditSvc, kpub)
	budgetHandler := grpc_server.NewBudgetHandler(budgetSvc)
	CashFlowHandler := grpc_server.NewCashFlowGRPCServer(CashFlowSvc)
	ConsolidationHandler := grpc_server.NewConsolidationHandler(ConsolidationSvc)
	CreditDebitNoteHandler := grpc_server.NewGRPCServer(CreditDebitNoteSvc)
	ExchangeRateHandler := grpc_server.NewExchangeRateHandler(ExchangeRateSvc)

	// ---------------- gRPC Server ----------------
	grpcServer := grpc.NewServer()

	// Register services
	pb.RegisterLedgerServiceServer(grpcServer, ledgerHandler)
	pb.RegisterAccrualServiceServer(grpcServer, accrualHandler)
	pb.RegisterAllocationAutomationServiceServer(grpcServer, allocationHandler)
	pb.RegisterAuditTrailServiceServer(grpcServer, auditHandler)
	pb.RegisterBudgetServiceServer(grpcServer, budgetHandler)
	pb.RegisterCashFlowServiceServer(grpcServer, CashFlowHandler)
	pb.RegisterConsolidationServiceServer(grpcServer, ConsolidationHandler)
	pb.RegisterCreditDebitNoteServiceServer(grpcServer, CreditDebitNoteHandler)
	pb.RegisterFxServiceServer(grpcServer, ExchangeRateHandler)

	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("🚀 gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
