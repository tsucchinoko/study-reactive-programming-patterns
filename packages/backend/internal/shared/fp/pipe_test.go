package fp

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	input := []int{1, 2, 3}
	result := Map(input, func(n int) int { return n * 2 })
	expected := []int{2, 4, 6}
	assertSliceEqual(t, "Map", result, expected)
}

func TestMap_Empty(t *testing.T) {
	result := Map([]int{}, func(n int) int { return n * 2 })
	if len(result) != 0 {
		t.Errorf("Map on empty slice should return empty slice, got len=%d", len(result))
	}
}

func TestMap_TypeConversion(t *testing.T) {
	input := []int{1, 2, 3}
	result := Map(input, func(n int) string { return fmt.Sprintf("%d", n) })
	expected := []string{"1", "2", "3"}
	assertSliceEqual(t, "Map type conversion", result, expected)
}

func TestFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result := Filter(input, func(n int) bool { return n%2 == 0 })
	expected := []int{2, 4}
	assertSliceEqual(t, "Filter", result, expected)
}

func TestFilter_NoMatch(t *testing.T) {
	input := []int{1, 3, 5}
	result := Filter(input, func(n int) bool { return n%2 == 0 })
	if len(result) != 0 {
		t.Errorf("Filter with no matches should return empty slice, got len=%d", len(result))
	}
}

func TestFilter_AllMatch(t *testing.T) {
	input := []int{2, 4, 6}
	result := Filter(input, func(n int) bool { return n%2 == 0 })
	assertSliceEqual(t, "Filter all match", result, input)
}

func TestReduce(t *testing.T) {
	input := []int{1, 2, 3, 4}
	result := Reduce(input, 0, func(acc, n int) int { return acc + n })
	if result != 10 {
		t.Errorf("Reduce sum = %d, want 10", result)
	}
}

func TestReduce_Empty(t *testing.T) {
	result := Reduce([]int{}, 42, func(acc, n int) int { return acc + n })
	if result != 42 {
		t.Errorf("Reduce on empty slice should return initial, got %d", result)
	}
}

func TestFind_Found(t *testing.T) {
	input := []int{1, 2, 3, 4}
	val, found := Find(input, func(n int) bool { return n == 3 })
	if !found {
		t.Error("Find should return found=true")
	}
	if val != 3 {
		t.Errorf("Find value = %d, want 3", val)
	}
}

func TestFind_NotFound(t *testing.T) {
	input := []int{1, 2, 3}
	val, found := Find(input, func(n int) bool { return n == 99 })
	if found {
		t.Error("Find should return found=false when not found")
	}
	if val != 0 {
		t.Errorf("Find should return zero value when not found, got %d", val)
	}
}

func TestAny_True(t *testing.T) {
	input := []int{1, 2, 3}
	if !Any(input, func(n int) bool { return n > 2 }) {
		t.Error("Any should return true when predicate matches")
	}
}

func TestAny_False(t *testing.T) {
	input := []int{1, 2, 3}
	if Any(input, func(n int) bool { return n > 10 }) {
		t.Error("Any should return false when no element matches")
	}
}

func TestAny_Empty(t *testing.T) {
	if Any([]int{}, func(n int) bool { return true }) {
		t.Error("Any on empty slice should return false")
	}
}

func TestAll_True(t *testing.T) {
	input := []int{2, 4, 6}
	if !All(input, func(n int) bool { return n%2 == 0 }) {
		t.Error("All should return true when all elements match")
	}
}

func TestAll_False(t *testing.T) {
	input := []int{2, 3, 4}
	if All(input, func(n int) bool { return n%2 == 0 }) {
		t.Error("All should return false when some element doesn't match")
	}
}

func TestAll_Empty(t *testing.T) {
	if !All([]int{}, func(n int) bool { return false }) {
		t.Error("All on empty slice should return true")
	}
}

func TestFlatMap(t *testing.T) {
	input := []int{1, 2, 3}
	result := FlatMap(input, func(n int) []int { return []int{n, n * 10} })
	expected := []int{1, 10, 2, 20, 3, 30}
	assertSliceEqual(t, "FlatMap", result, expected)
}

func TestFlatMap_Empty(t *testing.T) {
	result := FlatMap([]int{}, func(n int) []int { return []int{n} })
	if len(result) != 0 {
		t.Errorf("FlatMap on empty slice should return empty slice, got len=%d", len(result))
	}
}

func TestGroupBy(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	groups := GroupBy(input, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if len(groups) != 2 {
		t.Fatalf("GroupBy should produce 2 groups, got %d", len(groups))
	}
	assertSliceEqual(t, "GroupBy even", groups["even"], []int{2, 4, 6})
	assertSliceEqual(t, "GroupBy odd", groups["odd"], []int{1, 3, 5})
}

func TestGroupBy_Empty(t *testing.T) {
	groups := GroupBy([]int{}, func(n int) string { return "x" })
	if len(groups) != 0 {
		t.Errorf("GroupBy on empty slice should return empty map, got len=%d", len(groups))
	}
}

func TestPipe2(t *testing.T) {
	double := func(n int) int { return n * 2 }
	toString := func(n int) string { return fmt.Sprintf("%d", n) }
	composed := Pipe2(double, toString)
	if result := composed(5); result != "10" {
		t.Errorf("Pipe2(double, toString)(5) = %q, want \"10\"", result)
	}
}

func TestPipe3(t *testing.T) {
	add1 := func(n int) int { return n + 1 }
	double := func(n int) int { return n * 2 }
	toString := func(n int) string { return fmt.Sprintf("%d", n) }
	composed := Pipe3(add1, double, toString)
	// (5+1)*2 = 12
	if result := composed(5); result != "12" {
		t.Errorf("Pipe3(add1, double, toString)(5) = %q, want \"12\"", result)
	}
}

// assertSliceEqual はスライスの内容が一致するか検証するヘルパー
func assertSliceEqual[T comparable](t *testing.T, name string, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: len = %d, want %d", name, len(got), len(want))
		return
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("%s: [%d] = %v, want %v", name, i, got[i], want[i])
		}
	}
}
