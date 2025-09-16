package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBrokers       []string
	KafkaTopic         string
	KafkaPartitions    int
	KafkaReplicas      int
	PostgresDSN        string
	MongoURI           string
	MongoDatabase      string
	HTTPServerPort     string
	KafkaConsumerGroup string
	LogLevel           string
	CacheMaxSize       int64
	CacheTTLMinutes    int
}

func Load() *Config {
	err := godotenv.Load("configs/config.env")
	if err != nil {
		panic("Failed to load config file: " + err.Error())
	}

	return &Config{
		KafkaBrokers:       strings.Split(getEnvRequired("KAFKA_BROKERS"), ","),
		KafkaTopic:         getEnvRequired("KAFKA_TOPIC"),
		KafkaPartitions:    getEnvAsIntRequired("KAFKA_PARTITIONS"),
		KafkaReplicas:      getEnvAsIntRequired("KAFKA_REPLICAS"),
		PostgresDSN:        getEnvRequired("POSTGRES_DSN"),
		MongoURI:           getEnvRequired("MONGODB_URI"),
		MongoDatabase:      getEnvRequired("MONGODB_DATABASE"),
		HTTPServerPort:     getEnvRequired("HTTP_SERVER_PORT"),
		KafkaConsumerGroup: getEnvRequired("KAFKA_CONSUMER_GROUP"),
		LogLevel:           getEnvRequired("LOG_LEVEL"),
		CacheMaxSize:       getEnvAsInt64Required("CACHE_MAX_SIZE"),
		CacheTTLMinutes:    getEnvAsIntRequired("CACHE_TTL_MINUTES"),
	}
}

func getEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic("Environment variable " + key + " is required")
	}
	return value
}

func getEnvAsIntRequired(key string) int {
	valueStr := getEnvRequired(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic("Environment variable " + key + " must be an integer: " + err.Error())
	}
	return value
}

func getEnvAsInt64Required(key string) int64 {
	valueStr := getEnvRequired(key)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		panic("Environment variable " + key + " must be an integer: " + err.Error())
	}
	return value
}
