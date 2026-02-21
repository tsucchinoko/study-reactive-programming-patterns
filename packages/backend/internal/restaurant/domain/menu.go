package domain

import (
	"fmt"
	"time"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/fp"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/option"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/result"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// Menu はレストランに属する集約を表す。
type Menu struct {
	id    types.MenuID
	items []MenuItem
}

// NewMenu は新しい Menu を作成する。
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

// FindItem は指定された ID を持つ MenuItem を返す。
func (m Menu) FindItem(itemID types.MenuItemID) option.Option[MenuItem] {
	item, found := fp.Find(m.items, func(mi MenuItem) bool {
		return mi.ID().String() == itemID.String()
	})
	if !found {
		return option.None[MenuItem]()
	}
	return option.Some(item)
}

// AvailableItems は現在注文可能な商品のみを返す。純粋関数。
func (m Menu) AvailableItems() []MenuItem {
	return fp.Filter(m.items, func(mi MenuItem) bool {
		return mi.Available()
	})
}

// EstimatePreparationTime は指定された商品セットの合計調理時間を見積もる。
// 全商品を並行して調理すると想定し、最大の調理時間を返す。純粋関数。
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

// ValidateMenuItems は要求された全商品 ID が存在し、注文可能であることを検証する。純粋関数。
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
