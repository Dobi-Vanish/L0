package main

import (
	models "L0/internal/model"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/logger"
	"L0/internal/repository"
	"L0/internal/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Load()

	var mongoLogger *logger.MongoLogger
	var err error
	maxRetries := 5
	initialDelay := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		mongoLogger, err = logger.New(cfg.MongoURI, cfg.MongoDatabase)
		if err == nil {
			log.Printf("Successfully connected to MongoDB on attempt %d", attempt)
			break
		}

		if attempt == maxRetries {
			log.Printf("Failed to connect to MongoDB after %d attempts: %v", maxRetries, err)
			mongoLogger = &logger.MongoLogger{}
			break
		}

		delay := initialDelay * time.Duration(1<<uint(attempt-1))
		log.Printf("MongoDB connection attempt %d failed: %v. Retrying in %v", attempt, err, delay)
		time.Sleep(delay)
	}

	defer func() {
		if mongoLogger != nil {
			mongoLogger.Close()
		}
	}()

	repo, err := repository.New(cfg.PostgresDSN)
	if err != nil {
		mongoLogger.Log("ERROR", "server", "Failed to connect to database: "+err.Error())
		log.Fatal("Failed to connect to database:", err)
	}
	defer repo.Close()

	cache := cache.New(cfg.CacheMaxSize, time.Duration(cfg.CacheTTLMinutes)*time.Minute)
	defer cache.Stop()

	orderService := service.NewOrderService(repo, cache)

	if err := orderService.RestoreCacheFromDB(); err != nil {
		mongoLogger.Log("ERROR", "server", "Failed to restore cache: "+err.Error())
		log.Fatal("Failed to restore cache:", err)
	}

	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		uid := r.URL.Path[len("/order/"):]
		if uid == "" {
			http.Error(w, "Order ID required", http.StatusBadRequest)
			return
		}

		order, err := orderService.GetOrderByID(uid)
		if err != nil {
			mongoLogger.Log("ERROR", "server", "Error getting order: "+err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if order == nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	})

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := repo.Ping(); err != nil {
			http.Error(w, "Database not connected", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/add_order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			mongoLogger.Log("WARN", "server", "Invalid method for /add_order: "+r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			mongoLogger.Log("WARN", "server", "Invalid Content-Type for /add_order: "+r.Header.Get("Content-Type"))
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		var order models.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			mongoLogger.Log("ERROR", "server", "Error decoding JSON: "+err.Error())
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := orderService.CreateOrder(&order); err != nil {
			mongoLogger.Log("ERROR", "server", "Error creating order: "+err.Error())

			if strings.Contains(err.Error(), "validation error") {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else if strings.Contains(err.Error(), "already exists") {
				http.Error(w, err.Error(), http.StatusConflict)
			} else {
				http.Error(w, "Error creating order: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		mongoLogger.Log("INFO", "server", "Order created successfully: "+order.OrderUID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "success",
			"message":   "Order created successfully",
			"order_uid": order.OrderUID,
		})
	})

	http.HandleFunc("/cache/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		stats := map[string]interface{}{
			"size":        cache.CurrentSize(),
			"items_count": cache.ItemsCount(),
			"max_size":    cache.MaxSize(),
			"hits":        cache.Hits(),
			"misses":      cache.Misses(),
			"hit_ratio":   cache.HitRatio(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	http.HandleFunc("/cache/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		orderUID := r.URL.Query().Get("order_uid")
		if orderUID == "" {
			http.Error(w, "order_uid parameter is required", http.StatusBadRequest)
			return
		}

		order, exists := cache.Get(orderUID)
		if !exists {
			http.Error(w, "Order not found in cache", http.StatusNotFound)
			return
		}

		dbOrder, err := repo.GetOrderByID(orderUID)
		if err != nil {
			mongoLogger.Log("ERROR", "server", "Error getting order from DB: "+err.Error())
		}

		response := map[string]interface{}{
			"in_cache": true,
			"in_db":    dbOrder != nil,
			"order":    order,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/cache/clear", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cache.Clear()
		mongoLogger.Log("INFO", "server", "Cache cleared manually")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cache cleared successfully"))
	})

	log.Printf("HTTP server started on :%s", cfg.HTTPServerPort)
	mongoLogger.Log("INFO", "server", "HTTP server started on :"+cfg.HTTPServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPServerPort, nil))
}
