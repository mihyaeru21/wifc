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
		num.AddMut(&Uint320{uint64(i)})
	}

	return num, nil
}

// 末尾の 5 byte 目を見るためのマスク
const mask uint64 = 0x000000ff00000000

func (n *Uint320) IsValid() bool {
	// ここが 0x01 ではない場合は invalid 確定なので hash を見る必要がない
	// 255/256 はこっちを通る
	if n[0]&mask != 0x0000000100000000 {
		return false
	}

	return n.isValidHash()
}

// IsValid のほとんどのケースでは stack 領域にすら hash 計算用のメモリを確保する必要がない
// hash 計算部分を別の関数にしておくことで IsValid の時点でメモリが確保されるのを回避する
func (n *Uint320) isValidHash() bool {
	raw := n.Bytes()
	data := raw[2:36]    // 上位2 byte が無駄に多いのを抜いて 32 byte + 前後の 2 byte
	checksum := raw[36:] // 末尾の 4 byte

	hash := sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])

	return checksum[0] == hash[0] &&
		checksum[1] == hash[1] &&
		checksum[2] == hash[2] &&
		checksum[3] == hash[3]
}

func (n *Uint320) Bytes() [40]byte {
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
func (x *Uint320) AddMut(y *Uint320) {
	var c uint64
	x[0], c = bits.Add64(x[0], y[0], 0)
	x[1], c = bits.Add64(x[1], y[1], c)
	x[2], c = bits.Add64(x[2], y[2], c)
	x[3], c = bits.Add64(x[3], y[3], c)
	x[4], _ = bits.Add64(x[4], y[4], c)
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
// より速度が欲しくなったら assembly 実装をパクってくる
func mul(x, y, c uint64) (z1, z0 uint64) {
	hi, lo := bits.Mul64(x, y)
	lo, cc := bits.Add64(lo, c, 0)
	return hi + cc, lo
}
