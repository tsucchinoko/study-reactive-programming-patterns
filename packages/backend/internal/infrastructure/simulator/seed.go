package simulator

import (
	"time"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/domain"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// SeedRestaurants はシミュレーション用の20件の定義済みレストランとメニューを返す。
func SeedRestaurants() []domain.Restaurant {
	restaurants := []struct {
		name    string
		cuisine domain.CuisineType
		lat     float64
		lng     float64
		address string
		items   []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}
	}{
		{"すし匠", domain.CuisineJapanese, 35.6895, 139.6917, "東京都新宿区西新宿1-1-1", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"特上にぎりセット", "寿司", 2800, 600},
			{"サーモン丼", "丼", 1200, 300},
			{"味噌汁", "汁物", 300, 120},
			{"枝豆", "前菜", 400, 60},
		}},
		{"Trattoria Milano", domain.CuisineItalian, 35.6762, 139.6503, "東京都渋谷区神南1-2-3", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"マルゲリータ", "ピザ", 1500, 480},
			{"カルボナーラ", "パスタ", 1300, 420},
			{"ティラミス", "デザート", 600, 180},
			{"ミネストローネ", "スープ", 700, 300},
		}},
		{"龍華飯店", domain.CuisineChinese, 35.6938, 139.7035, "東京都豊島区東池袋1-4-5", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"麻婆豆腐", "メイン", 900, 360},
			{"餃子（6個）", "点心", 600, 300},
			{"チャーハン", "ご飯もの", 800, 240},
			{"担々麺", "麺", 950, 360},
		}},
		{"Seoul Kitchen", domain.CuisineKorean, 35.7023, 139.7745, "東京都荒川区西日暮里2-6-7", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"サムギョプサル", "焼肉", 1600, 480},
			{"ビビンバ", "ご飯もの", 1100, 300},
			{"トッポギ", "軽食", 700, 240},
			{"キムチチゲ", "鍋", 1000, 420},
		}},
		{"タイ屋台 バンコク", domain.CuisineThai, 35.6614, 139.7039, "東京都目黒区上目黒1-8-9", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"パッタイ", "麺", 1000, 360},
			{"グリーンカレー", "カレー", 1200, 420},
			{"トムヤムクン", "スープ", 900, 300},
			{"ガパオライス", "ご飯もの", 950, 300},
		}},
		{"Mumbai Spice", domain.CuisineIndian, 35.6580, 139.7016, "東京都品川区東五反田2-10-11", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"バターチキンカレー", "カレー", 1100, 420},
			{"ナンセット", "パン", 400, 180},
			{"タンドリーチキン", "グリル", 1300, 480},
			{"ラッシー", "ドリンク", 400, 60},
		}},
		{"NYC Diner", domain.CuisineAmerican, 35.6647, 139.7100, "東京都港区六本木3-12-13", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"クラシックバーガー", "バーガー", 1400, 360},
			{"フライドポテト", "サイド", 500, 240},
			{"クラブサンドイッチ", "サンド", 1200, 300},
			{"チーズケーキ", "デザート", 700, 120},
		}},
		{"Bistro Parisien", domain.CuisineFrench, 35.6700, 139.7696, "東京都中央区銀座4-14-15", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"オニオングラタンスープ", "スープ", 900, 360},
			{"鴨のコンフィ", "メイン", 2200, 600},
			{"クレームブリュレ", "デザート", 800, 240},
			{"ニース風サラダ", "前菜", 1000, 180},
		}},
		{"Taqueria Sol", domain.CuisineMexican, 35.6590, 139.6988, "東京都渋谷区恵比寿南1-16-17", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"タコス（3個）", "タコス", 900, 300},
			{"ブリトー", "メイン", 1100, 360},
			{"ワカモレ＆チップス", "前菜", 600, 120},
			{"チュロス", "デザート", 500, 180},
		}},
		{"天ぷら 天よし", domain.CuisineJapanese, 35.6813, 139.7671, "東京都中央区日本橋2-18-19", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"天ぷら定食", "定食", 1500, 480},
			{"海老天丼", "丼", 1300, 420},
			{"かき揚げ", "一品", 700, 300},
			{"茶碗蒸し", "一品", 500, 240},
		}},
		{"Pasta House NAPOLI", domain.CuisineItalian, 35.6890, 139.7020, "東京都新宿区歌舞伎町1-20-21", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"ペスカトーレ", "パスタ", 1400, 420},
			{"四種のチーズピザ", "ピザ", 1600, 480},
			{"カプレーゼ", "前菜", 800, 120},
			{"パンナコッタ", "デザート", 550, 180},
		}},
		{"焼肉 炎", domain.CuisineKorean, 35.6707, 139.7037, "東京都渋谷区道玄坂2-22-23", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"カルビ定食", "定食", 1800, 420},
			{"ホルモン盛り合わせ", "焼肉", 1400, 360},
			{"冷麺", "麺", 900, 300},
			{"ナムル盛り合わせ", "前菜", 600, 120},
		}},
		{"上海点心楼", domain.CuisineChinese, 35.6715, 139.7651, "東京都中央区築地4-24-25", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"小籠包（6個）", "点心", 800, 360},
			{"エビチリ", "メイン", 1200, 420},
			{"酸辣湯", "スープ", 700, 240},
			{"杏仁豆腐", "デザート", 400, 120},
		}},
		{"Burger Lab", domain.CuisineAmerican, 35.6610, 139.6680, "東京都世田谷区三軒茶屋1-26-27", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"ダブルチーズバーガー", "バーガー", 1600, 420},
			{"オニオンリング", "サイド", 500, 240},
			{"BBQチキンウイング", "サイド", 800, 360},
			{"シェイク", "ドリンク", 600, 120},
		}},
		{"ラーメン一番", domain.CuisineJapanese, 35.6950, 139.7037, "東京都豊島区西池袋3-28-29", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"豚骨ラーメン", "ラーメン", 900, 300},
			{"味玉トッピング", "トッピング", 150, 30},
			{"チャーシュー丼", "丼", 600, 240},
			{"餃子（5個）", "サイド", 450, 240},
		}},
		{"Curry House Bengal", domain.CuisineIndian, 35.6840, 139.7590, "東京都千代田区神田神保町1-30-31", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"チキンビリヤニ", "ご飯もの", 1300, 480},
			{"キーマカレー", "カレー", 1000, 360},
			{"サモサ（2個）", "前菜", 500, 180},
			{"チャイ", "ドリンク", 350, 60},
		}},
		{"café de Lyon", domain.CuisineFrench, 35.6636, 139.7141, "東京都港区赤坂5-32-33", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"キッシュロレーヌ", "メイン", 1100, 360},
			{"クロックムッシュ", "軽食", 900, 240},
			{"ガトーショコラ", "デザート", 700, 180},
			{"カフェオレ", "ドリンク", 500, 60},
		}},
		{"Thai Smile", domain.CuisineThai, 35.6991, 139.7748, "東京都台東区上野4-34-35", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"カオマンガイ", "ご飯もの", 950, 300},
			{"ソムタム", "サラダ", 700, 180},
			{"マッサマンカレー", "カレー", 1100, 420},
			{"マンゴーもち米", "デザート", 600, 180},
		}},
		{"うどん処 讃岐", domain.CuisineJapanese, 35.6817, 139.6736, "東京都渋谷区代々木2-36-37", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"かけうどん", "うどん", 500, 180},
			{"天ぷらうどん", "うどん", 900, 300},
			{"肉うどん", "うどん", 850, 300},
			{"おにぎりセット", "サイド", 300, 60},
		}},
		{"El Mexicano", domain.CuisineMexican, 35.6758, 139.6989, "東京都渋谷区神泉町2-38-39", []struct {
			name     string
			category string
			price    int64
			prepSec  int
		}{
			{"エンチラーダ", "メイン", 1200, 420},
			{"ケサディーヤ", "軽食", 800, 240},
			{"メキシカンサラダ", "サラダ", 700, 120},
			{"フラン", "デザート", 500, 180},
		}},
	}

	result := make([]domain.Restaurant, len(restaurants))
	for i, r := range restaurants {
		items := make([]domain.MenuItem, len(r.items))
		for j, item := range r.items {
			items[j] = domain.NewMenuItem(
				types.NewMenuItemID(),
				item.name,
				item.category,
				types.JPY(item.price),
				time.Duration(item.prepSec)*time.Second,
				true,
			)
		}
		menu := domain.NewMenu(types.NewMenuID(), items)
		location := domain.NewLocation(r.lat, r.lng, r.address)
		result[i] = domain.NewRestaurant(types.NewRestaurantID(), r.name, r.cuisine, location, menu, true)
	}

	return result
}
