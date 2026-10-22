// Package crypto_test validates cryptographic hashing, chaining, HMAC and password routines.
package crypto_test

import (
	"strings"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/pkg/crypto"
)

func TestComputeWorkNoteHash_DeterministicAndChained(t *testing.T) {
	ts := time.Now().Unix()
	hash1 := crypto.ComputeWorkNoteHash("", "TICKET", "TCK-1001", "STAFF-01", "Initial fault reported", ts)
	if len(hash1) != 64 {
		t.Fatalf("expected 64 char hex hash, got %d", len(hash1))
	}

	// Chaining test: next note uses hash1 as prev_hash
	hash2 := crypto.ComputeWorkNoteHash(hash1, "TICKET", "TCK-1001", "STAFF-02", "Investigating onsite", ts+60)
	if hash1 == hash2 {
		t.Fatalf("expected different hash values across chain, got duplicate: %s", hash1)
	}

	// Verify valid chain verification
	valid := crypto.VerifyWorkNoteHash(hash1, crypto.GenesisHash, "TICKET", "TCK-1001", "STAFF-01", "Initial fault reported", ts)
	if !valid {
		t.Fatalf("expected hash1 verification to succeed")
	}

	// Negative security test: tampered payload must fail
	tampered := crypto.VerifyWorkNoteHash(hash1, crypto.GenesisHash, "TICKET", "TCK-1001", "STAFF-01", "Tampered text", ts)
	if tampered {
		t.Fatalf("tampered payload must not verify against original hash")
	}
}

func TestRowHMAC_SignAndVerify(t *testing.T) {
	secret := "master-enterprise-secret-key-32b!"
	payload := `{"asset_code":"TR-JOS-01","capacity":500,"state":"IN_SERVICE"}`

	hmacSig, err := crypto.ComputeRowHMAC(secret, payload)
	if err != nil {
		t.Fatalf("unexpected HMAC error: %v", err)
	}
	if len(hmacSig) != 64 {
		t.Fatalf("expected 64 char HMAC, got %d", len(hmacSig))
	}

	// Positive verification
	if !crypto.VerifyRowHMAC(secret, payload, hmacSig) {
		t.Fatalf("expected valid HMAC verification to pass")
	}

	// Negative test: Tampered payload
	tamperedPayload := `{"asset_code":"TR-JOS-01","capacity":999,"state":"IN_SERVICE"}`
	if crypto.VerifyRowHMAC(secret, tamperedPayload, hmacSig) {
		t.Fatalf("HMAC verification must fail for tampered payload")
	}

	// Negative test: Wrong secret
	if crypto.VerifyRowHMAC("wrong-secret-key", payload, hmacSig) {
		t.Fatalf("HMAC verification must fail with wrong secret")
	}

	// Negative test: Empty secret handling
	_, err = crypto.ComputeRowHMAC("", payload)
	if err == nil {
		t.Fatalf("expected error on empty secret key")
	}
}

func TestUUID_Generation(t *testing.T) {
	u1 := crypto.MustUUID()
	u2 := crypto.MustUUID()

	if u1 == u2 {
		t.Fatalf("UUID collision detected: %s vs %s", u1, u2)
	}
	if len(u1) != 36 || !strings.Contains(u1, "-") {
		t.Fatalf("invalid UUID format: %s", u1)
	}
}

func TestPasswordHashing_SaltAndVerify(t *testing.T) {
	password := "SecureNescoAdmin#2026"
	encodedHash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !crypto.VerifyPassword(password, encodedHash) {
		t.Fatalf("valid password must verify")
	}

	// Negative test: incorrect password
	if crypto.VerifyPassword("WrongPassword#123", encodedHash) {
		t.Fatalf("invalid password must fail verification")
	}

	// Negative test: malformed hash
	if crypto.VerifyPassword(password, "invalid$format") {
		t.Fatalf("malformed hash must fail safely")
	}
}

