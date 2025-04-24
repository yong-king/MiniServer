package base62

import (
	"math"
	"strings"
)



var (
	base62String string
)

func MustInit(bs string){
	if len(bs) == 0{
		panic("need base string")
	}
	base62String = bs
}

// Int2String 十进制转62进制字符串
func Int2String(seq uint64) string{
	if seq < 61 {
		return string(base62String[seq])
	}
	bl := []byte{}
	for seq > 0{
		mod := seq % 62
		div := seq / 62
		bl = append(bl, base62String[mod])
		seq = div
	}
	return string(reverse(bl))

}

// String2Int 62进制字符串转十进制
func String2Int(s string) (seq uint64){
	bl := []byte(s)
	bl = reverse(bl)
	for i, b := range bl{
		base := math.Pow(62, float64(i))
		seq += uint64(base) * uint64(strings.Index(base62String, string(b)))
	}
	return seq
}

func reverse(s []byte) []byte{
	for i, j := 0, len(s) - 1; i < len(s) / 2; i, j = i+1, j-1{
		s[i], s[j] = s[j], s[i]
	}
	return s
}