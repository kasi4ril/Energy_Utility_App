// Package domain defines core entities, value objects, and audit chain models for
// the S02 Ollama Local AI Guardian bounded context.
package domain

import (
	"time"
)

// ActionType identifies database record mutation types recorded in the immutable audit chain.
type ActionType string

const (
	ActionInsert              ActionType = "INSERT"
	ActionUpdate              ActionType = "UPDATE"
	ActionDelete              ActionType = "DELETE"
	ActionMakerCheckerApproval ActionType = "MAKER_CHECKER_APPROVAL"
)

// ChangeCategory identifies sensitive operational changes requiring Maker-Checker oversight.
type ChangeCategory string

const (
	ChangeTariffBandChange  ChangeCategory = "TARIFF_BAND_CHANGE"
	ChangeDebtAdjustment    ChangeCategory = "DEBT_ADJUSTMENT"
	ChangeMeterSwap         ChangeCategory = "METER_SWAP"
	ChangeStoreStockEdit    ChangeCategory = "STORE_STOCK_EDIT"
	ChangeAllowanceOverride ChangeCategory = "ALLOWANCE_OVERRIDE"
)

// RequestStatus defines the operational lifecycle state of a Maker-Checker request.
type RequestStatus string

const (
	StatusPending   RequestStatus = "PENDING"
	StatusApproved  RequestStatus = "APPROVED"
	StatusRejected  RequestStatus = "REJECTED"
	StatusCancelled RequestStatus = "CANCELLED"
)

// AuditRecord encapsulates a cryptographically chained entry in the immutable audit ledger.
// Formula: current_hash = SHA256(prev_hash + ":" + table_name + ":" + record_id + ":" + action + ":" + new_json + ":" + maker_id + ":" + timestamp)
type AuditRecord struct {
	ID          string     `json:"id"`
	TableName   string     `json:"table_name"`
	RecordID    string     `json:"record_id"`
	ActionType  ActionType `json:"action_type"`
	OldJSON     string     `json:"old_json,omitempty"`
	NewJSON     string     `json:"new_json"`
	MakerID     string     `json:"maker_id"`
	CheckerID   *string    `json:"checker_id,omitempty"`
	PrevHash    string     `json:"prev_hash"`
	CurrentHash string     `json:"current_hash"`
	Timestamp   time.Time  `json:"timestamp"`
}

// MakerCheckerRequest encapsulates a pending or evaluated operational change request.
type MakerCheckerRequest struct {
	ID              string         `json:"id"`
	Module          string         `json:"module"`
	ChangeCategory  ChangeCategory `json:"change_category"`
	TargetRecordID  string         `json:"target_record_id"`
	PayloadJSON     string         `json:"payload_json"`
	Justification   string         `json:"justification"`
	OllamaRiskScore float64        `json:"ollama_risk_score"`
	OllamaReasoning string         `json:"ollama_reasoning"`
	Status          RequestStatus  `json:"status"`
	MakerID         string         `json:"maker_id"`
	CheckerID       *string        `json:"checker_id,omitempty"`
	CheckerNote     string         `json:"checker_note,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// ChatMessage represents a community message or @GridAssistant AI response.
type ChatMessage struct {
	ID                 string    `json:"id"`
	ChannelID          string    `json:"channel_id"`
	SenderID           *string   `json:"sender_id,omitempty"`
	SenderName         string    `json:"sender_name"`
	MessageText        string    `json:"message_text"`
	IsAIResponse       bool      `json:"is_ai_response"`
	ReferencedEntityID *string   `json:"referenced_entity_id,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

// RiskEvaluation captures the evaluation result returned by the Python AI Guardian sidecar.
type RiskEvaluation struct {
	RiskScore float64 `json:"risk_score"`
	Reasoning string  `json:"reasoning"`
}

