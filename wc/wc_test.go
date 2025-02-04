package main

import (
	"io"
	"os"
	"runtime"
	"testing"
)

func BenchmarkWC(b *testing.B) {
	b.ReportAllocs()

	args := wcArgs{}
	file, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		result := wc(args, data)
		runtime.KeepAlive(result)
	}
}
