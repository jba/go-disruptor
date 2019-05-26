package implementations

import (
	"fmt"
	"testing"
)

var caps = []int{ /*1, 8, */ 64, 1024 /*, 1024 * 16*/}

func BenchmarkChannel(b *testing.B) {
	const n = 1e6
	for _, cap := range caps {
		b.Run(fmt.Sprint(cap), func(b *testing.B) {
			c := make(chan int64, cap)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				go func() {
					for j := 0; j < n; j++ {
						c <- int64(j)
					}
				}()
				for j := 0; j < n; j++ {
					_ = <-c
				}
			}
		})
	}
}

func BenchmarkLocked(b *testing.B) {
	const n = 1e6
	for _, cap := range caps {
		b.Run(fmt.Sprint(cap), func(b *testing.B) {
			q := NewLockedQueue(cap)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				go func() {
					for j := 0; j < n; j++ {
						q.Put(int64(j))
					}
				}()
				for j := 0; j < n; j++ {
					_ = q.Get()
				}
			}
		})
	}
}

func BenchmarkLockFree(b *testing.B) {
	const n = 1e6
	for _, cap := range caps {
		b.Run(fmt.Sprint(cap), func(b *testing.B) {
			q := NewLockFreeQueue(cap)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				go func() {
					for j := 0; j < n; j++ {
						q.Put(int64(j))
					}
				}()
				for j := 0; j < n; j++ {
					_ = q.Get()
				}
			}
		})
	}
}
