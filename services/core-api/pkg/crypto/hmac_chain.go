// Package crypto provides enterprise cryptographic functions adhering to the
// NESCO Utility Platform specification. It provides SHA-256 hash chaining for
// immutable WorkNotes and audit records, row-level HMAC signatures for sensitive
// database entities, RFC 4122 v4 UUIDs, and salted password hashing.
package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	// ErrEmptySecret indicates a secret key was not provided for HMAC computation.
	ErrEmptySecret = errors.New("crypto: hmac secret key cannot be empty")
	// ErrInvalidHashFormat indicates a hash string format violation.
	ErrInvalidHashFormat = errors.New("crypto: invalid hash format")
)

// GenesisHash represents the initial seed hash for brand-new hash chains.
const GenesisHash = "0000000000000000000000000000000000000000000000000000000000000000"

// ComputeWorkNoteHash calculates an unalterable SHA-256 hash chaining a new WorkNote
// to its preceding ledger entry.
// Formula: SHA256(prev_hash + ":" + entity_type + ":" + entity_id + ":" + author_id + ":" + note_text + ":" + timestamp)
func ComputeWorkNoteHash(prevHash, entityType, entityID, authorID, noteText string, timestamp int64) string {
	if strings.TrimSpace(prevHash) == "" {
		prevHash = GenesisHash
	}
	payload := fmt.Sprintf("%s:%s:%s:%s:%s:%d",
		prevHash, entityType, entityID, authorID, noteText, timestamp)

	h := sha256.New()
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyWorkNoteHash validates that a given currentHash matches the expected SHA-256
// chaining from the given prevHash and content fields.
func VerifyWorkNoteHash(currentHash, prevHash, entityType, entityID, authorID, noteText string, timestamp int64) bool {
	expected := ComputeWorkNoteHash(prevHash, entityType, entityID, authorID, noteText, timestamp)
	return subtle.ConstantTimeCompare([]byte(currentHash), []byte(expected)) == 1
}

// ComputeRowHMAC calculates a cryptographic HMAC-SHA256 signature over a serialized
// database record payload to guarantee tamper evidence.
func ComputeRowHMAC(secretKey, payload string) (string, error) {
	if secretKey == "" {
		return "", ErrEmptySecret
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// VerifyRowHMAC performs a constant-time cryptographic verification of an HMAC signature.
func VerifyRowHMAC(secretKey, payload, expectedHMAC string) bool {
	if secretKey == "" || expectedHMAC == "" {
		return false
	}
	actual, err := ComputeRowHMAC(secretKey, payload)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(actual), []byte(expectedHMAC)) == 1
}

// NewUUID generates a cryptographically secure random RFC 4122 version 4 UUID.
func NewUUID() (string, error) {
	var uuid [16]byte
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		return "", fmt.Errorf("crypto: failed to read secure random bytes: %w", err)
	}
	// Set UUID version 4
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set variant bits (RFC 4122)
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}

// MustUUID generates a UUID and panics if the secure random reader encounters an error.
func MustUUID() string {
	u, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return u
}

// HashPassword generates a salted SHA-256 password digest formatted as:
// "v1$<hex_salt>$<hex_digest>" using 16 bytes of cryptographically secure random salt.
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("crypto: failed to generate salt: %w", err)
	}

	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(password))
	digest := h.Sum(nil)

	return fmt.Sprintf("v1$%s$%s", hex.EncodeToString(salt), hex.EncodeToString(digest)), nil
}

// VerifyPassword verifies whether a plain text password matches a salted password digest.
func VerifyPassword(password, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 3 || parts[0] != "v1" {
		return false
	}

	salt, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(password))
	expectedDigest := hex.EncodeToString(h.Sum(nil))

	return subtle.ConstantTimeCompare([]byte(parts[2]), []byte(expectedDigest)) == 1
}

