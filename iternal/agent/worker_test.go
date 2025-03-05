package agent

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/OnYyon/GoroutineRPNServer/iternal/orchestrator"
)

func TestEvaluateExpression(t *testing.T) {
	tests := []struct {
		name       string
		task       *orchestrator.Task
		expected   float64
		expectErr  bool
		errorValue error
	}{
		{
			name:      "valid addition",
			task:      &orchestrator.Task{Arg1: "3", Arg2: "4", Operation: "+"},
			expected:  7,
			expectErr: false,
		},
		{
			name:      "valid subtraction",
			task:      &orchestrator.Task{Arg1: "5", Arg2: "3", Operation: "-"},
			expected:  2,
			expectErr: false,
		},
		{
			name:      "valid multiplication",
			task:      &orchestrator.Task{Arg1: "2", Arg2: "3", Operation: "*"},
			expected:  6,
			expectErr: false,
		},
		{
			name:      "valid division",
			task:      &orchestrator.Task{Arg1: "6", Arg2: "2", Operation: "/"},
			expected:  3,
			expectErr: false,
		},
		{
			name:       "division by zero",
			task:       &orchestrator.Task{Arg1: "6", Arg2: "0", Operation: "/"},
			expected:   0,
			expectErr:  true,
			errorValue: fmt.Errorf("zero dividion"),
		},
		{
			name:      "invalid argument 1",
			task:      &orchestrator.Task{Arg1: "invalid", Arg2: "3", Operation: "+"},
			expected:  0,
			expectErr: true,
		},
		{
			name:      "invalid argument 2",
			task:      &orchestrator.Task{Arg1: "3", Arg2: "invalid", Operation: "+"},
			expected:  0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluateExpression(tt.task)

			if tt.expectErr && err == nil {
				t.Errorf("expected an error but got none")
			} else if !tt.expectErr && err != nil {
				t.Errorf("did not expect an error but got: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected result %v but got %v", tt.expected, result)
			}
		})
	}
}
func TestFetchTask(t *testing.T) {
	tests := []struct {
		name             string
		mockResponseBody string
		mockResponseCode int
		expectedTask     *orchestrator.Task
		expectedError    error
	}{
		{
			name:             "successful task fetch",
			mockResponseBody: `{"ID":"1", "Arg1":"2", "Arg2":"2", "Operation":"+"}`,
			mockResponseCode: http.StatusOK,
			expectedTask:     &orchestrator.Task{ID: "1", Arg1: "2", Arg2: "2", Operation: "+"},
			expectedError:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockResponseCode)
				w.Write([]byte(tt.mockResponseBody))
			}))
			defer mockServer.Close()

			parsedURL, err := url.Parse(mockServer.URL)
			if err != nil {
				t.Fatalf("failed to parse URL: %v", err)
			}

			defer func() { http.DefaultClient = &http.Client{} }()
			http.DefaultClient = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(parsedURL),
				},
			}

			task, err := FetchTask()

			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v but got %v", tt.expectedError, err)
			}
			if err == nil && *task != *tt.expectedTask {
				t.Errorf("expected task %+v but got %+v", tt.expectedTask, task)
			}
		})
	}
}

func TestSendResult(t *testing.T) {
	tests := []struct {
		name             string
		result           Res
		mockResponseCode int
		mockResponseBody string
		expectedError    error
	}{
		{
			name:             "successful result send",
			result:           Res{ID: "1", Result: 5, Timeout: false, Errors: ""},
			mockResponseCode: http.StatusOK,
			mockResponseBody: "",
			expectedError:    nil,
		},
		{
			name:             "send result returns error",
			result:           Res{ID: "1", Result: 5, Timeout: false, Errors: ""},
			mockResponseCode: http.StatusInternalServerError,
			mockResponseBody: "",
			expectedError:    errors.New("unexpected status code: 500"),
		},
		{
			name:             "failed to marshal result",
			result:           Res{ID: "", Result: 0, Timeout: false, Errors: ""},
			mockResponseCode: http.StatusOK,
			mockResponseBody: "",
			expectedError:    errors.New("json: cannot marshal object"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockResponseCode)
				w.Write([]byte(tt.mockResponseBody))
			}))
			defer mockServer.Close()

			parsedURL, err := url.Parse(mockServer.URL)
			if err != nil {
				t.Fatalf("failed to parse URL: %v", err)
			}

			defer func() { http.DefaultClient = &http.Client{} }()
			http.DefaultClient = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(parsedURL),
				},
			}

			err = SendResult(tt.result)

			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v but got %v", tt.expectedError, err)
			}
		})
	}
}
