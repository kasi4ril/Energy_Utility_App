// Package transport provides request/response data transfer objects and HTTP endpoints
// for the S02 Ollama Local AI Guardian module.
package transport

import (
	"errors"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/domain"
)

// SubmitMakerCheckerRequestDTO defines payload for POST /api/v1/maker-checker/request.
type SubmitMakerCheckerRequestDTO struct {
	Module         string                `json:"module"`
	ChangeCategory domain.ChangeCategory `json:"change_category"`
	TargetRecordID string                `json:"target_record_id"`
	PayloadJSON    string                `json:"payload_json"`
	Justification  string                `json:"justification"`
}

// Validate checks that required submission fields are populated.
func (d SubmitMakerCheckerRequestDTO) Validate() error {
	if strings.TrimSpace(d.Module) == "" || strings.TrimSpace(string(d.ChangeCategory)) == "" ||
		strings.TrimSpace(d.TargetRecordID) == "" || strings.TrimSpace(d.PayloadJSON) == "" ||
		strings.TrimSpace(d.Justification) == "" {
		return errors.New("transport: module, change_category, target_record_id, payload_json, and justification are mandatory")
	}
	return nil
}

// ApproveRequestDTO defines payload for POST /api/v1/maker-checker/{id}/approve.
type ApproveRequestDTO struct {
	CheckerNote string `json:"checker_note"`
}

// RejectRequestDTO defines payload for POST /api/v1/maker-checker/{id}/reject.
type RejectRequestDTO struct {
	Reason string `json:"reason"`
}

// Validate checks rejection reason is provided.
func (d RejectRequestDTO) Validate() error {
	if strings.TrimSpace(d.Reason) == "" {
		return errors.New("transport: rejection reason cannot be empty")
	}
	return nil
}

// ChatMessageDTO defines payload for POST /api/v1/chat/message.
type ChatMessageDTO struct {
	ChannelID   string `json:"channel_id"`
	MessageText string `json:"message_text"`
}

// Validate ensures non-empty message content.
func (d ChatMessageDTO) Validate() error {
	if strings.TrimSpace(d.MessageText) == "" {
		return errors.New("transport: message_text cannot be empty")
	}
	return nil
}

// AuditChainVerificationDTO formats the response for GET /api/v1/audit/chain/{tableName}/{recordId}.
type AuditChainVerificationDTO struct {
	TableName    string               `json:"table_name"`
	RecordID     string               `json:"record_id"`
	ChainCount   int                  `json:"chain_count"`
	IsValidChain bool                 `json:"is_valid_chain"`
	Records      []domain.AuditRecord `json:"records"`
}

// ErrorResponse formats standard HTTP error messages.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

