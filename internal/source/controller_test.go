package source

import (
	"context"
	"testing"
)

func TestStartPreheatFailureRollsBack(t *testing.T) {
	c := NewController(10)
	c.ConfigurePreheatFailures(1) // first Preheat attempt fails
	if err := c.Start(context.Background(), 5); err == nil {
		t.Fatal("expected first Start to fail")
	}
	if got := c.AdmissionUsed(); got != 0 {
		t.Fatalf("admission not released, used=%d", got)
	}
	// second Start must not be rejected by stale adapter occupancy
	if err := c.Start(context.Background(), 5); err != nil {
		t.Fatalf("second Start rejected: %v", err)
	}
	if got := c.AdmissionUsed(); got != 5 {
		t.Fatalf("admission not acquired, used=%d", got)
	}
	c.Stop()
	if got := c.AdmissionUsed(); got != 0 {
		t.Fatalf("Stop did not release admission, used=%d", got)
	}
	// after Stop the device must be free again
	if err := c.Start(context.Background(), 5); err != nil {
		t.Fatalf("Start after Stop rejected: %v", err)
	}
}
