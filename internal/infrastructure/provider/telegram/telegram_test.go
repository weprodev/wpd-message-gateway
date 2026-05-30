package telegram

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// TestNew verifies the constructor applies default BaseURL when none is given
// and uses a custom BaseURL when provided.
func TestNew(t *testing.T) {
	t.Run("defaults base URL when empty", func(t *testing.T) {
		p, err := New(Config{APIKey: "test-key"})
		if err != nil {
			t.Fatal(err)
		}
		if p.baseURL != defaultBaseURL {
			t.Errorf("baseURL = %q, want %q", p.baseURL, defaultBaseURL)
		}
		if p.apiKey != "test-key" {
			t.Errorf("apiKey = %q, want %q", p.apiKey, "test-key")
		}
	})

	t.Run("uses custom base URL", func(t *testing.T) {
		custom := "https://custom.example.com/"
		p, err := New(Config{APIKey: "test-key", BaseURL: custom})
		if err != nil {
			t.Fatal(err)
		}
		if p.baseURL != custom {
			t.Errorf("baseURL = %q, want %q", p.baseURL, custom)
		}
	})
}

// TestName checks that Name() returns the expected provider name constant.
func TestName(t *testing.T) {
	p, _ := New(Config{APIKey: "test-key"})
	if got := p.Name(); got != ProviderName {
		t.Errorf("Name() = %q, want %q", got, ProviderName)
	}
}

