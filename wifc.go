package wifc

import (
	"errors"
	"math/bits"
)

func IsValid(key string) bool {
	return true
}

// key 32 bytes + head 0x08 byte + tail 0x01 byte + hash 4 bytes = 38 bytes
type Uint320 [5]uint64

func newEmptyUint320() Uint320 {
	return Uint320{}
}

func BuildUint320(key string) (Uint320, error) {
	num := newEmptyUint320()

	if len(key) != 52 {
		return num, errors.New("Invalid key length.")
	}

	return num, nil
}

func (n Uint320) IsValid() bool {
	return true
}

// 用途的に不要なので桁あふれは呼び出し元に返さない
func (x Uint320) Add(y Uint320) Uint320 {
	var z Uint320
	var c uint64
	for i := 0; i < 5; i++ {
		zi, cc := bits.Add64(x[i], y[i], c)
		z[i] = zi
		c = cc
	}
	return z
}

// 用途的に不要なので桁あふれは呼び出し元に返さない
func (x Uint320) Mul(y uint64) Uint320 {
	var z Uint320
	var c uint64
	for i := 0; i < 5; i++ {
		z1, z0 := mul(x[i], y, z[i])
		lo, cc := bits.Add64(z0, c, 0)
		c, z[i] = cc, lo
		c += z1
	}
	return z
}

// Add/Mul/mul の中身は math/big/arith.go の実装をパクってる
func mul(x, y, c uint64) (z1, z0 uint64) {
	hi, lo := bits.Mul64(x, y)
	lo, cc := bits.Add64(lo, c, 0)
	return hi + cc, lo
}
