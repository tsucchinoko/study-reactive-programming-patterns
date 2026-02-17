package types

import "time"

// Timestamp wraps time.Time as a value object for domain events and entities.
type Timestamp struct {
	value time.Time
}

// Now creates a Timestamp for the current moment.
func Now() Timestamp {
	return Timestamp{value: time.Now()}
}

// TimestampFrom creates a Timestamp from a time.Time.
func TimestampFrom(t time.Time) Timestamp {
	return Timestamp{value: t}
}

// Time returns the underlying time.Time.
func (t Timestamp) Time() time.Time { return t.value }

// Before returns true if this timestamp is before other.
func (t Timestamp) Before(other Timestamp) bool {
	return t.value.Before(other.value)
}

// After returns true if this timestamp is after other.
func (t Timestamp) After(other Timestamp) bool {
	return t.value.After(other.value)
}

// Add returns a new Timestamp shifted by the given duration.
func (t Timestamp) Add(d time.Duration) Timestamp {
	return Timestamp{value: t.value.Add(d)}
}

// Since returns the duration since this timestamp.
func (t Timestamp) Since() time.Duration {
	return time.Since(t.value)
}

// String returns ISO 8601 formatted string.
func (t Timestamp) String() string {
	return t.value.Format(time.RFC3339)
}

// IsZero returns true if the timestamp is the zero value.
func (t Timestamp) IsZero() bool {
	return t.value.IsZero()
}
