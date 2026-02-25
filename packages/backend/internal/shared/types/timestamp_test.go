package types

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	before := time.Now()
	ts := Now()
	after := time.Now()

	if ts.Time().Before(before) || ts.Time().After(after) {
		t.Error("Now() should return current time")
	}
}

func TestTimestampFrom(t *testing.T) {
	fixed := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	ts := TimestampFrom(fixed)
	if !ts.Time().Equal(fixed) {
		t.Errorf("Time() = %v, want %v", ts.Time(), fixed)
	}
}

func TestTimestampBefore(t *testing.T) {
	t1 := TimestampFrom(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	t2 := TimestampFrom(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC))

	if !t1.Before(t2) {
		t.Error("t1 should be before t2")
	}
	if t2.Before(t1) {
		t.Error("t2 should not be before t1")
	}
}

func TestTimestampAfter(t *testing.T) {
	t1 := TimestampFrom(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	t2 := TimestampFrom(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC))

	if !t2.After(t1) {
		t.Error("t2 should be after t1")
	}
	if t1.After(t2) {
		t.Error("t1 should not be after t2")
	}
}

func TestTimestampAdd(t *testing.T) {
	base := TimestampFrom(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	shifted := base.Add(24 * time.Hour)

	expected := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	if !shifted.Time().Equal(expected) {
		t.Errorf("Add(24h): got %v, want %v", shifted.Time(), expected)
	}
	// 不変性の確認
	if !base.Time().Equal(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Error("Add should not mutate the original")
	}
}

func TestTimestampSince(t *testing.T) {
	ts := TimestampFrom(time.Now().Add(-1 * time.Second))
	elapsed := ts.Since()
	if elapsed < 1*time.Second {
		t.Errorf("Since() = %v, expected >= 1s", elapsed)
	}
}

func TestTimestampString(t *testing.T) {
	ts := TimestampFrom(time.Date(2025, 3, 15, 10, 30, 0, 0, time.UTC))
	expected := "2025-03-15T10:30:00Z"
	if got := ts.String(); got != expected {
		t.Errorf("String() = %q, want %q", got, expected)
	}
}

func TestTimestampIsZero(t *testing.T) {
	if !TimestampFrom(time.Time{}).IsZero() {
		t.Error("zero time should be IsZero")
	}
	if Now().IsZero() {
		t.Error("Now() should not be IsZero")
	}
}