// TestNormalizeSenderUsernames exercises all branches of the pure function:
// empty usernames, single broadcast, matching counts, mismatched counts, zero phones.
func TestNormalizeSenderUsernames(t *testing.T) {
	tests := []struct {
		name       string
		phones     []string
		usernames  []string
		want       []string
		wantErr    bool
		errMessage string
	}{
		{
			name:      "empty usernames returns empty slice",
			phones:    []string{"+12025550123", "+12025550124"},
			usernames: []string{},
			want:      []string{"", ""},
		},
		{
			name:      "single username broadcast to all",
			phones:    []string{"+12025550123", "+12025550124", "+12025550125"},
			usernames: []string{"@bot"},
			want:      []string{"@bot", "@bot", "@bot"},
		},
		{
			name:      "matching username count",
			phones:    []string{"+12025550123", "+12025550124"},
			usernames: []string{"@bot1", "@bot2"},
			want:      []string{"@bot1", "@bot2"},
		},
		{
			name:       "mismatched count returns error",
			phones:     []string{"+12025550123", "+12025550124", "+12025550125"},
			usernames:  []string{"@bot1", "@bot2"},
			wantErr:    true,
			errMessage: "sender_username count (2) must be 0, 1, or match phone_number count (3)",
		},
		{
			name:      "single phone with empty usernames",
			phones:    []string{"+12025550123"},
			usernames: []string{},
			want:      []string{""},
		},
		{
			name:      "single phone with single username",
			phones:    []string{"+12025550123"},
			usernames: []string{"@bot"},
			want:      []string{"@bot"},
		},
		{
			name:      "zero phones returns empty slice",
			phones:    []string{},
			usernames: []string{},
			want:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeSenderUsernames(tt.phones, tt.usernames)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.errMessage {
					t.Errorf("error = %q, want %q", err.Error(), tt.errMessage)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestSend covers the public Send method end-to-end via a test HTTP server.
func TestSend(t *testing.T) {
	// Sending with an unparseable phone number should fail before any HTTP call.
	t.Run("invalid phone number returns error", func(t *testing.T) {
		p, _ := New(Config{APIKey: "test-key"})
		_, err := p.Send(context.Background(), &contracts.OTP{
			PhoneNumber: []string{"not-a-phone"},
		})
		if err == nil {
			t.Fatal("expected error for invalid phone number")
		}
	})

	// A single valid phone routes through sendSingle and returns a successful result.
	t.Run("single phone success", func(t *testing.T) {
		var req struct {
			PhoneNumber string `json:"phone_number"`
			CodeLength  int    `json:"code_length"`
		}
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			json.NewDecoder(r.Body).Decode(&req)
			return http.StatusOK, map[string]interface{}{
				"ok": true,
				"result": map[string]interface{}{
					"request_id": "req-123",
				},
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		result, err := p.Send(context.Background(), &contracts.OTP{
			PhoneNumber: []string{"+12025550123"},
			CodeLength:  6,
		})
		if err != nil {
			t.Fatal(err)
		}
		if result.ID != "req-123" {
			t.Errorf("ID = %q, want %q", result.ID, "req-123")
		}
		if result.StatusCode != 200 {
			t.Errorf("StatusCode = %d, want %d", result.StatusCode, 200)
		}
		if req.PhoneNumber != "+12025550123" {
			t.Errorf("phone_number = %q, want %q", req.PhoneNumber, "+12025550123")
		}
		if req.CodeLength != 6 {
			t.Errorf("code_length = %d, want %d", req.CodeLength, 6)
		}
	})

	// sender_username is passed through in the request body when provided.
	t.Run("single phone with sender username", func(t *testing.T) {
		var req struct {
			PhoneNumber    string `json:"phone_number"`
			SenderUsername string `json:"sender_username"`
		}
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			json.NewDecoder(r.Body).Decode(&req)
			return http.StatusOK, map[string]interface{}{
				"ok": true,
				"result": map[string]interface{}{
					"request_id": "req-456",
				},
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		result, err := p.Send(context.Background(), &contracts.OTP{
			PhoneNumber:    []string{"+12025550123"},
			SenderUsername: []string{"@mybot"},
		})
		if err != nil {
			t.Fatal(err)
		}
		if result.ID != "req-456" {
			t.Errorf("ID = %q, want %q", result.ID, "req-456")
		}
		if req.SenderUsername != "@mybot" {
			t.Errorf("sender_username = %q, want %q", req.SenderUsername, "@mybot")
		}
	})

	// When code_length is zero the field should be omitted from the request body.
	t.Run("zero code_length omits field", func(t *testing.T) {
		var req map[string]interface{}
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			json.NewDecoder(r.Body).Decode(&req)
			return http.StatusOK, map[string]interface{}{
				"ok": true,
				"result": map[string]interface{}{
					"request_id": "req-789",
				},
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		_, err := p.Send(context.Background(), &contracts.OTP{
			PhoneNumber: []string{"+12025550123"},
		})
		if err != nil {
			t.Fatal(err)
		}
		if _, exists := req["code_length"]; exists {
			t.Error("code_length should not be present when zero")
		}
	})

	// Multiple phones trigger the concurrent errgroup path; all items are populated.
	t.Run("multiple phones success", func(t *testing.T) {
		reqCh := make(chan map[string]interface{}, 2)
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			reqCh <- body
			return http.StatusOK, map[string]interface{}{
				"ok": true,
				"result": map[string]interface{}{
					"request_id": "req-" + body["phone_number"].(string),
				},
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		result, err := p.Send(context.Background(), &contracts.OTP{
			PhoneNumber: []string{"+12025550123", "+447911123456"},
			CodeLength:  4,
		})
		if err != nil {
			t.Fatal(err)
		}
		if result.StatusCode != 200 {
			t.Errorf("StatusCode = %d, want %d", result.StatusCode, 200)
		}
		if len(result.Items) != 2 {
			t.Fatalf("got %d items, want 2", len(result.Items))
		}
		close(reqCh)
		got := 0
		for range reqCh {
			got++
		}
		if got != 2 {
			t.Errorf("got %d requests, want 2", got)
		}
		reqIDs := map[string]bool{}
		for _, it := range result.Items {
			reqIDs[it.RequestID] = true
		}
		if !reqIDs["req-+12025550123"] {
			t.Error("missing request_id for +12025550123")
		}
		if !reqIDs["req-+447911123456"] {
			t.Error("missing request_id for +447911123456")
		}
	})

	// Every request includes the Authorization: Bearer <apiKey> header.
	t.Run("authorization header is set", func(t *testing.T) {
		var authHeader string
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			authHeader = r.Header.Get("Authorization")
			return http.StatusOK, map[string]interface{}{
				"ok": true,
				"result": map[string]interface{}{
					"request_id": "req-auth",
				},
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "my-secret-key", BaseURL: srv.URL + "/"})
		_, err := p.Send(context.Background(), &contracts.OTP{
			PhoneNumber: []string{"+12025550123"},
		})
		if err != nil {
			t.Fatal(err)
		}
		if authHeader != "Bearer my-secret-key" {
			t.Errorf("Authorization = %q, want %q", authHeader, "Bearer my-secret-key")
		}
	})

	// An API response with ok:false should be surfaced as a Go error.
	t.Run("API error response returns error", func(t *testing.T) {
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			return http.StatusOK, map[string]interface{}{
				"ok":    false,
				"error": "invalid phone number",
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		_, err := p.Send(context.Background(), &contracts.OTP{
			PhoneNumber: []string{"+12025550123"},
		})
		if err == nil {
			t.Fatal("expected error from API response")
		}
	})

	// In a batch send, a single failure is reflected in StatusCode 500 and per-item errors.
	t.Run("partial batch failure", func(t *testing.T) {
		var callCount atomic.Int32
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			n := callCount.Add(1)
			if n == 2 {
				return http.StatusOK, map[string]interface{}{
					"ok":    false,
					"error": "rate limited",
				}
			}
			return http.StatusOK, map[string]interface{}{
				"ok": true,
				"result": map[string]interface{}{
					"request_id": "ok",
				},
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		result, err := p.Send(context.Background(), &contracts.OTP{
			PhoneNumber: []string{"+12025550123", "+12025550124", "+12025550125"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("result should not be nil")
		}
		if result.StatusCode != 500 {
			t.Errorf("StatusCode = %d, want %d", result.StatusCode, 500)
		}
		if len(result.Items) != 3 {
			t.Fatalf("got %d items, want 3", len(result.Items))
		}
		hasFailure := false
		for _, it := range result.Items {
			if it.Error != "" {
				hasFailure = true
				break
			}
		}
		if !hasFailure {
			t.Error("expected at least one item with an error")
		}
	})
}

// TestCheckStatus verifies the CheckStatus endpoint: full field parsing,
// nil-safe handling of optional fields, and API error propagation.
func TestCheckStatus(t *testing.T) {
	// All response fields are correctly parsed, including nested delivery_status
	// and verification_status objects.
	t.Run("success with all fields", func(t *testing.T) {
		var reqBody struct {
			RequestID string `json:"request_id"`
		}
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			json.NewDecoder(r.Body).Decode(&reqBody)
			return http.StatusOK, map[string]interface{}{
				"ok": true,
				"result": map[string]interface{}{
					"request_id":       "req-123",
					"phone_number":     "+12025550123",
					"request_cost":     0.015,
					"is_refunded":      false,
					"remaining_balance": 99.985,
					"delivery_status": map[string]interface{}{
						"status":          "delivered",
						"last_updated_at": float64(1717000000),
					},
					"verification_status": map[string]interface{}{
						"status":      "verified",
						"verified_at": float64(1717000100),
					},
				},
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		status, err := p.CheckStatus(context.Background(), "req-123")
		if err != nil {
			t.Fatal(err)
		}
		if reqBody.RequestID != "req-123" {
			t.Errorf("request_id = %q, want %q", reqBody.RequestID, "req-123")
		}
		if status.RequestID != "req-123" {
			t.Errorf("RequestID = %q, want %q", status.RequestID, "req-123")
		}
		if status.PhoneNumber != "+12025550123" {
			t.Errorf("PhoneNumber = %q, want %q", status.PhoneNumber, "+12025550123")
		}
		if status.RequestCost != 0.015 {
			t.Errorf("RequestCost = %f, want %f", status.RequestCost, 0.015)
		}
		if status.IsRefunded == nil || *status.IsRefunded != false {
			t.Errorf("IsRefunded = %v, want false", status.IsRefunded)
		}
		if status.RemainingBalance == nil || *status.RemainingBalance != 99.985 {
			t.Errorf("RemainingBalance = %v, want 99.985", status.RemainingBalance)
		}
		if status.DeliveryStatus == nil || status.DeliveryStatus.Status != "delivered" {
			t.Errorf("DeliveryStatus.Status = %q, want %q", status.DeliveryStatus.Status, "delivered")
		}
		if status.VerificationStatus == nil || status.VerificationStatus.Status != "verified" {
			t.Errorf("VerificationStatus.Status = %q, want %q", status.VerificationStatus.Status, "verified")
		}
	})

	// When optional fields are absent the corresponding pointer fields stay nil.
	t.Run("success with minimal fields", func(t *testing.T) {
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			return http.StatusOK, map[string]interface{}{
				"ok": true,
				"result": map[string]interface{}{
					"request_id":   "req-minimal",
					"phone_number": "+12025550123",
					"request_cost": 0.0,
				},
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		status, err := p.CheckStatus(context.Background(), "req-minimal")
		if err != nil {
			t.Fatal(err)
		}
		if status.RequestID != "req-minimal" {
			t.Errorf("RequestID = %q, want %q", status.RequestID, "req-minimal")
		}
		if status.IsRefunded != nil {
			t.Error("IsRefunded should be nil when not present")
		}
		if status.RemainingBalance != nil {
			t.Error("RemainingBalance should be nil when not present")
		}
		if status.DeliveryStatus != nil {
			t.Error("DeliveryStatus should be nil when not present")
		}
		if status.VerificationStatus != nil {
			t.Error("VerificationStatus should be nil when not present")
		}
	})

	// An ok:false response propagates as a Go error.
	t.Run("API error returns error", func(t *testing.T) {
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			return http.StatusOK, map[string]interface{}{
				"ok":    false,
				"error": "request not found",
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		_, err := p.CheckStatus(context.Background(), "req-missing")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// TestRevoke verifies the Revoke endpoint: success path and error propagation.
func TestRevoke(t *testing.T) {
	// A successful revoke returns a SendResult with StatusCode 200.
	t.Run("success", func(t *testing.T) {
		var reqBody struct {
			RequestID string `json:"request_id"`
		}
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			json.NewDecoder(r.Body).Decode(&reqBody)
			return http.StatusOK, map[string]interface{}{
				"ok": true,
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		result, err := p.Revoke(context.Background(), "req-123")
		if err != nil {
			t.Fatal(err)
		}
		if reqBody.RequestID != "req-123" {
			t.Errorf("request_id = %q, want %q", reqBody.RequestID, "req-123")
		}
		if result.StatusCode != http.StatusOK {
			t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusOK)
		}
	})

	// An ok:false response propagates as a Go error.
	t.Run("API error returns error", func(t *testing.T) {
		srv := newTestServer(t, func(r *http.Request) (int, interface{}) {
			return http.StatusOK, map[string]interface{}{
				"ok":    false,
				"error": "invalid request_id",
			}
		})
		defer srv.Close()

		p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
		_, err := p.Revoke(context.Background(), "req-bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// TestSend_ContextCancellation ensures a cancelled context is respected by Send.
func TestSend_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := p.Send(ctx, &contracts.OTP{
		PhoneNumber: []string{"+12025550123"},
	})
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}

// TestCheckStatus_ContextCancellation ensures a cancelled context is respected by CheckStatus.
func TestCheckStatus_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := p.CheckStatus(ctx, "req-123")
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}

// TestRevoke_ContextCancellation ensures a cancelled context is respected by Revoke.
func TestRevoke_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	p, _ := New(Config{APIKey: "test-key", BaseURL: srv.URL + "/"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := p.Revoke(ctx, "req-123")
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}

// newTestServer creates an httptest.Server that calls handlerFn for every request.
// handlerFn receives the request and returns (statusCode, responseBody).
func newTestServer(t *testing.T, handlerFn func(*http.Request) (int, interface{})) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status, body := handlerFn(r)
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
}
