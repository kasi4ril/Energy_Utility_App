// Package repository implements the PostgreSQL persistence layer for S02 Ollama Guardian.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/domain"
)

// PostgresRepository provides production PostgreSQL operations for the security_ai schema.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a new Postgres repository handler for security_ai.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateAuditRecord appends a record to security_ai.immutable_audit_chain.
func (r *PostgresRepository) CreateAuditRecord(ctx context.Context, rec domain.AuditRecord) error {
	query := `
		INSERT INTO security_ai.immutable_audit_chain 
		(id, table_name, record_id, action_type, old_json, new_json, maker_id, checker_id, prev_hash, current_hash, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	var oldVal *string
	if rec.OldJSON != "" {
		oldVal = &rec.OldJSON
	}

	_, err := r.db.ExecContext(ctx, query,
		rec.ID, rec.TableName, rec.RecordID, rec.ActionType, oldVal, rec.NewJSON, rec.MakerID, rec.CheckerID, rec.PrevHash, rec.CurrentHash, rec.Timestamp)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateAuditRecord failed: %w", err)
	}
	return nil
}

// GetLatestAuditRecord retrieves the latest audit chain entry for a table record.
func (r *PostgresRepository) GetLatestAuditRecord(ctx context.Context, tableName, recordID string) (*domain.AuditRecord, error) {
	query := `
		SELECT id, table_name, record_id, action_type, old_json, new_json, maker_id, checker_id, prev_hash, current_hash, timestamp
		FROM security_ai.immutable_audit_chain
		WHERE table_name = $1 AND record_id = $2
		ORDER BY timestamp DESC LIMIT 1`

	var rec domain.AuditRecord
	var oldJSON, checkerID sql.NullString
	err := r.db.QueryRowContext(ctx, query, tableName, recordID).Scan(
		&rec.ID, &rec.TableName, &rec.RecordID, &rec.ActionType, &oldJSON, &rec.NewJSON, &rec.MakerID, &checkerID, &rec.PrevHash, &rec.CurrentHash, &rec.Timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No prior audit row
		}
		return nil, fmt.Errorf("postgres_repo: GetLatestAuditRecord failed: %w", err)
	}
	if oldJSON.Valid {
		rec.OldJSON = oldJSON.String
	}
	if checkerID.Valid {
		rec.CheckerID = &checkerID.String
	}
	return &rec, nil
}

// ListAuditRecords retrieves all audit chain rows for a record ordered chronologically.
func (r *PostgresRepository) ListAuditRecords(ctx context.Context, tableName, recordID string) ([]domain.AuditRecord, error) {
	query := `
		SELECT id, table_name, record_id, action_type, old_json, new_json, maker_id, checker_id, prev_hash, current_hash, timestamp
		FROM security_ai.immutable_audit_chain
		WHERE table_name = $1 AND record_id = $2
		ORDER BY timestamp ASC`

	rows, err := r.db.QueryContext(ctx, query, tableName, recordID)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListAuditRecords failed: %w", err)
	}
	defer rows.Close()

	var records []domain.AuditRecord
	for rows.Next() {
		var rec domain.AuditRecord
		var oldJSON, checkerID sql.NullString
		if err := rows.Scan(&rec.ID, &rec.TableName, &rec.RecordID, &rec.ActionType, &oldJSON, &rec.NewJSON, &rec.MakerID, &checkerID, &rec.PrevHash, &rec.CurrentHash, &rec.Timestamp); err != nil {
			return nil, err
		}
		if oldJSON.Valid {
			rec.OldJSON = oldJSON.String
		}
		if checkerID.Valid {
			rec.CheckerID = &checkerID.String
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

// CreateMakerCheckerRequest registers an operational change request.
func (r *PostgresRepository) CreateMakerCheckerRequest(ctx context.Context, req domain.MakerCheckerRequest) error {
	query := `
		INSERT INTO security_ai.maker_checker_requests
		(id, module, change_category, target_record_id, payload_json, justification, ollama_risk_score, ollama_reasoning, status, maker_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		req.ID, req.Module, req.ChangeCategory, req.TargetRecordID, req.PayloadJSON, req.Justification, req.OllamaRiskScore, req.OllamaReasoning, req.Status, req.MakerID, req.CreatedAt, req.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateMakerCheckerRequest failed: %w", err)
	}
	return nil
}

