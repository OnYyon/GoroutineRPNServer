package orchestrator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddNewExpression(t *testing.T) {
	api := &API{
		Expressions: make(map[string]*Expression),
	}

	tests := []struct {
		name             string
		input            string
		expectedStatus   int
		expectedResponse map[string]string
		answer           float64
	}{
		{
			name:           "valid expression",
			input:          `{"expression": "2 + 2"}`,
			expectedStatus: http.StatusCreated,
			answer:         4,
		},
		{
			name:           "invalid expression format",
			input:          `{"expr": "2 + 2"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: map[string]string{
				"error": "Oppps something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/v1/expressions", bytes.NewBufferString(tt.input))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(api.AddNewExpression)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetSliceOfExpressions(t *testing.T) {
	api := &API{
		Expressions: make(map[string]*Expression),
	}
	expression := &Expression{
		ID:     "1",
		Status: StatusNew,
		Input:  "2 + 2",
	}
	api.Expressions[expression.ID] = expression

	req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetSliceOfExpressions)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	var response struct {
		Expressions []Expression `json:"expressions"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response.Expressions) != 1 || response.Expressions[0].ID != expression.ID {
		t.Errorf("expected expressions to be %v but got %v", expression, response.Expressions)
	}
}
