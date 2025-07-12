package models

type CreateAccountInput struct {
	DocumentNumber string `json:"document_number" example:"12345678900"`
}

type CreateTransactionInput struct {
	AccountID       int     `json:"account_id" example:"1"`
	OperationTypeID int     `json:"operation_type_id" example:"4"`
	Amount          float64 `json:"amount" example:"123.45"`
}
