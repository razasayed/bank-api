package router

import (
	"bank-api/handlers"
	"net/http"
)

func InitRoutes() {
	http.HandleFunc("/accounts", handlers.CreateAccountHandler)
	http.HandleFunc("/accounts/", handlers.GetAccountHandler)
}
