// Package auth implements enterprise JWT authentication and Role-Based Access
// Control (RBAC) middleware for the NESCO Utility Platform.
// It verifies staff identities, extracts security claims, and enforces
// department-level access controls across Executive, Unit Head, and Field tiers.
package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	// ErrInvalidToken indicates an ill-formed or tampered JWT token.
	ErrInvalidToken = errors.New("auth: invalid or malformed token")
	// ErrExpiredToken indicates token validity window has passed.
	ErrExpiredToken = errors.New("auth: token has expired")
	// ErrUnauthorized indicates missing or invalid authentication credentials.
	ErrUnauthorized = errors.New("auth: unauthorized request")
	// ErrForbidden indicates lack of required RBAC role permissions.
	ErrForbidden = errors.New("auth: insufficient permissions for this operation")
)

type contextKey string

const (
	claimsContextKey contextKey = "nesco_auth_claims"
	// TokenDuration defines the strict 15-minute access token lifespan.
	TokenDuration = 15 * time.Minute
)

// Claims defines the authenticated staff context encoded into the JWT payload.
type Claims struct {
	StaffID                    string `json:"staff_id"`
	StaffNumber                string `json:"staff_number"`
	FullName                   string `json:"full_name"`
	PrimaryRole                string `json:"primary_role"`
	OrgUnitID                  string `json:"org_unit_id"`
	CanViewExecutiveDashboard bool   `json:"can_view_executive_dashboard"`
	AssignedAssetID            string `json:"assigned_asset_id,omitempty"`
	IssuedAt                   int64  `json:"iat"`
	ExpiresAt                  int64  `json:"exp"`
}

// Valid validates whether the claims are active and non-expired.
func (c *Claims) Valid() error {
	if c.StaffID == "" || c.PrimaryRole == "" {
		return ErrInvalidToken
	}
	if time.Now().Unix() > c.ExpiresAt {
		return ErrExpiredToken
	}
	return nil
}

// Header represents the standard JWT header specifying the HS256 algorithm.
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// GenerateToken signs a new HMAC-SHA256 JWT using the provided secret key.
func GenerateToken(claims Claims, secret string) (string, error) {
	if secret == "" {
		return "", errors.New("auth: secret key cannot be empty")
	}

	header := Header{Alg: "HS256", Typ: "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("auth: failed to marshal header: %w", err)
	}

	if claims.IssuedAt == 0 {
		claims.IssuedAt = time.Now().Unix()
	}
	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = claims.IssuedAt + int64(TokenDuration.Seconds())
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("auth: failed to marshal claims: %w", err)
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsignedToken := headerEncoded + "." + claimsEncoded

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsignedToken))
	sigEncoded := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return unsignedToken + "." + sigEncoded, nil
}

// ParseAndValidateToken decodes a token, validates signature and expiration, and extracts claims.
func ParseAndValidateToken(tokenStr, secret string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("auth: secret key cannot be empty")
	}

	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	unsignedToken := parts[0] + "." + parts[1]
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsignedToken))
	expectedSig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, ErrInvalidToken
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if err := claims.Valid(); err != nil {
		return nil, err
	}

	return &claims, nil
}

// ContextWithClaims attaches Claims to the request context.
func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// ClaimsFromContext extracts authenticated Claims from the request context.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok && claims != nil
}

// RequireAuth returns an HTTP middleware verifying the Bearer JWT token in Authorization header.
func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"unauthorized","message":"missing or malformed authorization header"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := ParseAndValidateToken(tokenStr, secret)
			if err != nil {
				if errors.Is(err, ErrExpiredToken) {
					http.Error(w, `{"error":"token_expired","message":"session has expired, please log in again"}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error":"invalid_token","message":"token validation failed"}`, http.StatusUnauthorized)
				return
			}

			ctx := ContextWithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRoles returns an HTTP middleware checking if the authenticated user has one of the allowed roles.
func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				http.Error(w, `{"error":"unauthorized","message":"missing authentication context"}`, http.StatusUnauthorized)
				return
			}

			if _, allowed := roleSet[claims.PrimaryRole]; !allowed {
				http.Error(w, `{"error":"forbidden","message":"access denied for this role"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// IsExecutiveRole returns true if the specified role belongs to the C-Suite Executive Command Tier.
func IsExecutiveRole(role string) bool {
	switch role {
	case "MD", "HEAD_SERVICES", "HEAD_INTERNAL_AUDIT", "HEAD_COMMERCIAL_OPS", "GM_GENERATION", "HEAD_ENERGY_MGMT":
		return true
	default:
		return false
	}
}

