package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrdersProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_processed_total",
		Help: "Total number of processed orders",
	})

	OrdersFromCache = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_from_cache_total",
		Help: "Total number of orders served from cache",
	})

	OrdersFromDB = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_from_db_total",
		Help: "Total number of orders served from database",
	})

	CacheSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cache_size_bytes",
		Help: "Current size of cache in bytes",
	})

	CacheItems = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cache_items_count",
		Help: "Current number of items in cache",
	})

	KafkaMessagesReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "kafka_messages_received_total",
		Help: "Total number of messages received from Kafka",
	})

	DBErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_errors_total",
		Help: "Total number of database errors",
	})

	CacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "Total number of cache hits",
	})

	CacheMisses = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "Total number of cache misses",
	})
)
