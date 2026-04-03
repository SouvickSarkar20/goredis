package store

import (
	"fmt"
	"testing"
)

func BenchmarkStoreSet(b *testing.B) {
	s := NewStore()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(fmt.Sprintf("key-%d", i), "value", 0)
	}
}

func BenchmarkStoreGet(b *testing.B) {
	s := NewStore()
	for i := 0; i < 1000; i++ {
		s.Set(fmt.Sprintf("key-%d", i), "value", 0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get(fmt.Sprintf("key-%d", i%1000))
	}
}

func BenchmarkStoreDelete(b *testing.B) {
	s := NewStore()
	// Pre-populate keys
	for i := 0; i < b.N; i++ {
		s.Set(fmt.Sprintf("key-%d", i), "value", 0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Delete(fmt.Sprintf("key-%d", i))
	}
}
