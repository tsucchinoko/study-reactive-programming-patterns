package domain

import (
	"fmt"
	"time"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/fp"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/option"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/result"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// Menu is an aggregate that belongs to a Restaurant.
type Menu struct {
	id    types.MenuID
	items []MenuItem
}

// NewMenu creates a new Menu.
func NewMenu(id types.MenuID, items []MenuItem) Menu {
	copied := make([]MenuItem, len(items))
	copy(copied, items)
	return Menu{id: id, items: copied}
}

func (m Menu) ID() types.MenuID { return m.id }
func (m Menu) Items() []MenuItem {
	copied := make([]MenuItem, len(m.items))
	copy(copied, m.items)
	return copied
}

// FindItem returns the MenuItem with the given ID.
func (m Menu) FindItem(itemID types.MenuItemID) option.Option[MenuItem] {
	item, found := fp.Find(m.items, func(mi MenuItem) bool {
		return mi.ID().String() == itemID.String()
	})
	if !found {
		return option.None[MenuItem]()
	}
	return option.Some(item)
}

// AvailableItems returns only items that are currently available. Pure function.
func (m Menu) AvailableItems() []MenuItem {
	return fp.Filter(m.items, func(mi MenuItem) bool {
		return mi.Available()
	})
}

// EstimatePreparationTime estimates total prep time for a set of items.
// Uses the maximum prep time among all items (parallel preparation).
// Pure function.
func EstimatePreparationTime(items []MenuItem) time.Duration {
	if len(items) == 0 {
		return 0
	}
	return fp.Reduce(items, time.Duration(0), func(max time.Duration, item MenuItem) time.Duration {
		if item.PrepTime() > max {
			return item.PrepTime()
		}
		return max
	})
}

// ValidateMenuItems checks that all requested item IDs exist and are available.
// Pure function.
func ValidateMenuItems(menu Menu, itemIDs []types.MenuItemID) result.Result[[]MenuItem] {
	var found []MenuItem
	for _, id := range itemIDs {
		opt := menu.FindItem(id)
		if opt.IsNone() {
			return result.Err[[]MenuItem](fmt.Errorf("menu item not found: %s", id))
		}
		item := opt.Unwrap()
		if !item.Available() {
			return result.Err[[]MenuItem](fmt.Errorf("menu item unavailable: %s", item.Name()))
		}
		found = append(found, item)
	}
	return result.Ok(found)
}
