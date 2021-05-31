package wifc_test

import (
	"math"
	"testing"

	"github.com/mihyaeru21/wifc"
)

func Test_BuildUint320(t *testing.T) {
	_, err := wifc.BuildUint320("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q_")
	if err == nil {
		t.Error("Error must be present.")
	}
	if err.Error() != "Invalid key length." {
		t.Errorf("Error message is invalid: %v", err.Error())
	}

	_, err = wifc.BuildUint320("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33_")
	if err == nil {
		t.Error("Error must be present.")
	}
	if err.Error() != "Invalid character: 95" { // _ は 95
		t.Errorf("Error message is invalid: %v", err.Error())
	}

	_, err = wifc.BuildUint320("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q")
	if err != nil {
		t.Error(err)
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

func Test_Uint320_Mul(t *testing.T) {
	tests := []struct {
		x wifc.Uint320
		y uint64
		z wifc.Uint320
	}{
		{
			wifc.Uint320{1, 0, 0, 0, 0},
			0,
			wifc.Uint320{0, 0, 0, 0, 0},
		},
		{
			wifc.Uint320{1, 2, 3, 4, 5},
			1,
			wifc.Uint320{1, 2, 3, 4, 5},
		},
		{
			wifc.Uint320{1, 2, 3, 4, 5},
			2,
			wifc.Uint320{2, 4, 6, 8, 10},
		},
		{
			wifc.Uint320{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64},
			math.MaxUint64,
			wifc.Uint320{1, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64},
		},
	}

	for _, tt := range tests {
		z := tt.x.Mul(tt.y)
		if z != tt.z {
			t.Errorf("expected: %v, but got %v", tt.z, z)
		}
	}
}

func Test_Uint320_Bytes(t *testing.T) {
	tests := []struct {
		x     wifc.Uint320
		bytes [40]byte
	}{
		{
			wifc.Uint320{1, 2, 3, 4, 5},
			[40]byte{0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			wifc.Uint320{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64},
			[40]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
	}

	for _, tt := range tests {
		bytes := tt.x.Bytes()
		if bytes != tt.bytes {
			t.Errorf("expected: %v, but got %v", tt.bytes, bytes)
		}
	}
}

// func Benchmark_BuildUint320(b *testing.B) {
// 	key := "L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q"

// 	base, _ := wifc.BuildUint320(key)

// 	// base := big.NewInt(0)
// 	// for i := 0; i < 52; i++ {
// 	// 	base.Mul(base, big.NewInt(58))
// 	// 	base.Add(base, big.NewInt(int64(i)))
// 	// }

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		// _, _ = wifc.BuildUint320(key)

// 		a := wifc.Uint320{30}
// 		_ = base.Add(a)

// 		// a := big.NewInt(30)
// 		// _ = base.Add(base, a)
// 	}
// }

// func Benchmark_Uint320_Add(b *testing.B) {
// 	x := wifc.Uint320{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64}
// 	y := wifc.Uint320{1, 1, 1, 1, 1}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = x.Add(y)
// 	}
// }

// func Benchmark_Uint320_Mul(b *testing.B) {
// 	x := wifc.Uint320{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = x.Mul(58)
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
