package client

import (
	"context"
	"testing"
	"time"
)

// TestDefaultWaiterTimeout verifies that DefaultWaiterTimeout is set to 30 minutes.
func TestDefaultWaiterTimeout(t *testing.T) {
	expected := 30 * time.Minute
	if DefaultWaiterTimeout != expected {
		t.Errorf("DefaultWaiterTimeout = %v, want %v", DefaultWaiterTimeout, expected)
	}
}

// TestWaitForResourceDeleted_Maps404ToDeleted verifies that WaitForResourceDeleted treats a 404 error
// as a successful deletion (maps it to StateDeleted).
func TestWaitForResourceDeleted_Maps404ToDeleted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	callCount := 0
	refreshFunc := func() (interface{}, string, error) {
		callCount++
		// Simulate the resource being gone (404) on the first call
		return struct{}{}, "", errorString("404 Not Found")
	}

	err := WaitForResourceDeleted(ctx, refreshFunc)
	if err != nil {
		t.Fatalf("WaitForResourceDeleted should succeed when 404 is returned, got: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

// TestWaitForResource_Maps404ToDeleted verifies that WaitForResource maps 404 to StateDeleted
// when StateDeleted is in the target states.
func TestWaitForResource_Maps404ToDeleted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	callCount := 0
	refreshFunc := func() (interface{}, string, error) {
		callCount++
		return struct{}{}, "", errorString("404 Resource not found")
	}

	err := WaitForResource(ctx, []string{StateDeleting}, []string{StateDeleted}, refreshFunc)
	if err != nil {
		t.Fatalf("WaitForResource should succeed when 404 is returned and StateDeleted is a target state, got: %v", err)
	}
}

// TestWaitForResource_DoesNotMap404WhenNotDeleting verifies that 404 is NOT mapped to StateDeleted
// when StateDeleted is not in the target states — the error should propagate.
func TestWaitForResource_DoesNotMap404WhenNotDeleting(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	refreshFunc := func() (interface{}, string, error) {
		return struct{}{}, "", errorString("404 Not Found")
	}

	// StateDeleted is NOT in targetStates, so 404 should propagate as error
	err := WaitForResource(ctx, []string{StateCreating}, []string{StateActive}, refreshFunc)
	if err == nil {
		t.Fatal("WaitForResource should fail when 404 is returned and StateDeleted is not a target state")
	}
}

// TestWaitForStatus_CompletesWhenStateTransitions verifies that WaitForStatus succeeds when the
// refresh function transitions from a pending state to a target state. Uses short explicit delays.
func TestWaitForStatus_CompletesWhenStateTransitions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	callCount := 0
	refreshFunc := func() (interface{}, string, error) {
		callCount++
		if callCount == 1 {
			return struct{}{}, StateCreating, nil
		}
		return struct{}{}, StateActive, nil
	}

	// Use a short explicit timeout — verifies the timeout parameter works
	err := WaitForStatus(ctx, nil, []string{StateCreating}, []string{StateActive}, refreshFunc, 2*time.Second, 100*time.Millisecond, 100*time.Millisecond, 3)
	if err != nil {
		t.Fatalf("WaitForStatus should succeed when state transitions to target: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

// TestWaitForStatus_TimesOut verifies that WaitForStatus fails when the refresh function never
// transitions to a target state. Uses a short explicit timeout to avoid waiting 30 minutes.
func TestWaitForStatus_TimesOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	refreshFunc := func() (interface{}, string, error) {
		return struct{}{}, StateCreating, nil // stuck in pending state
	}

	// Short timeout to verify timeout behavior without waiting 30 min
	err := WaitForStatus(ctx, nil, []string{StateCreating}, []string{StateActive}, refreshFunc, 1*time.Second, 100*time.Millisecond, 100*time.Millisecond, 3)
	if err == nil {
		t.Fatal("WaitForStatus should time out when state never transitions")
	}
}

// TestWaitForResourceUpdated_UnexpectedStateFailsImmediately verifies that the waiter fails on the first call
// when refreshFunc returns a state not in pending or target (e.g., ERROR instead of EDITING).
func TestWaitForResourceUpdated_UnexpectedStateFailsImmediately(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	callCount := 0
	refreshFunc := func() (interface{}, string, error) {
		callCount++
		// Returns ERROR — not in pending (EDITING) and not in target (ACTIVE)
		return struct{}{}, StateError, nil
	}

	err := WaitForResourceUpdated(ctx, refreshFunc)
	if err == nil {
		t.Fatal("WaitForResourceUpdated should fail when refreshFunc returns an unexpected state (ERROR)")
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (immediate failure), got %d", callCount)
	}
}

// TestWaitForResourceUpdated_UnexpectedStateAfterPending verifies that the waiter fails after cycling through
// pending states before hitting an unexpected state (EDITING → EDITING → ERROR). Uses default delays.
func TestWaitForResourceUpdated_UnexpectedStateAfterPending(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	callCount := 0
	refreshFunc := func() (interface{}, string, error) {
		callCount++
		// First 2 calls return pending state, 3rd call returns ERROR
		if callCount <= 2 {
			return struct{}{}, StateEditing, nil
		}
		return struct{}{}, StateError, nil
	}

	err := WaitForResourceUpdated(ctx, refreshFunc)
	if err == nil {
		t.Fatal("WaitForResourceUpdated should fail when state transitions from pending to ERROR")
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls (2 pending + 1 error), got %d", callCount)
	}
}

// TestWaitForStatus_MultipleTargetStates verifies that WaitForStatus succeeds when the refresh function
// transitions to any one of multiple target states (e.g., ACTIVE or EDITING). Uses short explicit delays.
func TestWaitForStatus_MultipleTargetStates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	callCount := 0
	refreshFunc := func() (interface{}, string, error) {
		callCount++
		if callCount == 1 {
			return struct{}{}, StateCreating, nil
		}
		// Second call hits the second target state — should succeed
		return struct{}{}, StateEditing, nil
	}

	err := WaitForStatus(ctx, nil, []string{StateCreating}, []string{StateActive, StateEditing}, refreshFunc, 2*time.Second, 100*time.Millisecond, 100*time.Millisecond, 3)
	if err != nil {
		t.Fatalf("WaitForStatus should succeed when state transitions to one of multiple target states: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

// TestWaitForStatus_ContextCancellation verifies that the waiter fails fast when the context is
// cancelled, regardless of the configured timeout. Uses short explicit delays.
func TestWaitForStatus_ContextCancellation(t *testing.T) {
	// Cancel the context immediately — should fail fast regardless of the 30-min default
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	refreshFunc := func() (interface{}, string, error) {
		return struct{}{}, StateCreating, nil
	}

	err := WaitForStatus(ctx, nil, []string{StateCreating}, []string{StateActive}, refreshFunc, 2*time.Second, 100*time.Millisecond, 100*time.Millisecond, 3)
	if err == nil {
		t.Fatal("WaitForStatus should fail when context is cancelled")
	}
}

// TestWaitForResourceCreated verifies that WaitForResourceCreated transitions from CREATING to ACTIVE.
// Uses default delays (~5s between polls, growing to 10s).
func TestWaitForResourceCreated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	callCount := 0
	refreshFunc := func() (interface{}, string, error) {
		callCount++
		if callCount == 1 {
			return struct{}{}, StateCreating, nil
		}
		return struct{}{}, StateActive, nil
	}

	err := WaitForResourceCreated(ctx, refreshFunc)
	if err != nil {
		t.Fatalf("WaitForResourceCreated should succeed: %v", err)
	}
}

// TestWaitForResourceUpdated verifies that WaitForResourceUpdated transitions from EDITING to ACTIVE
// after multiple pending polls. Uses default delays (~5s between polls, growing to 10s).
func TestWaitForResourceUpdated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	callCount := 0
	refreshFunc := func() (interface{}, string, error) {
		callCount++
		if callCount <= 2 {
			return struct{}{}, StateEditing, nil
		}
		return struct{}{}, StateActive, nil
	}

	err := WaitForResourceUpdated(ctx, refreshFunc)
	if err != nil {
		t.Fatalf("WaitForResourceUpdated should succeed: %v", err)
	}
}

// errorString is a minimal error implementation for testing.
type errorString string

func (e errorString) Error() string { return string(e) }
