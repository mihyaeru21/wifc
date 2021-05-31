package wifc

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/bits"
)

const characters = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var characterBytes = []byte(characters)

// key 32 bytes + head 0x08 byte + tail 0x01 byte + hash 4 bytes = 38 bytes
type Uint320 [5]uint64

func BuildUint320(key string) (Uint320, error) {
	num := Uint320{}

	if len(key) != 52 {
		return num, errors.New("Invalid key length.")
	}

	for _, c := range key {
		i := bytes.IndexByte(characterBytes, byte(c))
		if i == -1 {
			return num, fmt.Errorf("Invalid character: %v", c)
		}
		num = num.Mul(58)
		num = num.Add(Uint320{uint64(i)})
	}

	return num, nil
}

func (n Uint320) IsValid() bool {
	raw := n.Bytes()

	// TODO: Bytes() せずに判定したい
	// ここが 0x01 ではない場合は invalid 確定なので hash を見る必要がない
	if raw[35] != 0x01 {
		return false
	}

	data := raw[2:36]    // 上位2 byte が無駄に多いを抜いて 32 byte + 前後の 2 byte
	checksum := raw[36:] // 末尾の 4 byte

	hasher := sha256.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(hash)
	hash = hasher.Sum(nil)

	return checksum[0] == hash[0] &&
		checksum[1] == hash[1] &&
		checksum[2] == hash[2] &&
		checksum[3] == hash[3]
}

func (n Uint320) Bytes() [40]byte {
	buf := [40]byte{}

	i := 40
	for _, d := range n {
		for j := 0; j < 8; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}

	return buf
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
