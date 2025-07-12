package utils

// http error constants
const (
	ErrMethodNotAllowed       = "method not allowed"
	ErrInvalidJSON            = "invalid JSON body"
	ErrMissingAccountFields   = "account_id, operation_type_id, and amount are required"
	ErrMissingDocument        = "document number is required"
	ErrAccountNotFound        = "account not found"
	ErrInvalidOperationTypeID = "invalid operation type id"
	ErrMissingOperationTypeID = "operation type ID is required"
	ErrInvalidAmount          = "amount must be greater than zero"
)

const OperationTypePurchase = 1
const OperationTypeInstallmentPurchase = 2
const OperationTypeWithdrawal = 3
const OperationTypePayment = 4
