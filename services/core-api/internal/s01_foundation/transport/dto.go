// Package transport implements the HTTP transport layer, request validation DTOs,
// and REST API endpoints for S01 Foundation.
package transport

import (
	"errors"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// LoginRequest defines credentials for POST /api/v1/auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate checks required login attributes.
func (r LoginRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" || strings.TrimSpace(r.Password) == "" {
		return errors.New("transport: email and password are required")
	}
	return nil
}

// LoginResponse contains the issued JWT token and claims.
type LoginResponse struct {
	Token  string      `json:"token"`
	Claims auth.Claims `json:"claims"`
}

// RegisterOrgUnitRequest defines attributes for POST /api/v1/foundation/org-units.
type RegisterOrgUnitRequest struct {
	UnitCode  string             `json:"unit_code"`
	Name      string             `json:"name"`
	UnitType  domain.OrgUnitType `json:"unit_type"`
	ParentID  *string            `json:"parent_id,omitempty"`
	Latitude  *float64           `json:"latitude,omitempty"`
	Longitude *float64           `json:"longitude,omitempty"`
}

// Validate ensures non-empty names and types.
func (r RegisterOrgUnitRequest) Validate() error {
	if strings.TrimSpace(r.UnitCode) == "" || strings.TrimSpace(r.Name) == "" || strings.TrimSpace(string(r.UnitType)) == "" {
		return errors.New("transport: unit_code, name, and unit_type are required")
	}
	return nil
}

// RegisterStaffRequest defines attributes for POST /api/v1/foundation/staff.
type RegisterStaffRequest struct {
	StaffNumber string          `json:"staff_number"`
	FullName    string          `json:"full_name"`
	Email       string          `json:"email"`
	Password    string          `json:"password"`
	Role        domain.UserRole `json:"role"`
	OrgUnitID   string          `json:"org_unit_id"`
}

// Validate ensures required staff parameters are present.
func (r RegisterStaffRequest) Validate() error {
	if strings.TrimSpace(r.StaffNumber) == "" || strings.TrimSpace(r.FullName) == "" ||
		strings.TrimSpace(r.Email) == "" || strings.TrimSpace(r.Password) == "" ||
		strings.TrimSpace(string(r.Role)) == "" || strings.TrimSpace(r.OrgUnitID) == "" {
		return errors.New("transport: all staff profile fields are mandatory")
	}
	return nil
}

// TransferStaffRequest defines payload for POST /api/v1/foundation/staff/transfer.
type TransferStaffRequest struct {
	StaffID      string          `json:"staff_id"`
	NewOrgUnitID string          `json:"new_org_unit_id"`
	NewRole      domain.UserRole `json:"new_role"`
}

// Validate ensures transfer targets are set.
func (r TransferStaffRequest) Validate() error {
	if strings.TrimSpace(r.StaffID) == "" || strings.TrimSpace(r.NewOrgUnitID) == "" || strings.TrimSpace(string(r.NewRole)) == "" {
		return errors.New("transport: staff_id, new_org_unit_id, and new_role are required")
	}
	return nil
}

// RegisterGridAssetRequest defines attributes for POST /api/v1/foundation/assets.
type RegisterGridAssetRequest struct {
	AssetCode        string            `json:"asset_code"`
	AssetType        domain.AssetType  `json:"asset_type"`
	ParentAssetID    *string           `json:"parent_asset_id,omitempty"`
	OrgUnitID        string            `json:"org_unit_id"`
	RatedCapacityKVA float64           `json:"rated_capacity_kva"`
	TariffBand       domain.TariffBand `json:"tariff_band,omitempty"`
	Latitude         float64           `json:"latitude"`
	Longitude        float64           `json:"longitude"`
}

// Validate ensures asset attributes are populated.
func (r RegisterGridAssetRequest) Validate() error {
	if strings.TrimSpace(r.AssetCode) == "" || strings.TrimSpace(string(r.AssetType)) == "" || strings.TrimSpace(r.OrgUnitID) == "" {
		return errors.New("transport: asset_code, asset_type, and org_unit_id are required")
	}
	return nil
}

// AppendWorkNoteRequest defines payload for POST /api/v1/notes/{entityType}/{entityId}.
type AppendWorkNoteRequest struct {
	AuthorStaffID string   `json:"author_staff_id"`
	TaggedUnitIDs []string `json:"tagged_unit_ids"`
	NoteText      string   `json:"note_text"`
}

// Validate ensures work note text and author are present.
func (r AppendWorkNoteRequest) Validate() error {
	if strings.TrimSpace(r.AuthorStaffID) == "" || strings.TrimSpace(r.NoteText) == "" {
		return errors.New("transport: author_staff_id and note_text are required")
	}
	return nil
}

// ErrorResponse formats standard error outputs.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

