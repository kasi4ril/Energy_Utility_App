// Package service coordinates business workflows, enforces domain rules, manages
// cryptographic integrity checks, and dispatches asynchronous domain events.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

var (
	// ErrInvalidCredentials indicates authentication password or user mismatch.
	ErrInvalidCredentials = errors.New("service: invalid email or password")
	// ErrHMACVerificationFailed indicates tampered row signature detected.
	ErrHMACVerificationFailed = errors.New("service: row HMAC tamper detection triggered")
)

// WorkflowService orchestrates all operations for S01 Foundation.
type WorkflowService struct {
	repo       repository.Repository
	bus        events.EventBus
	hmacSecret string
	jwtSecret  string
}

// NewWorkflowService constructs a new Foundation service orchestrator.
func NewWorkflowService(repo repository.Repository, bus events.EventBus, hmacSecret, jwtSecret string) *WorkflowService {
	return &WorkflowService{
		repo:       repo,
		bus:        bus,
		hmacSecret: hmacSecret,
		jwtSecret:  jwtSecret,
	}
}

// RegisterOrgUnit validates hierarchy and registers a new organizational unit.
func (s *WorkflowService) RegisterOrgUnit(ctx context.Context, code, name string, uType domain.OrgUnitType, parentID *string, gps *domain.GPSPoint) (*domain.OrgUnit, error) {
	var parentType *domain.OrgUnitType
	if parentID != nil && *parentID != "" {
		parent, err := s.repo.GetOrgUnitByID(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("service: parent org unit not found: %w", err)
		}
		parentType = &parent.UnitType
	}

	if err := domain.ValidateOrgUnitHierarchy(uType, parentType); err != nil {
		return nil, err
	}
	if gps != nil {
		if err := domain.ValidateGPSPoint(*gps); err != nil {
			return nil, err
		}
	}

	unit := domain.OrgUnit{
		ID:          crypto.MustUUID(),
		ParentID:    parentID,
		UnitCode:    code,
		Name:        name,
		UnitType:    uType,
		LocationGPS: gps,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateOrgUnit(ctx, unit); err != nil {
		return nil, err
	}
	_ = s.bus.Publish(ctx, "foundation.org_unit.registered", "s01_foundation", unit)
	return &unit, nil
}

// RegisterStaffProfile hashes credentials, signs record HMAC, and creates the employee.
func (s *WorkflowService) RegisterStaffProfile(ctx context.Context, staffNum, name, email, password string, role domain.UserRole, orgUnitID string) (*domain.StaffProfile, error) {
	if _, err := s.repo.GetOrgUnitByID(ctx, orgUnitID); err != nil {
		return nil, fmt.Errorf("service: assigned org unit not found: %w", err)
	}

	pwdHash, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}

	staffID := crypto.MustUUID()
	now := time.Now()
	hmacPayload := fmt.Sprintf("%s:%s:%s:%s", staffID, staffNum, email, role)
	sig, err := crypto.ComputeRowHMAC(s.hmacSecret, hmacPayload)
	if err != nil {
		return nil, err
	}

	staff := domain.StaffProfile{
		ID:                  staffID,
		StaffNumber:         staffNum,
		FullName:            name,
		Email:               email,
		Role:                role,
		OrgUnitID:           orgUnitID,
		CertificationStatus: domain.CertActive,
		ActiveShiftState:    domain.ShiftOffDuty,
		PasswordHash:        pwdHash,
		RowHMACSignature:    sig,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := domain.ValidateStaffProfile(staff); err != nil {
		return nil, err
	}
	if err := s.repo.CreateStaffProfile(ctx, staff); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "foundation.staff.created", "s01_foundation", staff)
	return &staff, nil
}

