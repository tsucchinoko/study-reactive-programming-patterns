package simulator

import (
	deliverydomain "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/domain"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// SeedDrivers はシミュレーション用の10名のドライバーを返す。
// 東京エリアの異なる初期位置を持つ。
func SeedDrivers() []deliverydomain.Driver {
	drivers := []struct {
		name  string
		phone string
		lat   float64
		lng   float64
	}{
		{"田中 太郎", "090-1234-5001", 35.6895, 139.6917}, // 新宿
		{"佐藤 花子", "090-1234-5002", 35.6762, 139.6503}, // 渋谷
		{"鈴木 一郎", "090-1234-5003", 35.6813, 139.7671}, // 日本橋
		{"高橋 美咲", "090-1234-5004", 35.6991, 139.7748}, // 上野
		{"山田 健太", "090-1234-5005", 35.6636, 139.7141}, // 赤坂
		{"中村 翔太", "090-1234-5006", 35.6585, 139.7454}, // 品川
		{"小林 愛", "090-1234-5007", 35.7101, 139.8107},   // 北千住
		{"加藤 大輔", "090-1234-5008", 35.6938, 139.7035},  // 四ツ谷
		{"吉田 さくら", "090-1234-5009", 35.6654, 139.7707}, // 豊洲
		{"渡辺 竜也", "090-1234-5010", 35.7063, 139.7513},  // 西日暮里
	}

	result := make([]deliverydomain.Driver, len(drivers))
	for i, d := range drivers {
		result[i] = deliverydomain.NewDriver(
			types.NewDriverID(),
			d.name,
			d.phone,
			deliverydomain.NewLocation(d.lat, d.lng),
		)
	}
	return result
}
