package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestRegisterReward(t *testing.T) {
	container := setupContainer()
	defer container.Terminate(context.Background())

	db, err := setupDB(container)
	if err != nil {
		t.Fatalf("failed to set up database: %v", err)
	}
	defer db.Close()

	err = runMigrations(db)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	serverCmd := startServer(container)
	defer stopServer(serverCmd)

	tests := []struct {
		name         string
		payload      interface{}
		expectedCode int
		expectedBody string
	}{
		{
			name:         "valid reward registration",
			payload:      map[string]interface{}{"match": "TestItem", "reward": 15, "reward_type": "%"},
			expectedCode: 200,
			expectedBody: "Reward registered successfully",
		},
		{
			name:         "missing match parameter",
			payload:      map[string]interface{}{"reward": 15, "reward_type": "%"},
			expectedCode: 400,
			expectedBody: "Missing match parameter",
		},
		{
			name:         "invalid reward type",
			payload:      map[string]interface{}{"match": "TestItem", "reward": 15, "reward_type": "dg"},
			expectedCode: 400,
			expectedBody: "Invalid reward type parameter",
		},
		{
			name:         "missing reward parameter",
			payload:      map[string]interface{}{"match": "TestItem", "reward_type": "%"},
			expectedCode: 400,
			expectedBody: "Missing reward parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer clearTable(db)

			client := resty.New()
			resp, err := client.R().
				SetBody(tt.payload).
				Post("http://localhost:8081/api/goods")

			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}

			// Check status code
			if resp.StatusCode() != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, resp.StatusCode())
			}

			// Trim whitespace and compare bodies
			actualBody := strings.TrimSpace(string(resp.Body()))
			expectedBody := strings.TrimSpace(tt.expectedBody)

			// Add debug logs to see the actual response body
			if actualBody != expectedBody {
				t.Errorf("expected body %q, got %q", expectedBody, actualBody)
				t.Logf("Actual response body (trimmed): %q", actualBody)
			}
		})
	}
}
