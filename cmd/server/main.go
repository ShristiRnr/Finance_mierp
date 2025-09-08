package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/kafka"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports/grpc_server"
)

func main() {
	// DB
	conn, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	queries := db.New(conn) // adjust if sqlc generated different ctor

	// Kafka publisher
	kpub := kafka.NewKafkaPublisher([]string{"localhost:9092"})

	// Repositories
	accRepo := repository.NewAccountRepository(queries, kpub)
	jrRepo := repository.NewJournalRepository(conn, queries, kpub)
	ldRepo := repository.NewLedgerRepository(queries)

	// Services
	accSvc := services.NewAccountService(accRepo, kpub)
	journalSvc := services.NewJournalService(jrRepo, kpub)
	ledgerSvc := services.NewLedgerService(ldRepo) // implement if needed

	// gRPC handler wiring
	handler := grpc_server.NewLedgerHandler(accSvc, journalSvc, ledgerSvc, kpub)

	// start gRPC server (your existing code)
	_ = handler
}