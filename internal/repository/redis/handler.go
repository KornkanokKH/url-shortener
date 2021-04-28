package redis

//go:generate mockgen -source=./handler.go -destination=./mocks/handler.go

import (
	"github.com/go-redis/redis"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"time"
)

var (
	consoleLog zerolog.Logger
)

type HandlerInterface interface {
	Connect(config Config) error
	Disconnect()
	Set(key string, value interface{}, exp time.Duration, txn *newrelic.Transaction) error
	Get(key string, txn *newrelic.Transaction) (string, error)
	Del(key string, txn *newrelic.Transaction) error
}
type Handler struct {
	client *redis.Client
}

func (handler *Handler) Connect(config Config) error {
	address := config.RedisServer.Address + ":" + config.RedisServer.Port

	client := redis.NewClient(&redis.Options{
		Addr:       address,
		DB:         config.RedisServer.Db,
		MaxRetries: 3,
	})

	_, err := client.Ping().Result()
	if err != nil {
		consoleLog.Error().Msgf("Unexpected error to connect Redis address: %v, err: %v", address, err)
		return err
	}
	handler.client = client
	consoleLog.Debug().Msg("Successfully connect Redis server")
	return nil
}

func (handler *Handler) Disconnect() {
	if err := handler.client.Close(); err != nil {
		consoleLog.Warn().Msgf("Unexpected error to close Redis client connection, error: %v", err)
	}
}

func (handler *Handler) Set(key string, value interface{}, exp time.Duration, txn *newrelic.Transaction) error {
	segment := newrelic.DatastoreSegment{
		StartTime:          txn.StartSegmentNow(),
		Product:            newrelic.DatastoreRedis,
		Operation:          "SET",
		ParameterizedQuery: key,
	}
	defer segment.End()

	return handler.client.Set(key, value, exp).Err()
}

func (handler *Handler) Get(key string, txn *newrelic.Transaction) (string, error) {

	segment := newrelic.DatastoreSegment{
		StartTime:          txn.StartSegmentNow(),
		Product:            newrelic.DatastoreRedis,
		Operation:          "GET",
		ParameterizedQuery: key,
	}

	defer segment.End()

	result, err := handler.client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return result, err
}

func (handler *Handler) Del(key string, txn *newrelic.Transaction) error {

	segment := newrelic.DatastoreSegment{
		StartTime:          txn.StartSegmentNow(),
		Product:            newrelic.DatastoreRedis,
		Operation:          "DELETE",
		ParameterizedQuery: key,
	}

	defer segment.End()

	err := handler.client.Del(key)
	if err != nil {
		return err.Err()
	}

	return nil
}
