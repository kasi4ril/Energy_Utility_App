// Package transport provides request/response data transfer objects and validation
// for the S03 Customer Care bounded context.
package transport

import (
	"errors"
	"strings"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
)

// OnboardCustomerDTO defines payload for POST /api/v1/customers/onboard.
type OnboardCustomerDTO struct {
	FullName      string               `json:"full_name"`
	Phone         string               `json:"phone"`
	NIN           string               `json:"nin,omitempty"`
	BuildingClass domain.BuildingClass `json:"building_class"`
	TariffClass   string               `json:"tariff_class,omitempty"`
	EstimatedLoad float64              `json:"estimated_load_kva,omitempty"`
	Latitude      float64              `json:"latitude"`
	Longitude     float64              `json:"longitude"`
	TransformerID string               `json:"transformer_id"`
	FeederID      string               `json:"feeder_id,omitempty"`
	MeterSerial   string               `json:"meter_serial,omitempty"`
	MAPSchemeType string               `json:"map_scheme_type,omitempty"`
	TariffBand    s01domain.TariffBand `json:"tariff_band"`
	StaffID       string               `json:"staff_id,omitempty"`
	StaffName     string               `json:"staff_name,omitempty"`
}

// Validate checks required customer onboarding fields.
func (d OnboardCustomerDTO) Validate() error {
	if strings.TrimSpace(d.FullName) == "" || strings.TrimSpace(d.Phone) == "" ||
		strings.TrimSpace(d.TransformerID) == "" || strings.TrimSpace(string(d.TariffBand)) == "" {
		return errors.New("transport: full_name, phone, transformer_id, and tariff_band are mandatory")
	}
	return nil
}

// LogTicketDTO defines payload for POST /api/v1/tickets.
type LogTicketDTO struct {
	CustomerID        *string               `json:"customer_id,omitempty"`
	Channel           domain.TicketChannel  `json:"channel"`
	Category          domain.TicketCategory `json:"category"`
	SubCategory       string                `json:"sub_category,omitempty"`
	IsSafetyEmergency bool                  `json:"is_safety_emergency"`
	Description       string                `json:"description"`
	AssignedUnitID    string                `json:"assigned_unit_id"`
	AssignedUnitName  string                `json:"assigned_unit_name,omitempty"`
	StaffID           string                `json:"staff_id,omitempty"`
	StaffName         string                `json:"staff_name,omitempty"`
}

// Validate checks ticket creation attributes.
func (d LogTicketDTO) Validate() error {
	if strings.TrimSpace(string(d.Channel)) == "" || strings.TrimSpace(string(d.Category)) == "" ||
		strings.TrimSpace(d.Description) == "" || strings.TrimSpace(d.AssignedUnitID) == "" {
		return errors.New("transport: channel, category, description, and assigned_unit_id are required")
	}
	return nil
}

// AttendTicketDTO defines payload for POST /api/v1/tickets/{id}/attend.
type AttendTicketDTO struct {
	ProofNote       string `json:"proof_note"`
	WorkDoneSummary string `json:"work_done_summary,omitempty"`
	StaffName       string `json:"staff_name,omitempty"`
}

// Validate checks repair proof note.
func (d AttendTicketDTO) Validate() error {
	if strings.TrimSpace(d.ProofNote) == "" {
		return errors.New("transport: proof_note cannot be empty")
	}
	return nil
}

// ConfirmClosureDTO defines payload for POST /api/v1/tickets/{id}/confirm-close.
type ConfirmClosureDTO struct {
	ConfirmedByCustomer bool   `json:"confirmed_by_customer"`
	FeedbackNotes       string `json:"feedback_notes,omitempty"`
	OfficerName         string `json:"officer_name,omitempty"`
}

// ScheduleMeterDTO defines payload for POST /api/v1/meters/schedule.
type ScheduleMeterDTO struct {
	CustomerID    string              `json:"customer_id"`
	TicketID      *string             `json:"ticket_id,omitempty"`
	ScheduleType  domain.ScheduleType `json:"schedule_type"`
	ScheduledDate time.Time           `json:"scheduled_date"`
	StaffID       string              `json:"staff_id,omitempty"`
	StaffName     string              `json:"staff_name,omitempty"`
}

// Validate checks meter appointment attributes.
func (d ScheduleMeterDTO) Validate() error {
	if strings.TrimSpace(d.CustomerID) == "" || strings.TrimSpace(string(d.ScheduleType)) == "" {
		return errors.New("transport: customer_id and schedule_type are required")
	}
	return nil
}

// GenerateKCTDTO defines payload for POST /api/v1/tokens/kct.
type GenerateKCTDTO struct {
	MeterSerial string `json:"meter_serial"`
}

// ErrorResponse formats standard HTTP error messages.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

