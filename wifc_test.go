package wifc_test

import (
	"math"
	"testing"

	"github.com/mihyaeru21/wifc"
)

func Test_BuildUint320(t *testing.T) {
	_, err := wifc.BuildUint320("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q")
	if err != nil {
		t.Error(err)
	}

	_, err = wifc.BuildUint320("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q_")
	if err == nil {
		t.Error("err to be present.")
	}
}

func Test_Uint320_Add(t *testing.T) {
	tests := []struct {
		x wifc.Uint320
		y wifc.Uint320
		z wifc.Uint320
	}{
		{
			wifc.Uint320{0, 0, 0, 0, 0},
			wifc.Uint320{0, 0, 0, 0, 0},
			wifc.Uint320{0, 0, 0, 0, 0},
		},
		{
			wifc.Uint320{1, 2, 3, 4, 5},
			wifc.Uint320{1, 1, 1, 1, 1},
			wifc.Uint320{2, 3, 4, 5, 6},
		},
		{
			wifc.Uint320{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64},
			wifc.Uint320{2, 0, 0, 0, 0},
			wifc.Uint320{1, 0, 0, 0, 0},
		},
		{
			wifc.Uint320{2, 0, 0, 0, 0},
			wifc.Uint320{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64},
			wifc.Uint320{1, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		z := tt.x.Add(tt.y)
		if z != tt.z {
			t.Errorf("expected: %v, but got %v", tt.z, z)
		}
	}
}

// func Benchmark_Uint320_Add(b *testing.B) {
// 	x := wifc.Uint320{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64}
// 	y := wifc.Uint320{1, 1, 1, 1, 1}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = x.Add(y)
// 	}
// }

// func BenchmarkIsValid(b *testing.B) {
// 	key := "L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q"

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = wifc.IsValid(key)
// 	}
// }

// func BenchmarkBase58CheckValid(b *testing.B) {
// 	key := "L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q"
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_, _ = base58check.Decode(key)
// 	}
// }

// func BenchmarkBase58CheckInvalid(b *testing.B) {
// 	key := "L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33X"
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_, _ = base58check.Decode(key)
// 	}
// }

// func BenchmarkConcatString(b *testing.B) {
// 	first := "AAAA"
// 	chunks := []string{
// 		"BBBBBBBBB",
// 		"CCCCCCCCCC",
// 		"DDDDDDDD",
// 		"EEEEEEEEEEEE",
// 		"FFFFFFFFF",
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = first + strings.Join(chunks, "")
// 	}
// }
