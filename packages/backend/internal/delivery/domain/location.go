package domain

import "math"

// Location は地理的な位置を表す Value Object。
type Location struct {
	lat float64
	lng float64
}

// NewLocation は新しい Location を生成する。
func NewLocation(lat, lng float64) Location {
	return Location{lat: lat, lng: lng}
}

func (l Location) Lat() float64 { return l.lat }
func (l Location) Lng() float64 { return l.lng }

// DistanceTo は2つの Location 間の距離をキロメートルで返す（Haversine 公式）。
// 純粋関数。
func (l Location) DistanceTo(other Location) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := l.lat * math.Pi / 180
	lat2Rad := other.lat * math.Pi / 180
	dLat := (other.lat - l.lat) * math.Pi / 180
	dLng := (other.lng - l.lng) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// LerpTo は現在地から目的地へ t (0.0〜1.0) の割合で補間した位置を返す。
// 純粋関数。
func (l Location) LerpTo(target Location, t float64) Location {
	if t <= 0 {
		return l
	}
	if t >= 1 {
		return target
	}
	return NewLocation(
		l.lat+(target.lat-l.lat)*t,
		l.lng+(target.lng-l.lng)*t,
	)
}
