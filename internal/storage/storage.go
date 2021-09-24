package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mohammadne/mosaicman/internal/models"
	"github.com/mohammadne/mosaicman/pkg/logger"
)

type Storage interface {
	Persist(context.Context, io.Reader, *models.Metadata) error
	Retrieve(context.Context, *models.Metadata) (*os.File, error)
}

type storage struct {
	config *Config
	logger logger.Logger
	pool   *redis.Pool
}

func New(cfg *Config, lg logger.Logger) (Storage, error) {
	s := &storage{config: cfg, logger: lg}

	s.pool = s.newPool()
	cleanUp(s.pool)
	go s.subscribe()

	return s, nil
}

func (s *storage) newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: s.config.IdleTimeout,
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			connection, err := redis.DialContext(ctx, "tcp", s.config.URL)
			if err != nil {
				return nil, err
			}
			return connection, err
		},
		TestOnBorrow: func(connection redis.Conn, t time.Time) error {
			_, err := connection.Do("PING")
			return err
		},
	}
}

func (s *storage) newConnection(ctx context.Context) (redis.Conn, error) {
	connection, err := s.pool.GetContext(ctx)
	if err != nil {
		s.logger.Error("error getting connection", logger.Error(err))
		return nil, err
	}
	return connection, nil
}

// TODO: REMOVE FILE
func (s *storage) subscribe() error {
	connection, err := s.newConnection(context.TODO())
	if err != nil {
		return err
	}
	defer connection.Close()

	psc := redis.PubSubConn{Conn: connection}
	psc.PSubscribe("__key*__:*")
	for {
		switch msg := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("Message: %s %s\n", msg.Channel, msg.Data)
		// case redis.PMessage:
		// 	fmt.Printf("PMessage: %s %s %s\n", msg.Pattern, msg.Channel, msg.Data)
		case redis.Subscription:
			fmt.Printf("Subscription: %s %s %d\n", msg.Kind, msg.Channel, msg.Count)
		case error:
			return fmt.Errorf("error: %v", msg)
		default:
			fmt.Println("DEFAULT")
		}
	}

}

func cleanUp(pool *redis.Pool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	// signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		pool.Close()
		os.Exit(0)
	}()
}
