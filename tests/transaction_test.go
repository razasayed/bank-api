package tests

import (
	"bank-api/handlers"
	"bank-api/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransactionEndpoints(t *testing.T) {
	ResetDB()

	accountID, _ := seedAccount(t, "12345678900")

	type want struct {
		code            int
		err             string
		accountID       int
		operationTypeID int
		amount          float64
	}
	type testCase struct {
		name, body string
		want
	}

	cases := []testCase{
		{
			name: "Purchase transaction",
			body: `{"account_id":%d,"operation_type_id":1,"amount":100}`,
			want: want{
				code:            http.StatusCreated,
				accountID:       accountID,
				operationTypeID: utils.OperationTypePurchase,
				amount:          -100,
			},
		},
		{
			name: "Installment Purchase transaction",
			body: `{"account_id":%d,"operation_type_id":2,"amount":100}`,
			want: want{
				code:            http.StatusCreated,
				accountID:       accountID,
				operationTypeID: utils.OperationTypeInstallmentPurchase,
				amount:          -100,
			},
		},
		{
			name: "Withdrawal transaction",
			body: `{"account_id":%d,"operation_type_id":3,"amount":100}`,
			want: want{
				code:            http.StatusCreated,
				accountID:       accountID,
				operationTypeID: utils.OperationTypeWithdrawal,
				amount:          -100,
			},
		},
		{
			name: "Payment transaction",
			body: `{"account_id":%d,"operation_type_id":4,"amount":200}`,
			want: want{
				code:            http.StatusCreated,
				accountID:       accountID,
				operationTypeID: utils.OperationTypePayment,
				amount:          200,
			},
		},
		{
			name: "Missing operation_type_id",
			body: `{"account_id":%d,"amount":100}`,
			want: want{
				code: http.StatusBadRequest,
				err:  utils.ErrMissingAccountFields,
			},
		},
		{
			name: "Invalid operation_type_id",
			body: `{"account_id":%d,"operation_type_id":99,"amount":100}`,
			want: want{
				code: http.StatusBadRequest,
				err:  utils.ErrInvalidOperationTypeID,
			},
		},
		{
			name: "Missing account_id",
			body: `{"operation_type_id":1,"amount":50}`,
			want: want{
				code: http.StatusBadRequest,
				err:  utils.ErrMissingAccountFields,
			},
		},
		{
			name: "Non existent account_id",
			body: `{"account_id":9999,"operation_type_id":1,"amount":50}`,
			want: want{
				code: http.StatusNotFound,
				err:  utils.ErrAccountNotFound,
			},
		},
		{
			name: "Missing amount",
			body: `{"operation_type_id":1,"account_id":%d}`,
			want: want{
				code: http.StatusBadRequest,
				err:  utils.ErrMissingAccountFields,
			},
		},
		{
			name: "Negative amount for payment",
			body: `{"account_id":%d,"operation_type_id":4,"amount":-50}`,
			want: want{
				code: http.StatusBadRequest,
				err:  utils.ErrInvalidAmount,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := fmt.Sprintf(tc.body, accountID)
			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handlers.CreateTransactionHandler(rec, req)

			if rec.Code != tc.want.code {
				t.Fatalf("expected status %d, got %d", tc.want.code, rec.Code)
			}

			if tc.want.err != "" {
				var m map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &m); err != nil {
					t.Fatalf("invalid error JSON: %v", err)
				}
				if m["error"] != tc.want.err {
					t.Errorf("expected error %q, got %q", tc.want.err, m["error"])
				}
				return
			}

			var resp struct {
				TransactionID   int     `json:"transaction_id"`
				AccountID       int     `json:"account_id"`
				OperationTypeID int     `json:"operation_type_id"`
				Amount          float64 `json:"amount"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("invalid JSON: %v", err)
			}

			if resp.TransactionID == 0 {
				t.Error("expected a valid transaction_id, got 0")
			}
			if tc.want.accountID != 0 && resp.AccountID != tc.want.accountID {
				t.Errorf("expected account_id %d, got %d", tc.want.accountID, resp.AccountID)
			}
			if tc.want.operationTypeID != 0 && resp.OperationTypeID != tc.want.operationTypeID {
				t.Errorf("expected operation_type_id %d, got %d", tc.want.operationTypeID, resp.OperationTypeID)
			}
			if tc.want.amount != 0 && resp.Amount != tc.want.amount {
				t.Errorf("expected amount %v, got %v", tc.want.amount, resp.Amount)
			}
		})
	}
}
