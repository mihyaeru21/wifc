package wifc

import (
	"errors"
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

// 本家は中身の中身は assembly だった
func (x Uint320) Add(y Uint320) Uint320 {
	var prev uint64
	var overflow bool

	for i := 0; i < 5; i++ {
		if overflow {
			prev = x[i]
			x[i]++
			overflow = false
			if prev > x[i] {
				overflow = true
			}
		}

		prev = x[i]
		x[i] += y[i]
		if prev > x[i] {
			overflow = true
		}
	}

	return x
}
