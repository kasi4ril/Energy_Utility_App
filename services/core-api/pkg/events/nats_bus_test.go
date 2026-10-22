// Package events_test validates asynchronous pub-sub, wildcard matching, and lifecycle.
package events_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/pkg/events"
)

type samplePayload struct {
	StaffID  string `json:"staff_id"`
	NewUnit  string `json:"new_unit"`
}

func TestMemoryEventBus_PublishSubscribe(t *testing.T) {
	bus := events.NewMemoryEventBus()
	defer bus.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	var received events.Event
	_, err := bus.Subscribe("staff.transferred", func(ctx context.Context, ev events.Event) error {
		received = ev
		wg.Done()
		return nil
	})
	if err != nil {
		t.Fatalf("failed to subscribe: %v", err)
	}

	payload := samplePayload{StaffID: "stf-123", NewUnit: "unit-gen-p1"}
	err = bus.Publish(context.Background(), "staff.transferred", "s01_foundation", payload)
	if err != nil {
		t.Fatalf("failed to publish: %v", err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		if received.Subject != "staff.transferred" {
			t.Fatalf("unexpected subject: %s", received.Subject)
		}
		if received.Producer != "s01_foundation" {
			t.Fatalf("unexpected producer: %s", received.Producer)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for event delivery")
	}
}

func TestMemoryEventBus_WildcardMatching(t *testing.T) {
	bus := events.NewMemoryEventBus()
	defer bus.Close()

	var count int
	var mu sync.Mutex

	_, _ = bus.Subscribe("grid.*", func(ctx context.Context, ev events.Event) error {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	})

	_ = bus.Publish(context.Background(), "grid.asset_registered", "foundation", map[string]string{"asset": "A1"})
	_ = bus.Publish(context.Background(), "grid.feeder_tripped", "dispatch", map[string]string{"feeder": "F2"})
	_ = bus.Publish(context.Background(), "unrelated.topic", "other", map[string]string{"k": "v"})

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 2 {
		t.Fatalf("expected 2 wildcard matches, got %d", count)
	}
}

func TestMemoryEventBus_UnsubscribeAndClose(t *testing.T) {
	bus := events.NewMemoryEventBus()

	sub, err := bus.Subscribe("test.topic", func(ctx context.Context, ev events.Event) error {
		return nil
	})
	if err != nil {
		t.Fatalf("failed to subscribe: %v", err)
	}

	sub.Unsubscribe()

	// Publishing should still succeed without delivering to unsubscribed
	err = bus.Publish(context.Background(), "test.topic", "test", "data")
	if err != nil {
		t.Fatalf("unexpected error publishing after unsubscribe: %v", err)
	}

	_ = bus.Close()

	// Operations on closed bus must fail
	err = bus.Publish(context.Background(), "test.topic", "test", "data")
	if err != events.ErrBusClosed {
		t.Fatalf("expected ErrBusClosed, got %v", err)
	}
}

