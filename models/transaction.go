package models

import (
	"bank-api/db"
	"bank-api/utils"
	"errors"
	"strings"
	"time"
)

type Transaction struct {
	TransactionID   int       `json:"transaction_id" example:"101"`
	AccountID       int       `json:"account_id" example:"1"`
	OperationTypeID int       `json:"operation_type_id" example:"4"`
	Amount          float64   `json:"amount" example:"100.00"`
	EventDate       time.Time `json:"event_date" example:"2025-10-01T12:00:00Z"`
}

func CreateTransaction(accountID, operationTypeID int, amount float64) (*Transaction, error) {
	// Enforce business rule: only payments (operation type 4) can be positive
	if operationTypeID != utils.OperationTypePayment {
		amount = -amount
	}

	var transaction Transaction
	err := db.DB.QueryRow(`
		INSERT INTO transactions (account_id, operation_type_id, amount)
		VALUES ($1, $2, $3)
		RETURNING transaction_id, account_id, operation_type_id, amount, event_date
	`, accountID, operationTypeID, amount).
		Scan(&transaction.TransactionID, &transaction.AccountID, &transaction.OperationTypeID, &transaction.Amount, &transaction.EventDate)

	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return nil, errors.New(utils.ErrAccountNotFound)
		}
		return nil, err
	}

	return &transaction, nil
}
