package handlers

import (
	"bank-api/models"
	"bank-api/utils"
	"encoding/json"
	"net/http"
)

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, utils.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		AccountID       int     `json:"account_id"`
		OperationTypeID int     `json:"operation_type_id"`
		Amount          float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteJSONError(w, utils.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if input.AccountID == 0 || input.OperationTypeID == 0 || input.Amount == 0 {
		utils.WriteJSONError(w, utils.ErrMissingAccountFields, http.StatusBadRequest)
		return
	}

	if !models.OperationTypeExists(input.OperationTypeID) {
		utils.WriteJSONError(w, utils.ErrInvalidOperationTypeID, http.StatusBadRequest)
		return
	}

	transaction, err := models.CreateTransaction(input.AccountID, input.OperationTypeID, input.Amount)
	if err != nil {
		utils.WriteJSONError(w, "Failed to create transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}
