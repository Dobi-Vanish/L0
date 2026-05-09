package cache

import (
	models "L0/internal/model"
	"testing"
	"time"
)

func BenchmarkCacheSet(b *testing.B) {
	c := New(1e9, 10*time.Minute)
	order := &models.Order{OrderUID: "uid", TrackNumber: "T"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(order)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	c := New(1e9, 10*time.Minute)
	order := &models.Order{OrderUID: "uid", TrackNumber: "T"}
	c.Set(order)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get("uid")
	}
}