// AuthenticateStaff verifies credentials and HMAC, returning a signed JWT token.
func (s *WorkflowService) AuthenticateStaff(ctx context.Context, email, password string) (*auth.Claims, string, error) {
	staff, err := s.repo.GetStaffProfileByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if !crypto.VerifyPassword(password, staff.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	// Verify cryptographic row HMAC integrity
	hmacPayload := fmt.Sprintf("%s:%s:%s:%s", staff.ID, staff.StaffNumber, staff.Email, staff.Role)
	if !crypto.VerifyRowHMAC(s.hmacSecret, hmacPayload, staff.RowHMACSignature) {
		return nil, "", ErrHMACVerificationFailed
	}

	claims := auth.Claims{
		StaffID:                    staff.ID,
		StaffNumber:                staff.StaffNumber,
		FullName:                   staff.FullName,
		PrimaryRole:                string(staff.Role),
		OrgUnitID:                  staff.OrgUnitID,
		CanViewExecutiveDashboard: auth.IsExecutiveRole(string(staff.Role)),
	}

	token, err := auth.GenerateToken(claims, s.jwtSecret)
	if err != nil {
		return nil, "", fmt.Errorf("service: token generation failed: %w", err)
	}

	return &claims, token, nil
}

// TransferStaff executes organizational re-assignment with cryptographic audit trail.
func (s *WorkflowService) TransferStaff(ctx context.Context, staffID, newOrgUnitID string, newRole domain.UserRole) (*domain.StaffProfile, error) {
	staff, err := s.repo.GetStaffProfileByID(ctx, staffID)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.GetOrgUnitByID(ctx, newOrgUnitID); err != nil {
		return nil, fmt.Errorf("service: destination org unit not found: %w", err)
	}

	newPayload := fmt.Sprintf("%s:%s:%s:%s", staff.ID, staff.StaffNumber, staff.Email, newRole)
	newHMAC, err := crypto.ComputeRowHMAC(s.hmacSecret, newPayload)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateStaffTransfer(ctx, staffID, newOrgUnitID, newRole, newHMAC); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "staff.transferred", "s01_foundation", map[string]interface{}{
		"staff_id":        staffID,
		"new_org_unit_id": newOrgUnitID,
		"new_role":        newRole,
	})

	return s.repo.GetStaffProfileByID(ctx, staffID)
}

// RegisterGridAsset validates grid topology, PostGIS coordinates, and signs HMAC.
func (s *WorkflowService) RegisterGridAsset(ctx context.Context, code string, aType domain.AssetType, parentAssetID *string, orgUnitID string, capKVA float64, tariff domain.TariffBand, gps domain.GPSPoint) (*domain.GridAsset, error) {
	if _, err := s.repo.GetOrgUnitByID(ctx, orgUnitID); err != nil {
		return nil, fmt.Errorf("service: asset org unit not found: %w", err)
	}

	var parentType *domain.AssetType
	if parentAssetID != nil && *parentAssetID != "" {
		parent, err := s.repo.GetGridAssetByID(ctx, *parentAssetID)
		if err != nil {
			return nil, fmt.Errorf("service: parent grid asset not found: %w", err)
		}
		parentType = &parent.AssetType
	}

	if err := domain.ValidateGridTopology(aType, parentType); err != nil {
		return nil, err
	}
	if err := domain.ValidateGPSPoint(gps); err != nil {
		return nil, err
	}

	assetID := crypto.MustUUID()
	now := time.Now()
	hmacPayload := fmt.Sprintf("%s:%s:%s:%f", assetID, code, aType, capKVA)
	sig, err := crypto.ComputeRowHMAC(s.hmacSecret, hmacPayload)
	if err != nil {
		return nil, err
	}

	asset := domain.GridAsset{
		ID:               assetID,
		ParentAssetID:    parentAssetID,
		OrgUnitID:        orgUnitID,
		AssetCode:        code,
		AssetType:        aType,
		RatedCapacityKVA: capKVA,
		TariffBand:       tariff,
		GPSPoint:         gps,
		OperationalState: domain.StateInService,
		RowHMACSignature: sig,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.CreateGridAsset(ctx, asset); err != nil {
		return nil, err
	}
	_ = s.bus.Publish(ctx, "foundation.grid_asset.registered", "s01_foundation", asset)
	return &asset, nil
}

// AppendWorkNote appends an immutable, SHA-256 chained note to the universal WorkNotes ledger.
func (s *WorkflowService) AppendWorkNote(ctx context.Context, entityType, entityID, authorID string, taggedUnitIDs []string, noteText string) (*domain.WorkNote, error) {
	if _, err := s.repo.GetStaffProfileByID(ctx, authorID); err != nil {
		return nil, fmt.Errorf("service: author staff profile not found: %w", err)
	}

	latest, err := s.repo.GetLatestWorkNote(ctx, entityType, entityID)
	if err != nil {
		return nil, err
	}

	prevHash := crypto.GenesisHash
	if latest != nil {
		prevHash = latest.CurrentHash
	}

	now := time.Now()
	currHash := crypto.ComputeWorkNoteHash(prevHash, entityType, entityID, authorID, noteText, now.Unix())

	note := domain.WorkNote{
		ID:            crypto.MustUUID(),
		EntityType:    entityType,
		EntityID:      entityID,
		AuthorStaffID: authorID,
		TaggedUnitIDs: taggedUnitIDs,
		NoteText:      noteText,
		PrevHash:      prevHash,
		CurrentHash:   currHash,
		CreatedAt:     now,
	}

	if err := domain.ValidateWorkNoteIntegrity(note, prevHash); err != nil {
		return nil, err
	}
	if err := s.repo.CreateWorkNote(ctx, note); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "foundation.work_note.appended", "s01_foundation", note)
	return &note, nil
}

