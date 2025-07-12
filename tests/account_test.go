package tests

import (
	"bank-api/handlers"
	"bank-api/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestAccountEndpoints(t *testing.T) {
	ResetDB()
	existingID, existingDocumentNumber := seedAccount(t, "12345678900")

	type want struct {
		code           int
		err            string
		accountID      int
		documentNumber string
	}
	type testCase struct {
		name, method, url, body string
		want
	}

	cases := []testCase{
		{
			name:   "Create OK",
			method: http.MethodPost,
			url:    "/accounts",
			body:   `{"document_number":"11122233344"}`,
			want:   want{code: http.StatusCreated, documentNumber: "11122233344"},
		},
		{
			name:   "Create missing document number",
			method: http.MethodPost,
			url:    "/accounts",
			body:   `{}`,
			want:   want{code: http.StatusBadRequest, err: utils.ErrMissingDocument},
		},
		{
			name:   "Get valid account",
			method: http.MethodGet,
			url:    "/accounts/" + strconv.Itoa(existingID),
			want:   want{code: http.StatusOK, accountID: existingID, documentNumber: existingDocumentNumber},
		},
		{
			name:   "Get invalid id",
			method: http.MethodGet,
			url:    "/accounts/abc",
			want:   want{code: http.StatusBadRequest, err: utils.ErrAccountNotFound},
		},
		{
			name:   "Get non existing id",
			method: http.MethodGet,
			url:    "/accounts/9999",
			want:   want{code: http.StatusNotFound, err: utils.ErrAccountNotFound},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			switch tc.method {
			case http.MethodPost:
				handlers.CreateAccountHandler(w, req)
			case http.MethodGet:
				handlers.GetAccountHandler(w, req)
			}

			if w.Code != tc.want.code {
				t.Fatalf("want status %d, got %d", tc.want.code, w.Code)
			}

			if tc.want.err != "" {
				var m map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
					t.Fatalf("invalid JSON error: %v", err)
				}
				if m["error"] != tc.want.err {
					t.Errorf("expected error %q, got %q", tc.want.err, m["error"])
				}
				return
			}

			var resp struct {
				AccountID      int    `json:"account_id"`
				DocumentNumber string `json:"document_number"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("invalid JSON response: %v", err)
			}

			if tc.want.accountID != 0 && resp.AccountID != tc.want.accountID {
				t.Errorf("expected account_id %d, got %d", tc.want.accountID, resp.AccountID)
			}
			if tc.want.documentNumber != "" && resp.DocumentNumber != tc.want.documentNumber {
				t.Errorf("expected document_number %q, got %q", tc.want.documentNumber, resp.DocumentNumber)
			}
		})
	}
}
