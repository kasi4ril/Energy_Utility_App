// Package transport implements REST API handlers for S03 Customer Care.
package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// HTTPHandler encapsulates REST endpoints for the Customer Care bounded context.
type HTTPHandler struct {
	svc       *service.WorkflowService
	repo      repository.Repository
	jwtSecret string
}

// NewHTTPHandler constructs a new HTTPHandler instance for S03.
func NewHTTPHandler(svc *service.WorkflowService, repo repository.Repository, jwtSecret string) *HTTPHandler {
	return &HTTPHandler{svc: svc, repo: repo, jwtSecret: jwtSecret}
}

// RegisterRoutes registers S03 Customer Care routes onto the ServeMux.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	authMW := auth.RequireAuth(h.jwtSecret)

	// Customer Management
	mux.Handle("POST /api/v1/customers/onboard", authMW(http.HandlerFunc(h.HandleOnboardCustomer)))
	mux.Handle("GET /api/v1/customers", authMW(http.HandlerFunc(h.HandleListCustomers)))
	mux.Handle("GET /api/v1/customers/{id}", authMW(http.HandlerFunc(h.HandleGetCustomer)))

	// Service Tickets
	mux.Handle("POST /api/v1/tickets", authMW(http.HandlerFunc(h.HandleLogTicket)))
	mux.Handle("GET /api/v1/tickets", authMW(http.HandlerFunc(h.HandleListTickets)))
	mux.Handle("GET /api/v1/tickets/{id}", authMW(http.HandlerFunc(h.HandleGetTicket)))
	mux.Handle("POST /api/v1/tickets/{id}/attend", authMW(http.HandlerFunc(h.HandleAttendTicket)))
	mux.Handle("POST /api/v1/tickets/{id}/confirm-close", authMW(http.HandlerFunc(h.HandleConfirmCloseTicket)))

	// Meter Schedules
	mux.Handle("POST /api/v1/meters/schedule", authMW(http.HandlerFunc(h.HandleScheduleMeter)))
	mux.Handle("GET /api/v1/meters/schedules", authMW(http.HandlerFunc(h.HandleListMeterSchedules)))
}

// HandleOnboardCustomer validates Anti-Ghost GPS proximity and onboards a new customer.
func (h *HTTPHandler) HandleOnboardCustomer(w http.ResponseWriter, r *http.Request) {
	var dto OnboardCustomerDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	gps := s01domain.GPSPoint{Latitude: dto.Latitude, Longitude: dto.Longitude}
	cust, err := h.svc.OnboardCustomer(
		r.Context(), dto.FullName, dto.Phone, dto.NIN, dto.BuildingClass,
		gps, dto.TransformerID, dto.MeterSerial, dto.TariffBand,
	)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "onboarding_rejected", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, cust)
}

// HandleListCustomers retrieves customers optionally filtered by DT.
func (h *HTTPHandler) HandleListCustomers(w http.ResponseWriter, r *http.Request) {
	txID := r.URL.Query().Get("transformer_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	customers, err := h.repo.ListCustomers(r.Context(), txID, limit)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, customers)
}

// HandleGetCustomer retrieves customer by UUID.
func (h *HTTPHandler) HandleGetCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cust, err := h.repo.GetCustomerByID(r.Context(), id)
	if err != nil {
		h.jsonError(w, http.StatusNotFound, "not_found", "customer not found")
		return
	}
	h.jsonResponse(w, http.StatusOK, cust)
}

// HandleLogTicket routes an incoming complaint directly to the designated Unit Head.
func (h *HTTPHandler) HandleLogTicket(w http.ResponseWriter, r *http.Request) {
	var dto LogTicketDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	ticket, err := h.svc.LogServiceTicket(r.Context(), dto.CustomerID, dto.Channel, dto.Category, dto.Description, dto.AssignedUnitID)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "ticket_creation_failed", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, ticket)
}

// HandleListTickets retrieves tickets filtered by status, assigned unit, and category.
func (h *HTTPHandler) HandleListTickets(w http.ResponseWriter, r *http.Request) {
	status := domain.TicketStatus(r.URL.Query().Get("status"))
	unitID := r.URL.Query().Get("assigned_unit_id")
	cat := domain.TicketCategory(r.URL.Query().Get("category"))

	tickets, err := h.repo.ListTickets(r.Context(), status, unitID, cat)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, tickets)
}

// HandleGetTicket retrieves ticket details by UUID.
func (h *HTTPHandler) HandleGetTicket(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ticket, err := h.repo.GetTicketByID(r.Context(), id)
	if err != nil {
		h.jsonError(w, http.StatusNotFound, "not_found", "service ticket not found")
		return
	}
	h.jsonResponse(w, http.StatusOK, ticket)
}

// HandleAttendTicket marks ticket ATTENDED_TO with mandatory technical proof note.
func (h *HTTPHandler) HandleAttendTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing staff claims")
		return
	}

	ticketID := r.PathValue("id")
	var dto AttendTicketDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	attended, err := h.svc.AttendToTicket(r.Context(), ticketID, dto.ProofNote, claims.StaffID)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "attend_failed", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, attended)
}

// HandleConfirmCloseTicket enforces Customer Care exclusive confirmation closure gate.
func (h *HTTPHandler) HandleConfirmCloseTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing staff claims")
		return
	}

	ticketID := r.PathValue("id")
	var dto ConfirmClosureDTO
	_ = json.NewDecoder(r.Body).Decode(&dto)

	closedTicket, err := h.svc.ConfirmAndCloseTicket(
		r.Context(), ticketID, claims.StaffID, claims.PrimaryRole, dto.ConfirmedByCustomer, dto.FeedbackNotes,
	)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "closure_rejected", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, closedTicket)
}

// HandleScheduleMeter schedules a metering activity.
func (h *HTTPHandler) HandleScheduleMeter(w http.ResponseWriter, r *http.Request) {
	var dto ScheduleMeterDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	sched, err := h.svc.ScheduleMeterActivity(r.Context(), dto.CustomerID, dto.TicketID, dto.ScheduleType, dto.ScheduledDate)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "schedule_failed", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, sched)
}

// HandleListMeterSchedules lists meter appointments.
func (h *HTTPHandler) HandleListMeterSchedules(w http.ResponseWriter, r *http.Request) {
	status := domain.ScheduleStatus(r.URL.Query().Get("status"))
	schedules, err := h.repo.ListMeterSchedules(r.Context(), status)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, schedules)
}

func (h *HTTPHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *HTTPHandler) jsonError(w http.ResponseWriter, status int, errCode, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: strings.ToLower(errCode), Message: msg})
}

