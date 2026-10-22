// Package events provides an asynchronous event bus adhering to the NATS JetStream
// specification for the NESCO Utility Platform.
// It decouples cross-departmental reactions (e.g. WorkNotes tagging, staff transfer
// notifications, asset registration updates) with thread-safe publish-subscribe semantics.
package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"disco-enterprise-platform/services/core-api/pkg/crypto"
)

var (
	// ErrBusClosed is returned when publishing or subscribing on a terminated event bus.
	ErrBusClosed = errors.New("events: event bus is closed")
	// ErrEmptySubject indicates a topic/subject name was not provided.
	ErrEmptySubject = errors.New("events: subject cannot be empty")
)

// Event encapsulates a discrete system occurrence published across departmental boundaries.
type Event struct {
	ID        string          `json:"id"`
	Subject   string          `json:"subject"`
	Producer  string          `json:"producer"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// HandlerFunc defines the callback signature for asynchronous event subscribers.
type HandlerFunc func(ctx context.Context, event Event) error

// Subscription represents an active topic listener that can be cleanly unsubscribed.
type Subscription interface {
	Unsubscribe()
}

// EventBus defines the contract for dispatching and receiving domain events.
type EventBus interface {
	Publish(ctx context.Context, subject string, producer string, payload interface{}) error
	Subscribe(subject string, handler HandlerFunc) (Subscription, error)
	Close() error
}

type memorySub struct {
	id      string
	subject string
	handler HandlerFunc
	bus     *MemoryEventBus
}

func (s *memorySub) Unsubscribe() {
	s.bus.removeSub(s)
}

// MemoryEventBus provides a thread-safe, non-blocking, in-memory implementation of EventBus
// that mimics NATS pub-sub for testing, edge nodes, and local development.
type MemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]*memorySub
	closed      bool
}

// NewMemoryEventBus constructs a fresh in-memory event bus.
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		subscribers: make(map[string][]*memorySub),
	}
}

// Publish serializes the payload and distributes the event to all matching subscribers.
func (b *MemoryEventBus) Publish(ctx context.Context, subject, producer string, payload interface{}) error {
	if subject == "" {
		return ErrEmptySubject
	}

	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return ErrBusClosed
	}

	dataBytes, err := json.Marshal(payload)
	if err != nil {
		b.mu.RUnlock()
		return fmt.Errorf("events: failed to marshal event data: %w", err)
	}

	event := Event{
		ID:        crypto.MustUUID(),
		Subject:   subject,
		Producer:  producer,
		Timestamp: time.Now().Unix(),
		Data:      dataBytes,
	}

	// Collect matching subscribers (exact match or wildcard suffix ".*")
	var matchedSubs []*memorySub
	for pattern, subs := range b.subscribers {
		if matchSubject(pattern, subject) {
			matchedSubs = append(matchedSubs, subs...)
		}
	}
	b.mu.RUnlock()

	// Dispatch non-blockingly to subscribers
	for _, sub := range matchedSubs {
		targetSub := sub
		go func() {
			_ = targetSub.handler(ctx, event)
		}()
	}

	return nil
}

// Subscribe registers a handler callback for the given subject pattern.
func (b *MemoryEventBus) Subscribe(subject string, handler HandlerFunc) (Subscription, error) {
	if subject == "" {
		return nil, ErrEmptySubject
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil, ErrBusClosed
	}

	sub := &memorySub{
		id:      crypto.MustUUID(),
		subject: subject,
		handler: handler,
		bus:     b,
	}

	b.subscribers[subject] = append(b.subscribers[subject], sub)
	return sub, nil
}

func (b *MemoryEventBus) removeSub(sub *memorySub) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subs := b.subscribers[sub.subject]
	for i, s := range subs {
		if s.id == sub.id {
			b.subscribers[sub.subject] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

// Close gracefully stops the event bus.
func (b *MemoryEventBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closed = true
	b.subscribers = make(map[string][]*memorySub)
	return nil
}

// matchSubject checks if a pattern matches a subject, supporting trailing wildcard `.*` or `*`.
func matchSubject(pattern, subject string) bool {
	if pattern == "*" || pattern == ">" || pattern == subject {
		return true
	}
	if len(pattern) > 2 && pattern[len(pattern)-2:] == ".*" {
		prefix := pattern[:len(pattern)-2]
		return len(subject) >= len(prefix) && subject[:len(prefix)] == prefix
	}
	return false
}

