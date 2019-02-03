package main

import (
	"testing"
)

var result uint32
var resultA uint16

func BenchmarkGetErrorCode(b *testing.B) {
	var r uint32
	for n := 0; n < b.N; n++ {
		r = getErrorCode(1234, 189, 4321)
	}
	result = r
}

func BenchmarkDecodeErrorCode(b *testing.B) {
	var ra uint16
	for n := 0; n < b.N; n++ {
		ra, _, _ = decodeErrorCode(1331331)
	}
	resultA = ra
}
