package wifc

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"math/bits"
)

const characters = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var characterBytes = []byte(characters)

// key 32 bytes + head 0x08 byte + tail 0x01 byte + hash 4 bytes = 38 bytes
type Decimal [5]uint64

// from https://github.com/anaskhan96/base58check/blob/master/base58check.go#L114
func BuildFromKey(key string) (Decimal, error) {
	num := Decimal{}

	if len(key) != 52 {
		return num, errors.New("invalid key length")
	}

	for _, c := range key {
		i := bytes.IndexByte(characterBytes, byte(c))
		if i == -1 {
			return num, fmt.Errorf("invalid character: %v", c)
		}
		num = num.Mul(58)
		num.AddMut(&Decimal{uint64(i)})
	}

	return num, nil
}

func (n *Decimal) IsValid() bool {
	// 255/256 はこっちを通る
	if !n.IsValidDigit() {
		return false
	}

	return n.IsValidHash()
}

// 末尾の 5 byte 目を見るためのマスク
const mask uint64 = 0x000000ff00000000

// ここが 0x01 ではない場合は invalid 確定なので hash を見る必要がない
func (n *Decimal) IsValidDigit() bool {
	return n[0]&mask == 0x0000000100000000
}

// IsValid のほとんどのケースでは stack 領域にすら hash 計算用のメモリを確保する必要がない
// hash 計算部分を別の関数にしておくことで IsValid の時点でメモリが確保されるのを回避する
func (n *Decimal) IsValidHash() bool {
	raw := n.Bytes()
	data := raw[:34]     // 32 byte + 前後の 2 byte
	checksum := raw[34:] // 末尾の 4 byte

	hash := sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])

	return checksum[0] == hash[0] &&
		checksum[1] == hash[1] &&
		checksum[2] == hash[2] &&
		checksum[3] == hash[3]
}

// 上位の2 byte は端数なので返さない
// ループをやめて展開したら2倍くらい速くなった
func (n *Decimal) Bytes() [38]byte {
	buf := [38]byte{}

	d := n[0]
	buf[37] = byte(d)
	buf[36] = byte(d >> 8)
	buf[35] = byte(d >> 16)
	buf[34] = byte(d >> 24)
	buf[33] = byte(d >> 32)
	buf[32] = byte(d >> 40)
	buf[31] = byte(d >> 48)
	buf[30] = byte(d >> 56)

	d = n[1]
	buf[29] = byte(d)
	buf[28] = byte(d >> 8)
	buf[27] = byte(d >> 16)
	buf[26] = byte(d >> 24)
	buf[25] = byte(d >> 32)
	buf[24] = byte(d >> 40)
	buf[23] = byte(d >> 48)
	buf[22] = byte(d >> 56)

	d = n[2]
	buf[21] = byte(d)
	buf[20] = byte(d >> 8)
	buf[19] = byte(d >> 16)
	buf[18] = byte(d >> 24)
	buf[17] = byte(d >> 32)
	buf[16] = byte(d >> 40)
	buf[15] = byte(d >> 48)
	buf[14] = byte(d >> 56)

	d = n[3]
	buf[13] = byte(d)
	buf[12] = byte(d >> 8)
	buf[11] = byte(d >> 16)
	buf[10] = byte(d >> 24)
	buf[9] = byte(d >> 32)
	buf[8] = byte(d >> 40)
	buf[7] = byte(d >> 48)
	buf[6] = byte(d >> 56)

	d = n[4]
	buf[5] = byte(d)
	buf[4] = byte(d >> 8)
	buf[3] = byte(d >> 16)
	buf[2] = byte(d >> 24)
	buf[1] = byte(d >> 32)
	buf[0] = byte(d >> 40)
	// 上位の2 byte は捨てる

	return buf
}

// 用途的に不要なので桁あふれは呼び出し元に返さない
func (x *Decimal) AddMut(y *Decimal) {
	var c uint64
	x[0], c = bits.Add64(x[0], y[0], 0)
	x[1], c = bits.Add64(x[1], y[1], c)
	x[2], c = bits.Add64(x[2], y[2], c)
	x[3], c = bits.Add64(x[3], y[3], c)
	x[4], _ = bits.Add64(x[4], y[4], c)
}

// 用途的に不要なので桁あふれは呼び出し元に返さない
func (x Decimal) Mul(y uint64) Decimal {
	var z Decimal
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

// from https://github.com/anaskhan96/base58check/blob/master/base58check.go#L99
func (n *Decimal) Base58() string {
	raw := n.Bytes()

	var encoded string
	decimalData := new(big.Int)
	decimalData.SetBytes(raw[:])
	divisor, zero := big.NewInt(58), big.NewInt(0)

	for decimalData.Cmp(zero) > 0 {
		mod := new(big.Int)
		decimalData.DivMod(decimalData, divisor, mod)
		encoded = string(characters[mod.Int64()]) + encoded
	}

	return encoded
}
