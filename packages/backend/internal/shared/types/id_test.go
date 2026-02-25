package types

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewOrderID(t *testing.T) {
	id := NewOrderID()
	if id.UUID() == uuid.Nil {
		t.Error("NewOrderID should generate non-nil UUID")
	}
	if id.String() == "" {
		t.Error("String() should return non-empty string")
	}
}

func TestOrderIDFrom(t *testing.T) {
	raw := uuid.New()
	id := OrderIDFrom(raw)
	if id.UUID() != raw {
		t.Errorf("UUID() = %v, want %v", id.UUID(), raw)
	}
}

func TestParseOrderID(t *testing.T) {
	raw := uuid.New()
	id, err := ParseOrderID(raw.String())
	if err != nil {
		t.Fatalf("ParseOrderID(%q) returned error: %v", raw.String(), err)
	}
	if id.UUID() != raw {
		t.Errorf("UUID() = %v, want %v", id.UUID(), raw)
	}
}

func TestParseOrderID_Invalid(t *testing.T) {
	_, err := ParseOrderID("invalid")
	if err == nil {
		t.Error("ParseOrderID(\"invalid\") should return error")
	}
}

func TestNewCustomerID(t *testing.T) {
	id := NewCustomerID()
	if id.UUID() == uuid.Nil {
		t.Error("NewCustomerID should generate non-nil UUID")
	}
}

func TestCustomerIDFrom(t *testing.T) {
	raw := uuid.New()
	id := CustomerIDFrom(raw)
	if id.UUID() != raw {
		t.Errorf("UUID() = %v, want %v", id.UUID(), raw)
	}
}

func TestParseCustomerID(t *testing.T) {
	raw := uuid.New()
	id, err := ParseCustomerID(raw.String())
	if err != nil {
		t.Fatalf("ParseCustomerID returned error: %v", err)
	}
	if id.UUID() != raw {
		t.Errorf("UUID() = %v, want %v", id.UUID(), raw)
	}
}

func TestParseCustomerID_Invalid(t *testing.T) {
	_, err := ParseCustomerID("not-a-uuid")
	if err == nil {
		t.Error("ParseCustomerID with invalid input should return error")
	}
}

func TestNewRestaurantID(t *testing.T) {
	id := NewRestaurantID()
	if id.UUID() == uuid.Nil {
		t.Error("NewRestaurantID should generate non-nil UUID")
	}
}

func TestRestaurantIDFrom(t *testing.T) {
	raw := uuid.New()
	id := RestaurantIDFrom(raw)
	if id.UUID() != raw {
		t.Errorf("UUID() = %v, want %v", id.UUID(), raw)
	}
}

func TestParseRestaurantID(t *testing.T) {
	raw := uuid.New()
	id, err := ParseRestaurantID(raw.String())
	if err != nil {
		t.Fatalf("ParseRestaurantID returned error: %v", err)
	}
	if id.UUID() != raw {
		t.Errorf("UUID() = %v, want %v", id.UUID(), raw)
	}
}

func TestParseRestaurantID_Invalid(t *testing.T) {
	_, err := ParseRestaurantID("")
	if err == nil {
		t.Error("ParseRestaurantID with empty string should return error")
	}
}

func TestNewMenuID(t *testing.T) {
	id := NewMenuID()
	if id.UUID() == uuid.Nil {
		t.Error("NewMenuID should generate non-nil UUID")
	}
	if id.String() == "" {
		t.Error("String() should return non-empty string")
	}
}

func TestNewMenuItemID(t *testing.T) {
	id := NewMenuItemID()
	if id.UUID() == uuid.Nil {
		t.Error("NewMenuItemID should generate non-nil UUID")
	}
}

func TestMenuItemIDFrom(t *testing.T) {
	raw := uuid.New()
	id := MenuItemIDFrom(raw)
	if id.UUID() != raw {
		t.Errorf("UUID() = %v, want %v", id.UUID(), raw)
	}
}

func TestNewDriverID(t *testing.T) {
	id := NewDriverID()
	if id.UUID() == uuid.Nil {
		t.Error("NewDriverID should generate non-nil UUID")
	}
}

func TestDriverIDFrom(t *testing.T) {
	raw := uuid.New()
	id := DriverIDFrom(raw)
	if id.UUID() != raw {
		t.Errorf("UUID() = %v, want %v", id.UUID(), raw)
	}
}

func TestParseDriverID(t *testing.T) {
	raw := uuid.New()
	id, err := ParseDriverID(raw.String())
	if err != nil {
		t.Fatalf("ParseDriverID returned error: %v", err)
	}
	if id.UUID() != raw {
		t.Errorf("UUID() = %v, want %v", id.UUID(), raw)
	}
}

func TestParseDriverID_Invalid(t *testing.T) {
	_, err := ParseDriverID("xyz")
	if err == nil {
		t.Error("ParseDriverID with invalid input should return error")
	}
}

func TestNewAssignmentID(t *testing.T) {
	id := NewAssignmentID()
	if id.UUID() == uuid.Nil {
		t.Error("NewAssignmentID should generate non-nil UUID")
	}
	if id.String() == "" {
		t.Error("String() should return non-empty string")
	}
}

func TestNewNotificationID(t *testing.T) {
	id := NewNotificationID()
	if id.UUID() == uuid.Nil {
		t.Error("NewNotificationID should generate non-nil UUID")
	}
	if id.String() == "" {
		t.Error("String() should return non-empty string")
	}
}

func TestIDUniqueness(t *testing.T) {
	id1 := NewOrderID()
	id2 := NewOrderID()
	if id1.UUID() == id2.UUID() {
		t.Error("two NewOrderID calls should generate different UUIDs")
	}
}

func TestIDStringRoundTrip(t *testing.T) {
	original := NewOrderID()
	parsed, err := ParseOrderID(original.String())
	if err != nil {
		t.Fatalf("ParseOrderID(original.String()) returned error: %v", err)
	}
	if parsed.UUID() != original.UUID() {
		t.Errorf("round-trip failed: got %v, want %v", parsed.UUID(), original.UUID())
	}
}
