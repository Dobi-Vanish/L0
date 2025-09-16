package cache

import (
	"L0/internal/model"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := New(1024*1024, time.Minute)

	order := &models.Order{
		OrderUID:    "test-order",
		TrackNumber: "test-track",
	}

	// Test setting and getting an order
	cache.Set(order)
	retrieved, exists := cache.Get("test-order")
	assert.True(t, exists)
	assert.Equal(t, order.OrderUID, retrieved.OrderUID)
}

func TestCache_Expiration(t *testing.T) {
	cache := New(1024*1024, time.Millisecond*10)

	order := &models.Order{
		OrderUID:    "test-order",
		TrackNumber: "test-track",
	}

	// Test that item expires after TTL
	cache.Set(order)
	time.Sleep(time.Millisecond * 20)
	_, exists := cache.Get("test-order")
	assert.False(t, exists)
}

func TestCache_Eviction(t *testing.T) {
	// Create a small cache that can hold about 2 items
	// Each item is approximately: 1 (OrderUID) + 4 (TrackNumber) + 1024 = 1029 bytes
	// So for 2 items we need at least 2058 bytes
	cache := New(2058, time.Minute)

	// Add 3 items - the oldest should be evicted
	for i := 0; i < 3; i++ {
		order := &models.Order{
			OrderUID:    strconv.Itoa(i), // "0", "1", "2"
			TrackNumber: "test",
		}
		cache.Set(order)
	}

	// The oldest item ("0") should have been evicted
	_, exists := cache.Get("0")
	assert.False(t, exists)

	// The newer items should still be there
	_, exists = cache.Get("1")
	assert.True(t, exists)

	_, exists = cache.Get("2")
	assert.True(t, exists)

	// Check that cache size doesn't exceed maximum
	assert.True(t, cache.CurrentSize() <= cache.MaxSize())
}

func TestCache_Restore(t *testing.T) {
	cache := New(1024*1024, time.Minute)

	orders := []*models.Order{
		{
			OrderUID:    "order-1",
			TrackNumber: "track-1",
		},
		{
			OrderUID:    "order-2",
			TrackNumber: "track-2",
		},
	}

	// Test restoring multiple orders
	cache.Restore(orders)

	// Check that both orders are in cache
	order1, exists := cache.Get("order-1")
	assert.True(t, exists)
	assert.Equal(t, "order-1", order1.OrderUID)

	order2, exists := cache.Get("order-2")
	assert.True(t, exists)
	assert.Equal(t, "order-2", order2.OrderUID)

	// Check items count
	assert.Equal(t, 2, cache.ItemsCount())
}

func TestCache_Stats(t *testing.T) {
	cache := New(1024*1024, time.Minute)

	// Initial stats should be zero
	assert.Equal(t, int64(0), cache.Hits())
	assert.Equal(t, int64(0), cache.Misses())
	assert.Equal(t, 0.0, cache.HitRatio())

	// Test miss
	_, exists := cache.Get("nonexistent")
	assert.False(t, exists)
	assert.Equal(t, int64(1), cache.Misses())
	assert.Equal(t, 0.0, cache.HitRatio())

	// Test hit
	order := &models.Order{
		OrderUID:    "test-order",
		TrackNumber: "test-track",
	}
	cache.Set(order)

	_, exists = cache.Get("test-order")
	assert.True(t, exists)
	assert.Equal(t, int64(1), cache.Hits())
	assert.Equal(t, int64(1), cache.Misses())
	assert.Equal(t, 0.5, cache.HitRatio())
}

func TestCache_Clear(t *testing.T) {
	cache := New(1024*1024, time.Minute)

	order := &models.Order{
		OrderUID:    "test-order",
		TrackNumber: "test-track",
	}

	cache.Set(order)
	assert.Equal(t, 1, cache.ItemsCount())

	// Test clearing cache
	cache.Clear()
	assert.Equal(t, 0, cache.ItemsCount())
	assert.Equal(t, int64(0), cache.CurrentSize())

	// Should not find the order after clear
	_, exists := cache.Get("test-order")
	assert.False(t, exists)
}
