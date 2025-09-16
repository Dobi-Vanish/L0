package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/logger"
	"L0/internal/repository"
	"L0/internal/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
)

func connectToDB(dsn string, maxAttempts int) (*repository.PostgresRepository, error) {
	var repo *repository.PostgresRepository
	var err error

	for i := 0; i < maxAttempts; i++ {
		repo, err = repository.New(dsn)
		if err == nil {
			return repo, nil
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxAttempts, err)
		time.Sleep(5 * time.Second)
	}

	return nil, err
}

func main() {
	cfg := config.Load()

	mongoLogger, err := logger.New(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		mongoLogger.Log("ERROR", "consumer", "Failed to connect to database: "+err.Error())
		log.Fatal("Failed to connect to database:", err)
	}
	defer mongoLogger.Close()

	repo, err := connectToDB(cfg.PostgresDSN, 5)
	if err != nil {
		mongoLogger.Log("ERROR", "consumer", "Failed to connect to database: "+err.Error())
		log.Fatal("Failed to connect to database:", err)
	}
	defer repo.Close()

	cache := cache.New(cfg.CacheMaxSize, time.Duration(cfg.CacheTTLMinutes)*time.Minute)
	defer cache.Stop()

	orderService := service.NewOrderService(repo, cache)

	if err := orderService.RestoreCacheFromDB(); err != nil {
		mongoLogger.Log("ERROR", "consumer", "Failed to restore cache: "+err.Error())
		log.Fatal("Failed to restore cache:", err)
	}

	go func() {
		healthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := repo.Ping(); err != nil {
				http.Error(w, "Database not connected", http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		http.Handle("/health", healthHandler)
		http.Handle("/metrics", promhttp.Handler())

		log.Println("Metrics server started on :8082")
		http.ListenAndServe(":8082", nil)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var reader *kafka.Reader
	for i := 0; i < 5; i++ {
		reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  cfg.KafkaBrokers,
			Topic:    cfg.KafkaTopic,
			GroupID:  cfg.KafkaConsumerGroup,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		})

		conn, err := kafka.Dial("tcp", cfg.KafkaBrokers[0])
		if err == nil {
			conn.Close()
			break
		}

		log.Printf("Failed to connect to Kafka (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)

		if i == 4 {
			mongoLogger.Log("ERROR", "consumer", "Failed to connect to Kafka: "+err.Error())
			log.Fatal("Failed to connect to Kafka:", err)
		}
	}
	defer reader.Close()

	log.Println("Consumer started")
	mongoLogger.Log("INFO", "consumer", "Consumer started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down consumer...")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				mongoLogger.Log("ERROR", "consumer", "Error reading message: "+err.Error())
				log.Println("Error reading message:", err)
				continue
			}

			if err := orderService.ProcessOrderMessage(msg.Value); err != nil {
				mongoLogger.Log("ERROR", "consumer", "Error processing message: "+err.Error())
				log.Println("Error processing message:", err)
				continue
			}

			log.Printf("Order %s processed successfully", string(msg.Key))
		}
	}
}
