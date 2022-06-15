package safemap

import "testing"

//BenchmarkSafeMap_Set
func BenchmarkSafeMap_Set(b *testing.B) {
	m := New[int, int]()
	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(1024)
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			m.Set(i, i)
		}
	})
}

//BenchmarkSafeMap_Get
func BenchmarkSafeMap_Get(b *testing.B) {
	m := New[int, int]()
	for i := 0; i < 1024; i++ {
		m.Set(i, i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(1024)
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			m.Get(i)
		}
	})
}
