package utils

// http error constants
const (
	ErrMethodNotAllowed     = "Method not allowed"
	ErrInvalidJSON          = "Invalid JSON body"
	ErrMissingAccountFields = "account_id, operation_type_id, and amount are required"
	ErrMissingDocument      = "Document number is required"
	ErrInvalidAccountID     = "Invalid account ID"
	ErrAccountNotFound      = "Account not found"
)
