package utils

import (
	"fmt"
	"testing"
)

func BenchmarkRandSeq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandSeq(32)
	}
}

func BenchmarkRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Rand(32)
	}
}

// 随机构建长度为20 的字符串
func ExampleRandSeq() {
	str := RandSeq(20)

	fmt.Println(str)
	fmt.Println(len(str))

	// Output:
	// Hujx08qkpUfTQe3jUn5S
	// 20
}

// 随机构造长度为20的 byte 序列
func ExampleRand() {
	bytes := Rand(20)

	fmt.Println(bytes)
	fmt.Println(len(bytes))

	// Output:
	// [113 160 217 150 46 16 189 211 0 166 226 74 30 141 16 190 242 153 212 246]
	// 20
}