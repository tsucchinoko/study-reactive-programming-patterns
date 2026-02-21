package domain

import (
	"time"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// MenuItem はメニュー上の一品を表す値オブジェクト。
type MenuItem struct {
	id        types.MenuItemID
	name      string
	category  string
	price     types.Money
	prepTime  time.Duration
	available bool
}

// NewMenuItem は新しい MenuItem 値オブジェクトを作成する。
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

// WithAvailability は在庫状態を更新した新しい MenuItem を返す。
func (m MenuItem) WithAvailability(available bool) MenuItem {
	return MenuItem{
		id: m.id, name: m.name, category: m.category,
		price: m.price, prepTime: m.prepTime, available: available,
	}
}
