package types

import "time"

// Timestamp はドメインイベントとエンティティのために time.Time を値オブジェクトとしてラップする。
type Timestamp struct {
	value time.Time
}

// Now は現在時刻の Timestamp を生成する。
func Now() Timestamp {
	return Timestamp{value: time.Now()}
}

// TimestampFrom は time.Time から Timestamp を生成する。
func TimestampFrom(t time.Time) Timestamp {
	return Timestamp{value: t}
}

// Time は内部の time.Time を返す。
func (t Timestamp) Time() time.Time { return t.value }

// Before はこのタイムスタンプが other より前の場合に true を返す。
func (t Timestamp) Before(other Timestamp) bool {
	return t.value.Before(other.value)
}

// After はこのタイムスタンプが other より後の場合に true を返す。
func (t Timestamp) After(other Timestamp) bool {
	return t.value.After(other.value)
}

// Add は指定した duration だけシフトした新しい Timestamp を返す。
func (t Timestamp) Add(d time.Duration) Timestamp {
	return Timestamp{value: t.value.Add(d)}
}

// Since はこのタイムスタンプからの経過時間を返す。
func (t Timestamp) Since() time.Duration {
	return time.Since(t.value)
}

// String は ISO 8601 形式の文字列を返す。
func (t Timestamp) String() string {
	return t.value.Format(time.RFC3339)
}

// IsZero はタイムスタンプがゼロ値の場合に true を返す。
func (t Timestamp) IsZero() bool {
	return t.value.IsZero()
}
