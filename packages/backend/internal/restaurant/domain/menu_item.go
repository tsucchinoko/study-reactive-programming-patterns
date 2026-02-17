package domain

import (
	"time"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// MenuItem is a value object representing a single dish on a menu.
type MenuItem struct {
	id        types.MenuItemID
	name      string
	category  string
	price     types.Money
	prepTime  time.Duration
	available bool
}

// NewMenuItem creates a new MenuItem value object.
func NewMenuItem(id types.MenuItemID, name, category string, price types.Money, prepTime time.Duration, available bool) MenuItem {
	return MenuItem{
		id:        id,
		name:      name,
		category:  category,
		price:     price,
		prepTime:  prepTime,
		available: available,
	}
}

func (m MenuItem) ID() types.MenuItemID    { return m.id }
func (m MenuItem) Name() string            { return m.name }
func (m MenuItem) Category() string        { return m.category }
func (m MenuItem) Price() types.Money      { return m.price }
func (m MenuItem) PrepTime() time.Duration { return m.prepTime }
func (m MenuItem) Available() bool         { return m.available }

// WithAvailability returns a new MenuItem with updated availability.
func (m MenuItem) WithAvailability(available bool) MenuItem {
	return MenuItem{
		id: m.id, name: m.name, category: m.category,
		price: m.price, prepTime: m.prepTime, available: available,
	}
}
