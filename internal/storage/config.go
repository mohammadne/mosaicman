package storage

import (
	"time"
)

type Mode uint8

const (
	Single  Mode = 0
	Cluster Mode = 1
)

// Config is the config of the cache client.
type Config struct {
	Mode               int           `split_words:"true" default:"0"`
	URL                string        `split_words:"true"`
	MasterURL          string        `split_words:"true"`
	SlaveURL           string        `split_words:"true"`
	Password           string        `default:"" split_words:"true"`
	Expiration         int           `default:"10"`
	PoolSize           int           `split_words:"true" default:"10"`
	MaxRetries         int           `split_words:"true" default:"0"`
	ReadTimeout        time.Duration `split_words:"true" default:"3s"`
	PoolTimeout        time.Duration `split_words:"true" default:"4s"`
	MinRetryBackoff    time.Duration `split_words:"true" default:"8ms"`
	MaxRetryBackoff    time.Duration `split_words:"true" default:"512ms"`
	IdleTimeout        time.Duration `split_words:"true" default:"300s"`
	IdleCheckFrequency time.Duration `split_words:"true" default:"60s"`
	SetMemberExpTime   time.Duration `split_words:"true" default:"300s"`
}

// newSingleRedis returns a new `RedisHandler` with a single Redis client.
// func newSingleRedis(cfg *Config) *redis.Client {
// 	return redis.NewClient(&redis.Options{
// 		Addr:               cfg.URL,
// 		Password:           cfg.Password,
// 		MaxRetries:         cfg.MaxRetries,
// 		MinRetryBackoff:    cfg.MinRetryBackoff,
// 		MaxRetryBackoff:    cfg.MaxRetryBackoff,
// 		ReadTimeout:        cfg.ReadTimeout,
// 		PoolSize:           cfg.PoolSize,
// 		PoolTimeout:        cfg.PoolTimeout,
// 		IdleTimeout:        cfg.IdleTimeout,
// 		IdleCheckFrequency: cfg.IdleCheckFrequency,
// 	})
// }

// newClusterRedis returns a new `RedisHandler` with a clustered Redis client.
// func newClusterRedis(cfg *Config) *redis.ClusterClient {
// 	return redis.NewClusterClient(&redis.ClusterOptions{
// 		Addrs:              []string{cfg.MasterURL, cfg.SlaveURL},
// 		Password:           cfg.Password,
// 		MaxRetries:         cfg.MaxRetries,
// 		MinRetryBackoff:    cfg.MinRetryBackoff,
// 		MaxRetryBackoff:    cfg.MaxRetryBackoff,
// 		ReadTimeout:        cfg.ReadTimeout,
// 		PoolSize:           cfg.PoolSize,
// 		PoolTimeout:        cfg.PoolTimeout,
// 		IdleTimeout:        cfg.IdleTimeout,
// 		IdleCheckFrequency: cfg.IdleCheckFrequency,
// 	})
// }
