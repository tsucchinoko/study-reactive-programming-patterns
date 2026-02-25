package domain

import (
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/option"
)

// SelectBestDriver は利用可能なドライバーの中からピックアップ地点に最も近いドライバーを選択する。
// 純粋関数 — 副作用なし。
func SelectBestDriver(drivers []Driver, pickup Location) option.Option[Driver] {
	if len(drivers) == 0 {
		return option.None[Driver]()
	}

	best := drivers[0]
	bestDistance := best.CurrentLocation().DistanceTo(pickup)

	for _, d := range drivers[1:] {
		distance := d.CurrentLocation().DistanceTo(pickup)
		if distance < bestDistance {
			best = d
			bestDistance = distance
		}
	}

	return option.Some(best)
}
