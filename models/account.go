package models

import (
	"bank-api/db"
	"database/sql"
	"errors"
)

var ErrAccountNotFound = errors.New("account not found")

type Account struct {
	AccountID      int    `json:"account_id"  example:"1"`
	DocumentNumber string `json:"document_number" example:"12345678900"`
}

func GetAccountByID(id int) (*Account, error) {
	var acc Account
	err := db.DB.QueryRow(`
	  SELECT account_id, document_number
	  FROM accounts
	  WHERE account_id = $1
	`, id).Scan(&acc.AccountID, &acc.DocumentNumber)

	switch {
	case err == sql.ErrNoRows:
		return nil, ErrAccountNotFound
	case err != nil:
		return nil, err
	default:
		return &acc, nil
	}
}

func CreateAccount(documentNumber string) (*Account, error) {
	var id int
	err := db.DB.QueryRow(`
		INSERT INTO accounts (document_number)
		VALUES ($1)
		RETURNING account_id
	`, documentNumber).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &Account{
		AccountID:      id,
		DocumentNumber: documentNumber,
	}, nil
}
