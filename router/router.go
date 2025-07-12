package router

import (
	"bank-api/handlers"
	"net/http"

	_ "bank-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRoutes() {
	http.HandleFunc("/accounts", handlers.CreateAccountHandler)
	http.HandleFunc("/accounts/", handlers.GetAccountHandler)

	http.HandleFunc("/transactions", handlers.CreateTransactionHandler)

	http.Handle("/swagger/", httpSwagger.WrapHandler)
}
