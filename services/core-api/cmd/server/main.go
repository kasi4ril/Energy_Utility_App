// Package main serves as the primary runtime entrypoint for the NESCO Core API Gateway.
// It initializes connection pools, event buses, foundational bounded contexts,
// registers REST/WebSocket routes, and manages graceful operating system shutdowns.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/service"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/transport"
	s02repo "disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/repository"
	s02svc "disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/service"
	s02trans "disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/transport"
	s03repo "disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	s03svc "disco-enterprise-platform/services/core-api/internal/s03_customer_care/service"
	s03trans "disco-enterprise-platform/services/core-api/internal/s03_customer_care/transport"
	s04repo "disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	s04svc "disco-enterprise-platform/services/core-api/internal/s04_services_maint/service"
	s04trans "disco-enterprise-platform/services/core-api/internal/s04_services_maint/transport"
	s05repo "disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	s05svc "disco-enterprise-platform/services/core-api/internal/s05_generation_master/service"
	s05trans "disco-enterprise-platform/services/core-api/internal/s05_generation_master/transport"
	s06repo "disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	s06svc "disco-enterprise-platform/services/core-api/internal/s06_energy_management/service"
	s06trans "disco-enterprise-platform/services/core-api/internal/s06_energy_management/transport"
	s07repo "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
	s07svc "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/service"
	s07trans "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/transport"
	s08repo "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
	s08svc "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/service"
	s08trans "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/transport"
	s09repo "disco-enterprise-platform/services/core-api/internal/s09_executive_command/repository"
	s09svc "disco-enterprise-platform/services/core-api/internal/s09_executive_command/service"
	s09trans "disco-enterprise-platform/services/core-api/internal/s09_executive_command/transport"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

func main() {
	log.Println("[INFO] Initializing NESCO Enterprise Utility Core API Gateway...")

	port := getEnv("PORT", "8080")
	hmacSecret := getEnv("HMAC_SECRET", "nesco-master-hmac-integrity-key-32b!")
	jwtSecret := getEnv("JWT_SECRET", "nesco-master-jwt-security-key-32b!!")

	// 1. Initialize Event Bus
	eventBus := events.NewMemoryEventBus()
	defer eventBus.Close()
	log.Println("[INFO] Asynchronous Event Bus online.")

	// 2. Initialize Foundation Repository
	// In-memory repo is used by default for resilient sandbox/edge environments
	repo := repository.NewMemoryRepository()
	log.Println("[INFO] Storage repository initialized.")

	// 3. Initialize S01 Foundation Workflow Service
	foundationService := service.NewWorkflowService(repo, eventBus, hmacSecret, jwtSecret)

	// 4. Register HTTP Routes
	mux := http.NewServeMux()

	// Platform Health & Liveness Checks
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy","service":"nesco-core-api","sprint":"s09_executive_command"}`))
	})

	// Mount S01 Foundation Handlers
	foundationHandler := transport.NewHTTPHandler(foundationService, repo, jwtSecret)
	foundationHandler.RegisterRoutes(mux)
	log.Println("[INFO] S01 Foundation routes registered.")

	// Mount S02 Ollama Guardian Handlers
	aiSidecarURL := getEnv("AI_SIDECAR_URL", "http://127.0.0.1:8001")
	s02Repo := s02repo.NewMemoryRepository()
	s02Service := s02svc.NewWorkflowService(s02Repo, eventBus, aiSidecarURL)
	s02Handler := s02trans.NewHTTPHandler(s02Service, s02Repo, jwtSecret)
	s02Handler.RegisterRoutes(mux)
	log.Println("[INFO] S02 Ollama AI Guardian & Maker-Checker routes registered.")

	// Mount S03 Customer Care Handlers
	s03Repo := s03repo.NewMemoryRepository()
	s03Service := s03svc.NewWorkflowService(s03Repo, repo, eventBus, hmacSecret)
	s03Handler := s03trans.NewHTTPHandler(s03Service, s03Repo, jwtSecret)
	s03Handler.RegisterRoutes(mux)
	log.Println("[INFO] S03 Customer Care & Ticketing routes registered.")

	// Mount S04 Services & Maintenance Handlers
	s04Repo := s04repo.NewMemoryRepository()
	s04Service := s04svc.NewWorkflowService(s04Repo, repo, s03Repo, eventBus, aiSidecarURL)
	s04Handler := s04trans.NewHTTPHandler(s04Service, s04Repo, jwtSecret)
	s04Handler.RegisterRoutes(mux)
	log.Println("[INFO] S04 Services & Maintenance routes registered.")

	// Mount S05 Generation Department & Master Station Handlers
	s05Repo := s05repo.NewMemoryRepository()
	s05Service := s05svc.NewWorkflowService(s05Repo, repo, eventBus, hmacSecret)
	s05Handler := s05trans.NewHTTPHandler(s05Service, s05Repo, jwtSecret)
	s05Handler.RegisterRoutes(mux)
	log.Println("[INFO] S05 Generation Master & Google Forms Webhook routes registered.")

	// Mount S06 Energy Management, Dispatch & PTW Handlers
	s06Repo := s06repo.NewMemoryRepository()
	s06Service := s06svc.NewWorkflowService(s06Repo, repo, eventBus, hmacSecret)
	s06Handler := s06trans.NewHTTPHandler(s06Service, s06Repo, jwtSecret)
	s06Handler.RegisterRoutes(mux)
	log.Println("[INFO] S06 Energy Management, Dispatch & PTW routes registered.")

	// Mount S07 Stores, Material Requisitions & Audit Handlers
	s07Repo := s07repo.NewMemoryRepository()
	s07Service := s07svc.NewWorkflowService(s07Repo, eventBus, hmacSecret, aiSidecarURL)
	s07Handler := s07trans.NewHTTPHandler(s07Service, s07Repo, jwtSecret)
	s07Handler.RegisterRoutes(mux)
	log.Println("[INFO] S07 Stores, Material Requisitions & Audit routes registered.")

	// Mount S08 Billing, Commercial Operations & Revenue Protection AI Handlers
	s08Repo := s08repo.NewMemoryRepository()
	s08Service := s08svc.NewWorkflowService(s08Repo, eventBus, s03Repo)
	s08Handler := s08trans.NewHTTPHandler(s08Service, s08Repo, jwtSecret)
	s08Handler.RegisterRoutes(mux)
	log.Println("[INFO] S08 Billing, Commercial Operations & Revenue Protection AI routes registered.")

	// Mount S09 Executive Command & Predictive Maintenance Handlers
	s09Repo := s09repo.NewMemoryRepository()
	s09Service := s09svc.NewWorkflowService(s09Repo, eventBus, s04Repo, s05Repo, s06Repo, s07Repo, s08Repo)
	s09Handler := s09trans.NewHTTPHandler(s09Service, s09Repo, jwtSecret)
	s09Handler.RegisterRoutes(mux)
	log.Println("[INFO] S09 Executive Command & Predictive Maintenance routes registered.")

	// 5. Configure HTTP Server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Graceful Shutdown Listener
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[INFO] Server listening on http://0.0.0.0:%s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Server listener error: %v", err)
		}
	}()

	<-stop
	log.Println("[INFO] Termination signal received. Initiating graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR] Graceful shutdown encountered an error: %v", err)
	}
	log.Println("[INFO] NESCO Core API stopped cleanly.")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
