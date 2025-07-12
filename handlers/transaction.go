package handlers

import (
	"bank-api/models"
	"bank-api/utils"
	"encoding/json"
	"net/http"
)

// CreateTransactionHandler godoc
// @Summary Create a new transaction
// @Description Creates a transaction for a given account with a valid operation type and amount
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param transaction body models.CreateTransactionInput true "Create Transaction input"
// @Success 201 {object} models.Transaction
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 405 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /transactions [post]
func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, utils.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var input models.CreateTransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteJSONError(w, utils.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if input.AccountID == 0 || input.OperationTypeID == 0 || input.Amount == 0 {
		utils.WriteJSONError(w, utils.ErrMissingAccountFields, http.StatusBadRequest)
		return
	}

	if input.Amount < 0 {
		utils.WriteJSONError(w, utils.ErrInvalidAmount, http.StatusBadRequest)
		return
	}

	if !models.OperationTypeExists(input.OperationTypeID) {
		utils.WriteJSONError(w, utils.ErrInvalidOperationTypeID, http.StatusBadRequest)
		return
	}

	transaction, err := models.CreateTransaction(input.AccountID, input.OperationTypeID, input.Amount)
	if err != nil {
		if err.Error() == utils.ErrAccountNotFound {
			utils.WriteJSONError(w, utils.ErrAccountNotFound, http.StatusNotFound)
			return
		}
		utils.WriteJSONError(w, "Failed to create transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}
