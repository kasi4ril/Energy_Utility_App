// Package domain implements business validation rules for Customer Care and Ticketing.
package domain

import (
	"errors"
	"fmt"
	"math"
	"strings"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
)

var (
	// ErrDirectClosureForbidden indicates technical units cannot close a ticket directly.
	ErrDirectClosureForbidden = errors.New("domain: technical units cannot close tickets; mandatory customer confirmation call required")
	// ErrUnauthorizedClosureRole indicates only Customer Care officers can execute the closure gate.
	ErrUnauthorizedClosureRole = errors.New("domain: only Customer Care officers are authorized to confirm customer resolution and close tickets")
	// ErrTicketNotAttended indicates a ticket must be marked ATTENDED_TO before closure confirmation.
	ErrTicketNotAttended = errors.New("domain: ticket must be in ATTENDED_TO state with technical proof before closure confirmation")
	// ErrCustomerConfirmationRequired indicates ticket cannot be closed without telephone confirmation.
	ErrCustomerConfirmationRequired = errors.New("domain: customer telephone call confirmation is strictly mandatory to mark ticket CLOSED")
	// ErrAntiGhostToleranceExceeded indicates customer premises GPS is too far from assigned Distribution Transformer.
	ErrAntiGhostToleranceExceeded = errors.New("domain: Anti-Ghost Lock: premises GPS distance exceeds transformer service radius")
	// ErrDebtorAccountBlocked indicates debtor accounts cannot onboard new service without clearance.
	ErrDebtorAccountBlocked = errors.New("domain: Anti-Ghost Lock: account or NIN is flagged with outstanding debtor balance")
)

// DefaultDTRadiusMeters defines the maximum allowable distance (500 meters) from customer premises to assigned DT.
const DefaultDTRadiusMeters = 500.0

// ResolveDirectRouting determines the designated technical unit head and supervising department head.
func ResolveDirectRouting(category TicketCategory) RoutingTarget {
	switch category {
	case CategoryFallenPole, CategoryLineSnap:
		return RoutingTarget{
			TargetUnitHeadRole:  string(s01domain.RoleUnitHeadFeeder),
			SupervisingDeptRole: string(s01domain.RoleHeadServices),
			SLADurationHours:    4,
		}
	case CategoryTransformerFault, CategoryFuseBlown, CategoryOilLeak:
		return RoutingTarget{
			TargetUnitHeadRole:  string(s01domain.RoleUnitHeadTx),
			SupervisingDeptRole: string(s01domain.RoleHeadServices),
			SLADurationHours:    6,
		}
	case CategoryFaultyMeter, CategoryNewMeterSchedule:
		return RoutingTarget{
			TargetUnitHeadRole:  string(s01domain.RoleUnitHeadMetering),
			SupervisingDeptRole: string(s01domain.RoleHeadServices),
			SLADurationHours:    24,
		}
	case CategoryWrongBilling, CategoryUncreditedPayment:
		return RoutingTarget{
			TargetUnitHeadRole:  string(s01domain.RoleUnitHeadBilling),
			SupervisingDeptRole: string(s01domain.RoleHeadCommercial),
			SLADurationHours:    24,
		}
	case CategorySafetyEmergency:
		return RoutingTarget{
			TargetUnitHeadRole:  string(s01domain.RoleUnitHeadFeeder),
			SupervisingDeptRole: string(s01domain.RoleHeadServices),
			SLADurationHours:    2,
		}
	default:
		return RoutingTarget{
			TargetUnitHeadRole:  string(s01domain.RoleCustomerCareHead),
			SupervisingDeptRole: string(s01domain.RoleCustomerCareHead),
			SLADurationHours:    12,
		}
	}
}

// ResolveDirectRoutingWithEmergency enforces < 2-hour SLA for life-threatening safety emergencies.
func ResolveDirectRoutingWithEmergency(category TicketCategory, isEmergency bool) RoutingTarget {
	target := ResolveDirectRouting(category)
	if isEmergency || category == CategorySafetyEmergency {
		target.SLADurationHours = 2
	}
	return target
}

// ValidateClosureGate enforces the exclusive Customer Care telephone confirmation barrier.
func ValidateClosureGate(ticket ServiceTicket, officerRole string, confirmedByCustomer bool) (TicketStatus, error) {
	// 1. Role Authorization Check: Only Customer Care can close tickets
	if officerRole != string(s01domain.RoleCustomerCareOp) && officerRole != string(s01domain.RoleCustomerCareHead) {
		return "", ErrUnauthorizedClosureRole
	}

	// 2. State Check: Ticket must already be ATTENDED_TO by technical unit
	if ticket.Status != StatusAttendedTo {
		return "", fmt.Errorf("%w: current status is %s", ErrTicketNotAttended, ticket.Status)
	}

	// 3. Customer Denial Branch: If customer denies resolution during confirmation call, escalate!
	if !confirmedByCustomer {
		return StatusReopenedEscalated, nil
	}

	return StatusClosed, nil
}

// ValidateAntiGhostDTDistance calculates spherical distance from customer to DT and enforces radius lock.
func ValidateAntiGhostDTDistance(premisesGPS, dtGPS s01domain.GPSPoint, maxRadiusMeters float64) (float64, error) {
	if maxRadiusMeters <= 0 {
		maxRadiusMeters = DefaultDTRadiusMeters
	}

	dist := calculateHaversineMeters(premisesGPS.Latitude, premisesGPS.Longitude, dtGPS.Latitude, dtGPS.Longitude)
	if dist > maxRadiusMeters {
		return dist, fmt.Errorf("%w: distance %.1fm exceeds allowable radius of %.1fm", ErrAntiGhostToleranceExceeded, dist, maxRadiusMeters)
	}
	return dist, nil
}

// ValidateCustomerOnboarding validates phone, full name, and debtor flags.
func ValidateCustomerOnboarding(c Customer) error {
	if strings.TrimSpace(c.FullName) == "" {
		return errors.New("domain: customer full_name cannot be empty")
	}
	if strings.TrimSpace(c.Phone) == "" {
		return errors.New("domain: customer phone cannot be empty")
	}
	if strings.TrimSpace(c.TransformerID) == "" {
		return errors.New("domain: assigned transformer_id is required")
	}
	if c.KYCStatus == KYCDebtorFlagged || c.LedgerBalance < -100.0 {
		return ErrDebtorAccountBlocked
	}
	return nil
}

func calculateHaversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000.0
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	rLat1 := lat1 * (math.Pi / 180.0)
	rLat2 := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2.0)*math.Sin(dLat/2.0) +
		math.Cos(rLat1)*math.Cos(rLat2)*math.Sin(dLon/2.0)*math.Sin(dLon/2.0)
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))
	return earthRadius * c
}