// GetMakerCheckerRequestByID fetches a request by UUID.
func (r *PostgresRepository) GetMakerCheckerRequestByID(ctx context.Context, id string) (*domain.MakerCheckerRequest, error) {
	query := `
		SELECT id, module, change_category, target_record_id, payload_json, justification, ollama_risk_score, ollama_reasoning, status, maker_id, checker_id, checker_note, created_at, updated_at
		FROM security_ai.maker_checker_requests WHERE id = $1`

	var req domain.MakerCheckerRequest
	var checkerID, checkerNote sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&req.ID, &req.Module, &req.ChangeCategory, &req.TargetRecordID, &req.PayloadJSON, &req.Justification, &req.OllamaRiskScore, &req.OllamaReasoning, &req.Status, &req.MakerID, &checkerID, &checkerNote, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetMakerCheckerRequestByID failed: %w", err)
	}
	if checkerID.Valid {
		req.CheckerID = &checkerID.String
	}
	if checkerNote.Valid {
		req.CheckerNote = checkerNote.String
	}
	return &req, nil
}

// UpdateMakerCheckerStatus updates request state to approved, rejected, or cancelled.
func (r *PostgresRepository) UpdateMakerCheckerStatus(ctx context.Context, id string, status domain.RequestStatus, checkerID string, note string) error {
	query := `
		UPDATE security_ai.maker_checker_requests
		SET status = $2, checker_id = $3, checker_note = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id, status, checkerID, note)
	if err != nil {
		return fmt.Errorf("postgres_repo: UpdateMakerCheckerStatus failed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListMakerCheckerRequests returns requests with optional status and module filtering.
func (r *PostgresRepository) ListMakerCheckerRequests(ctx context.Context, status domain.RequestStatus, module string) ([]domain.MakerCheckerRequest, error) {
	query := `
		SELECT id, module, change_category, target_record_id, payload_json, justification, ollama_risk_score, ollama_reasoning, status, maker_id, checker_id, checker_note, created_at, updated_at
		FROM security_ai.maker_checker_requests
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR module = $2)
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, string(status), module)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListMakerCheckerRequests failed: %w", err)
	}
	defer rows.Close()

	var requests []domain.MakerCheckerRequest
	for rows.Next() {
		var req domain.MakerCheckerRequest
		var checkerID, checkerNote sql.NullString
		if err := rows.Scan(&req.ID, &req.Module, &req.ChangeCategory, &req.TargetRecordID, &req.PayloadJSON, &req.Justification, &req.OllamaRiskScore, &req.OllamaReasoning, &req.Status, &req.MakerID, &checkerID, &checkerNote, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, err
		}
		if checkerID.Valid {
			req.CheckerID = &checkerID.String
		}
		if checkerNote.Valid {
			req.CheckerNote = checkerNote.String
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

// CreateChatMessage records an in-app community message.
func (r *PostgresRepository) CreateChatMessage(ctx context.Context, msg domain.ChatMessage) error {
	query := `
		INSERT INTO security_ai.community_chat_messages
		(id, channel_id, sender_id, sender_name, message_text, is_ai_response, referenced_entity_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, query,
		msg.ID, msg.ChannelID, msg.SenderID, msg.SenderName, msg.MessageText, msg.IsAIResponse, msg.ReferencedEntityID, msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateChatMessage failed: %w", err)
	}
	return nil
}

// ListChatMessages fetches recent messages for a channel.
func (r *PostgresRepository) ListChatMessages(ctx context.Context, channelID string, limit int) ([]domain.ChatMessage, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, channel_id, sender_id, sender_name, message_text, is_ai_response, referenced_entity_id, created_at
		FROM security_ai.community_chat_messages
		WHERE channel_id = $1
		ORDER BY created_at ASC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, channelID, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListChatMessages failed: %w", err)
	}
	defer rows.Close()

	var messages []domain.ChatMessage
	for rows.Next() {
		var msg domain.ChatMessage
		var senderID, refID sql.NullString
		if err := rows.Scan(&msg.ID, &msg.ChannelID, &senderID, &msg.SenderName, &msg.MessageText, &msg.IsAIResponse, &refID, &msg.CreatedAt); err != nil {
			return nil, err
		}
		if senderID.Valid {
			msg.SenderID = &senderID.String
		}
		if refID.Valid {
			msg.ReferencedEntityID = &refID.String
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

