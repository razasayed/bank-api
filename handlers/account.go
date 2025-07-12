package handlers

import (
	"bank-api/models"
	"bank-api/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// CreateAccountHandler godoc
// @Summary Create a new account
// @Description Creates a customer account with the provided document number
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param account body models.CreateAccountInput true "Create Account input"
// @Success 201 {object} models.Account
// @Failure 400 {object} models.ErrorResponse
// @Failure 405 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts [post]
func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, utils.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}

	var input struct {
		DocumentNumber string `json:"document_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteJSONError(w, utils.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if input.DocumentNumber == "" {
		utils.WriteJSONError(w, utils.ErrMissingDocument, http.StatusBadRequest)
		return
	}

	account, err := models.CreateAccount(input.DocumentNumber)
	if err != nil {
		utils.WriteJSONError(w, "Failed to create account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetAccountHandler godoc
// @Summary Get account by ID
// @Description Retrieve account information by account ID
// @Tags accounts
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} models.Account
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 405 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts/{id} [get]
func GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, utils.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/accounts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, utils.ErrAccountNotFound, http.StatusBadRequest)
		return
	}

	account, err := models.GetAccountByID(id)
	if err != nil {
		if errors.Is(err, models.ErrAccountNotFound) {
			utils.WriteJSONError(w, utils.ErrAccountNotFound, http.StatusNotFound)
		} else {
			utils.WriteJSONError(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}
