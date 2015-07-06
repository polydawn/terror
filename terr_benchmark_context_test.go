package terror

import (
	"testing"
)

//
// benchmarks for CONTEXT
//
// these attempt to give some semblance of meaning to numbers by making them
// comparable to other things you do in a real program.  after all, as
// a pragmatic and realistic soul, you only care about the performance of your
// error handling if it rises to more than 0.0001% of your runtime, right?
//

func Benchmark_Context_Copy1kb(b *testing.B) {
	src := make([]byte, 1024)
	dst := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		copy(dst, src)
	}
}

func Benchmark_Context_Copy1kbPiecemeal(b *testing.B) {
	src := make([]byte, 1024)
	dst := make([]byte, 1024)
	fn := func(src, dst []byte) {
		copy(dst, src)
	}
	for i := 0; i < b.N; i++ {
		fn(dst[0:256], src[0:256])
		fn(dst[256:512], src[256:512])
		fn(dst[512:768], src[512:768])
		fn(dst[768:1024], src[768:1024])
	}
}
