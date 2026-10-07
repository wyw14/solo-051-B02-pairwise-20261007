package source

import (
	"context"
	"testing"
)

func TestStartPreheatFailureRollsBackEverything(t *testing.T) {
	c := NewController(10)
	c.ConfigurePreheatFailures(1)

	if err := c.Start(context.Background(), 5); err == nil {
		t.Fatal("expected first Start to fail on preheat")
	}
	if got := c.AdmissionUsed(); got != 0 {
		t.Fatalf("admission not rolled back, used=%d", got)
	}
	if got := c.Power(); got != 0 {
		t.Fatalf("power not rolled back, power=%v", got)
	}
	// Adapter must be released: an immediate retry must not be rejected as busy.
	if err := c.Start(context.Background(), 5); err != nil {
		t.Fatalf("retry after failed Start rejected: %v", err)
	}
	if got := c.AdmissionUsed(); got != 5 {
		t.Fatalf("admission not acquired on retry, used=%d", got)
	}
	c.Stop()
	if got := c.AdmissionUsed(); got != 0 {
		t.Fatalf("admission not released on Stop, used=%d", got)
	}
	// Stop must also release the device.
	if err := c.Start(context.Background(), 5); err != nil {
		t.Fatalf("Start after Stop rejected: %v", err)
	}
}
