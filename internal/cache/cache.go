package cache

import (
	"L0/internal/metrics"
	"L0/internal/model"
	"sync"
	"time"
)

type CacheItem struct {
	Order     *models.Order
	ExpiresAt time.Time
	Size      int64
}

type Cache struct {
	sync.RWMutex
	items       map[string]*CacheItem
	maxSize     int64
	currentSize int64
	ttl         time.Duration
	stopChan    chan bool
	hits        int64
	misses      int64
}

func New(maxSize int64, ttl time.Duration) *Cache {
	cache := &Cache{
		items:    make(map[string]*CacheItem),
		maxSize:  maxSize,
		ttl:      ttl,
		stopChan: make(chan bool),
	}

	go cache.cleanup()

	return cache
}

func (c *Cache) Set(order *models.Order) {
	c.Lock()
	defer c.Unlock()

	if existing, exists := c.items[order.OrderUID]; exists {
		c.currentSize -= existing.Size
		delete(c.items, order.OrderUID)
	}

	size := int64(len(order.OrderUID)) + int64(len(order.TrackNumber)) + 1024

	if c.currentSize+size > c.maxSize {
		c.evictOldest()
	}

	c.items[order.OrderUID] = &CacheItem{
		Order:     order,
		ExpiresAt: time.Now().Add(c.ttl),
		Size:      size,
	}
	c.currentSize += size
	metrics.CacheSize.Set(float64(c.currentSize))
	metrics.CacheItems.Set(float64(len(c.items)))
}

func (c *Cache) Get(uid string) (*models.Order, bool) {
	c.RLock()
	defer c.RUnlock()

	item, exists := c.items[uid]
	if !exists || time.Now().After(item.ExpiresAt) {
		c.misses++
		metrics.CacheMisses.Inc()
		return nil, false
	}

	c.hits++
	metrics.CacheHits.Inc()
	return item.Order, true
}

func (c *Cache) Restore(orders []*models.Order) {
	c.Lock()
	defer c.Unlock()

	for _, order := range orders {
		size := int64(len(order.OrderUID)) + int64(len(order.TrackNumber)) + 1024
		c.items[order.OrderUID] = &CacheItem{
			Order:     order,
			ExpiresAt: time.Now().Add(c.ttl),
			Size:      size,
		}
		c.currentSize += size
	}
	metrics.CacheSize.Set(float64(c.currentSize))
	metrics.CacheItems.Set(float64(len(c.items)))
}

func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestTime.IsZero() || item.ExpiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.ExpiresAt
		}
	}

	if oldestKey != "" {
		c.currentSize -= c.items[oldestKey].Size
		delete(c.items, oldestKey)
	}
	metrics.CacheSize.Set(float64(c.currentSize))
	metrics.CacheItems.Set(float64(len(c.items)))
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Lock()
			now := time.Now()
			for key, item := range c.items {
				if now.After(item.ExpiresAt) {
					c.currentSize -= item.Size
					delete(c.items, key)
				}
			}
			c.Unlock()
		case <-c.stopChan:
			return
		}
	}

	metrics.CacheItems.Set(float64(len(c.items)))
}

func (c *Cache) Stop() {
	close(c.stopChan)
}

func (c *Cache) CurrentSize() int64 {
	c.RLock()
	defer c.RUnlock()
	return c.currentSize
}

func (c *Cache) ItemsCount() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.items)
}

func (c *Cache) MaxSize() int64 {
	return c.maxSize
}

func (c *Cache) Hits() int64 {
	c.RLock()
	defer c.RUnlock()
	return c.hits
}

func (c *Cache) Misses() int64 {
	c.RLock()
	defer c.RUnlock()
	return c.misses
}

func (c *Cache) HitRatio() float64 {
	c.RLock()
	defer c.RUnlock()
	total := c.hits + c.misses
	if total == 0 {
		return 0
	}
	return float64(c.hits) / float64(total)
}

func (c *Cache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.items = make(map[string]*CacheItem)
	c.currentSize = 0
	c.hits = 0
	c.misses = 0
}
