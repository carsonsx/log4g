package log4g

import "testing"

func BenchmarkLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Debug(i)
	}
}



