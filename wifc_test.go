package wifc_test

import (
	"math"
	"testing"

	"github.com/mihyaeru21/wifc"
)

func Test_BuildFromKey(t *testing.T) {
	_, err := wifc.BuildFromKey("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q_")
	if err == nil {
		t.Error("Error must be present.")
	}
	if err.Error() != "invalid key length" {
		t.Errorf("Error message is invalid: %v", err.Error())
	}

	_, err = wifc.BuildFromKey("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33_")
	if err == nil {
		t.Error("Error must be present.")
	}
	if err.Error() != "invalid character: 95" { // _ は 95
		t.Errorf("Error message is invalid: %v", err.Error())
	}

	_, err = wifc.BuildFromKey("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q")
	if err != nil {
		t.Error(err)
	}
}

func Test_Decimal_AddMut(t *testing.T) {
	tests := []struct {
		x wifc.Decimal
		y wifc.Decimal
		z wifc.Decimal
	}{
		{
			wifc.Decimal{0, 0, 0, 0, 0},
			wifc.Decimal{0, 0, 0, 0, 0},
			wifc.Decimal{0, 0, 0, 0, 0},
		},
		{
			wifc.Decimal{1, 2, 3, 4, 5},
			wifc.Decimal{1, 1, 1, 1, 1},
			wifc.Decimal{2, 3, 4, 5, 6},
		},
		{
			wifc.Decimal{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64},
			wifc.Decimal{2, 0, 0, 0, 0},
			wifc.Decimal{1, 0, 0, 0, 0},
		},
		{
			wifc.Decimal{2, 0, 0, 0, 0},
			wifc.Decimal{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64},
			wifc.Decimal{1, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		x := tt.x // ポインタに対するメソッドを呼び出せるようにするためいったん変数に入れる
		x.AddMut(&tt.y)
		if x != tt.z {
			t.Errorf("expected: %v, but got %v", tt.z, x)
		}
	}
}

func Test_Decimal_Mul(t *testing.T) {
	tests := []struct {
		x wifc.Decimal
		y uint64
		z wifc.Decimal
	}{
		{
			wifc.Decimal{1, 0, 0, 0, 0},
			0,
			wifc.Decimal{0, 0, 0, 0, 0},
		},
		{
			wifc.Decimal{1, 2, 3, 4, 5},
			1,
			wifc.Decimal{1, 2, 3, 4, 5},
		},
		{
			wifc.Decimal{1, 2, 3, 4, 5},
			2,
			wifc.Decimal{2, 4, 6, 8, 10},
		},
		{
			wifc.Decimal{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64},
			math.MaxUint64,
			wifc.Decimal{1, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64},
		},
	}

	for _, tt := range tests {
		z := tt.x.Mul(tt.y)
		if z != tt.z {
			t.Errorf("expected: %v, but got %v", tt.z, z)
		}
	}
}

func Test_Decimal_Bytes(t *testing.T) {
	tests := []struct {
		x     wifc.Decimal
		bytes [38]byte
	}{
		{
			wifc.Decimal{1, 2, 3, 4, 5},
			[38]byte{0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			wifc.Decimal{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64},
			[38]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
	}

	for _, tt := range tests {
		bytes := tt.x.Bytes()
		if bytes != tt.bytes {
			t.Errorf("expected: %v, but got %v", tt.bytes, bytes)
		}
	}
}

func Test_Decimal_IsValid(t *testing.T) {
	valid, _ := wifc.BuildFromKey("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q")
	if !valid.IsValid() {
		t.Error("!!!!!")
	}

	// checksum がダメなケース
	invalid1, _ := wifc.BuildFromKey("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33X")
	if invalid1.IsValid() {
		t.Error("!!!!!")
	}

	// 大事な桁が 0x01 じゃないケース(途中の w を W に変えてある)
	invalid2, _ := wifc.BuildFromKey("L3JLGe5rCiCsWFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q")
	if invalid2.IsValid() {
		t.Error("!!!!!")
	}
}

// func Benchmark_Decimal_AddMut(b *testing.B) {
// 	x := wifc.Decimal{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64}
// 	y := &wifc.Decimal{1, 1, 1, 1, 1}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		x.AddMut(y)
// 	}
// }

// func Benchmark_Decimal_Mul(b *testing.B) {
// 	x := wifc.Decimal{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = x.Mul(58)
// 	}
// }

// func Benchmark_Decimal_Bytes(b *testing.B) {
// 	key := "L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q"
// 	x, _ := wifc.BuildFromKey(key)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = x.Bytes()
// 	}
// }

// func Benchmark_Decimal_IsValid_full(b *testing.B) {
// 	key := "L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q"
// 	x, _ := wifc.BuildFromKey(key)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = x.IsValid()
// 	}
// }

// func Benchmark_Decimal_IsValid_fast(b *testing.B) {
// 	key := "L3JLGe5rCiCsWFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q"
// 	x, _ := wifc.BuildFromKey(key)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = x.IsValid()
// 	}
// }

// func Benchmark_Fast(b *testing.B) {
// 	x, _ := wifc.BuildFromKey("L3JLGe5rCiCswFyUKrLZc38iGunHULPk4aFFuHELHKUunt1Ke33Q")
// 	yy, _ := wifc.BuildFromKey("1111111111111111111111111121111111111111111111111111")
// 	y := &yy

// 	// 1/1 slow
// 	// y, _ := wifc.BuildFromKey("1111111111121111111111111111111111111111111111111111")

// 	// 1/2 slow
// 	// y, _ := wifc.BuildFromKey("1111111111112111111111111111111111111111111111111111")

// 	// 1/4 slow
// 	// y, _ := wifc.BuildFromKey("1111111111111211111111111111111111111111111111111111")

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		// y をループ中で足し続けることで58進数の特定の桁について1ずつ増やしていく操作になる
// 		x.AddMut(y)
// 		_ = x.IsValid()
// 	}
// }

// func Benchmark_Slow(b *testing.B) {
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
// 		key := first + strings.Join(chunks, "")
// 		_, _ = base58check.Decode(key)
// 	}
// }
