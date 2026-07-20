package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/database"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/event"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/event/handler"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/create_account"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/create_client"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/create_transaction"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_account"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_all_accounts"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_all_clients"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_all_transactions"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_client"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_transaction"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/get_account"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/get_client"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/list_accounts"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/list_clients"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/set_account_balance"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/web"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/web/webserver"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/events"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/kafka"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "wallet")
	kafkaServers := getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
	httpPort := getEnv("HTTP_PORT", "8080")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			break
		} else if i == 29 {
			panic(fmt.Sprintf("failed to connect to mysql: %v", err))
		}

		fmt.Println("waiting for mysql...")
		time.Sleep(2 * time.Second)
	}

	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": kafkaServers,
		"acks":              "all",
	}

	kafkaProducer := kafka.NewKafkaProducer(configMap)
	defer kafkaProducer.Close()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewBalanceUpdatedKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDB := database.NewClientDB(db)
	accountDB := database.NewAccountDB(db)
	transactionDB := database.NewTransactionDB(db)

	clientUseCase := create_client.NewCreateClientUseCase(clientDB)
	getClientUseCase := get_client.NewGetClientUseCase(clientDB)
	listClientsUseCase := list_clients.NewListClientsUseCase(clientDB)
	deleteClientUseCase := delete_client.NewDeleteClientUseCase(clientDB, accountDB)
	deleteAllClientsUseCase := delete_all_clients.NewDeleteAllClientsUseCase(clientDB, accountDB, transactionDB)

	accountUseCase := create_account.NewCreateAccountUseCase(accountDB, clientDB)
	setAccountBalanceUseCase := set_account_balance.NewSetAccountBalanceUseCase(accountDB)
	getAccountUseCase := get_account.NewGetAccountUseCase(accountDB)
	listAccountsUseCase := list_accounts.NewListAccountsUseCase(accountDB)
	deleteAccountUseCase := delete_account.NewDeleteAccountUseCase(accountDB, transactionDB)
	deleteAllAccountsUseCase := delete_all_accounts.NewDeleteAllAccountsUseCase(accountDB, transactionDB)

	transactionUnitOfWork := uow.NewSQLUnitOfWork(db)
	transactionUseCase := create_transaction.NewCreateTransactionUseCase(
		transactionUnitOfWork,
		eventDispatcher,
		transactionCreatedEvent,
		balanceUpdatedEvent,
	)
	deleteTransactionUseCase := delete_transaction.NewDeleteTransactionUseCase(transactionDB)
	deleteAllTransactionsUseCase := delete_all_transactions.NewDeleteAllTransactionsUseCase(transactionDB)

	clientHandler := web.NewWebClientHandler(*clientUseCase, *getClientUseCase, *listClientsUseCase, *deleteClientUseCase, *deleteAllClientsUseCase)
	accountHandler := web.NewWebAccountHandler(*accountUseCase, *setAccountBalanceUseCase, *getAccountUseCase, *listAccountsUseCase, *deleteAccountUseCase, *deleteAllAccountsUseCase)
	transactionHandler := web.NewWebTransactionHandler(*transactionUseCase, *deleteTransactionUseCase, *deleteAllTransactionsUseCase)

	webServer := webserver.NewWebServer(":" + httpPort)
	webServer.Get("/clients", clientHandler.ListClients)
	webServer.Get("/clients/{id}", clientHandler.GetClient)
	webServer.Post("/clients", clientHandler.CreateClient)
	webServer.Delete("/clients", clientHandler.DeleteAllClients)
	webServer.Delete("/clients/{id}", clientHandler.DeleteClient)
	webServer.Get("/accounts", accountHandler.ListAccounts)
	webServer.Get("/accounts/{id}", accountHandler.GetAccount)
	webServer.Post("/accounts", accountHandler.CreateAccount)
	webServer.Post("/accounts/balance", accountHandler.SetAccountBalance)
	webServer.Delete("/accounts", accountHandler.DeleteAllAccounts)
	webServer.Delete("/accounts/{id}", accountHandler.DeleteAccount)
	webServer.Post("/transactions", transactionHandler.CreateTransaction)
	webServer.Delete("/transactions", transactionHandler.DeleteAllTransactions)
	webServer.Delete("/transactions/{id}", transactionHandler.DeleteTransaction)
	webServer.Start()
}
