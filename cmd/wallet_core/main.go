package main

import (
	"database/sql"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/database"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/event"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/create_account"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/create_client"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/create_transaction"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/web"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/web/webserver"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/events"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/wallet")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDB := database.NewClientDB(db)
	accountDB := database.NewAccountDB(db)
	transactionDB := database.NewTransactionDB(db)

	clientUseCase := create_client.NewCreateClientUseCase(clientDB)
	accountUseCase := create_account.NewCreateAccountUseCase(accountDB, clientDB)
	transactionUseCase := create_transaction.NewCreateTransactionUseCase(transactionDB, accountDB, eventDispatcher, transactionCreatedEvent)

	clientHandler := web.NewWebClientHandler(*clientUseCase)
	accountHandler := web.NewWebAccountHandler(*accountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*transactionUseCase)

	webServer := webserver.NewWebServer(":8080")
	webServer.AddHandler("/clients", clientHandler.CreateClient)
	webServer.AddHandler("/accounts", accountHandler.CreateAccount)
	webServer.AddHandler("/transactions", transactionHandler.CreateTransaction)
	webServer.Start()
}
