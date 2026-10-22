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
	FullName      string                `json:"full_name"`
	Phone         string                `json:"phone"`
	NIN           string                `json:"nin,omitempty"`
	BuildingClass domain.BuildingClass  `json:"building_class"`
	Latitude      float64               `json:"latitude"`
	Longitude     float64               `json:"longitude"`
	TransformerID string                `json:"transformer_id"`
	MeterSerial   string                `json:"meter_serial,omitempty"`
	TariffBand    s01domain.TariffBand  `json:"tariff_band"`
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
	CustomerID     *string               `json:"customer_id,omitempty"`
	Channel        domain.TicketChannel  `json:"channel"`
	Category       domain.TicketCategory `json:"category"`
	Description    string                `json:"description"`
	AssignedUnitID string                `json:"assigned_unit_id"`
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
	ProofNote string `json:"proof_note"`
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
}

// ScheduleMeterDTO defines payload for POST /api/v1/meters/schedule.
type ScheduleMeterDTO struct {
	CustomerID    string              `json:"customer_id"`
	TicketID      *string             `json:"ticket_id,omitempty"`
	ScheduleType  domain.ScheduleType `json:"schedule_type"`
	ScheduledDate time.Time           `json:"scheduled_date"`
}

// Validate checks meter appointment attributes.
func (d ScheduleMeterDTO) Validate() error {
	if strings.TrimSpace(d.CustomerID) == "" || strings.TrimSpace(string(d.ScheduleType)) == "" {
		return errors.New("transport: customer_id and schedule_type are required")
	}
	return nil
}

// ErrorResponse formats standard HTTP error messages.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

