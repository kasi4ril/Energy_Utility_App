// Package auth_test tests JWT creation, parsing, middleware, and RBAC authorization.
package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/pkg/auth"
)

func TestGenerateAndValidateToken_Success(t *testing.T) {
	secret := "test-secret-key-32-characters-minimum!"
	claims := auth.Claims{
		StaffID:                   "staff-uuid-001",
		StaffNumber:               "NSC-ENG-042",
		FullName:                  "Engr. Danladi Musa",
		PrimaryRole:               "HEAD_SERVICES",
		OrgUnitID:                 "unit-uuid-hq",
		CanViewExecutiveDashboard: true,
	}

	token, err := auth.GenerateToken(claims, secret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	parsed, err := auth.ParseAndValidateToken(token, secret)
	if err != nil {
		t.Fatalf("failed to parse valid token: %v", err)
	}

	if parsed.StaffID != claims.StaffID || parsed.PrimaryRole != claims.PrimaryRole {
		t.Fatalf("mismatched claims: got %+v, want %+v", parsed, claims)
	}
	if !parsed.CanViewExecutiveDashboard {
		t.Fatalf("expected CanViewExecutiveDashboard to be true")
	}
}

func TestTokenValidation_TamperedAndExpired(t *testing.T) {
	secret := "test-secret-key-32-characters-minimum!"
	claims := auth.Claims{
		StaffID:     "staff-uuid-002",
		PrimaryRole: "CUSTOMER_CARE_OFFICER",
		IssuedAt:    time.Now().Unix() - 3600,
		ExpiresAt:   time.Now().Unix() - 1800, // Expired
	}

	expiredToken, err := auth.GenerateToken(claims, secret)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	_, err = auth.ParseAndValidateToken(expiredToken, secret)
	if err == nil {
		t.Fatalf("expected error on expired token, got nil")
	}

	// Tampered token test
	validClaims := auth.Claims{StaffID: "staff-uuid-003", PrimaryRole: "STORE_OFFICER"}
	validToken, _ := auth.GenerateToken(validClaims, secret)
	tamperedToken := validToken + "tamper"

	_, err = auth.ParseAndValidateToken(tamperedToken, secret)
	if err == nil {
		t.Fatalf("expected error on tampered token")
	}

	// Wrong secret test
	_, err = auth.ParseAndValidateToken(validToken, "completely-wrong-secret-key")
	if err == nil {
		t.Fatalf("expected signature verification failure on wrong secret")
	}
}

func TestRequireAuthMiddleware(t *testing.T) {
	secret := "super-secure-key-123456789012345"
	middleware := auth.RequireAuth(secret)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := auth.ClaimsFromContext(r.Context())
		if !ok || claims == nil {
			t.Fatalf("claims missing from context in downstream handler")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("authorized"))
	})

	handler := middleware(nextHandler)

	// Test 1: No Authorization header -> 401
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}

	// Test 2: Valid Bearer token -> 200
	token, _ := auth.GenerateToken(auth.Claims{
		StaffID:     "uuid-123",
		PrimaryRole: "MD",
	}, secret)

	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestRequireRolesMiddleware(t *testing.T) {
	secret := "super-secure-key-123456789012345"
	authMW := auth.RequireAuth(secret)
	roleMW := auth.RequireRoles("MD", "HEAD_SERVICES")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	pipeline := authMW(roleMW(nextHandler))

	// User with role "STORE_OFFICER" (Not allowed) -> 403 Forbidden
	tokenStoreOfficer, _ := auth.GenerateToken(auth.Claims{
		StaffID:     "uuid-store",
		PrimaryRole: "STORE_OFFICER",
	}, secret)

	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStoreOfficer)
	rec := httptest.NewRecorder()
	pipeline.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for STORE_OFFICER, got %d", rec.Code)
	}

	// User with role "MD" (Allowed) -> 200 OK
	tokenMD, _ := auth.GenerateToken(auth.Claims{
		StaffID:     "uuid-md",
		PrimaryRole: "MD",
	}, secret)

	req = httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+tokenMD)
	rec = httptest.NewRecorder()
	pipeline.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for MD, got %d", rec.Code)
	}
}

func TestIsExecutiveRole(t *testing.T) {
	executives := []string{"MD", "HEAD_SERVICES", "HEAD_INTERNAL_AUDIT", "HEAD_COMMERCIAL_OPS", "GM_GENERATION", "HEAD_ENERGY_MGMT"}
	for _, role := range executives {
		if !auth.IsExecutiveRole(role) {
			t.Fatalf("expected role %s to be recognized as executive", role)
		}
	}

	nonExecutives := []string{"CUSTOMER_CARE_OFFICER", "STORE_OFFICER", "TEAM_LEAD_FEEDER"}
	for _, role := range nonExecutives {
		if auth.IsExecutiveRole(role) {
			t.Fatalf("expected role %s NOT to be executive", role)
		}
	}
}
