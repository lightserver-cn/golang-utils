package datasource

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type RedisClient struct {
	client *redis.Client
}

type RedisOptions struct {
	Network            string        `json:"network"`
	Addr               string        `json:"addr"`
	Port               string        `json:"port"`
	Password           string        `json:"password"`
	DB                 int           `json:"db"`
	MaxRetries         int           `json:"max_retries"`
	MinRetryBackoff    time.Duration `json:"min_retry_backoff"`
	MaxRetryBackoff    time.Duration `json:"max_retry_backoff"`
	DialTimeout        time.Duration `json:"dial_timeout"`
	ReadTimeout        time.Duration `json:"read_timeout"`
	WriteTimeout       time.Duration `json:"write_timeout"`
	PoolSize           int           `json:"pool_size"`
	MinIdleConns       int           `json:"min_idle_conns"`
	MaxConnAge         time.Duration `json:"max_conn_age"`
	PoolTimeout        time.Duration `json:"pool_timeout"`
	IdleTimeout        time.Duration `json:"idle_timeout"`
	IdleCheckFrequency time.Duration `json:"idle_check_frequency"`
	ReadOnly           bool          `json:"read_only"`
}

func NewRedis(opts *RedisOptions) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Network:  opts.Network,
		Addr:     opts.Addr + ":" + opts.Port,
		Password: opts.Password,
		DB:       opts.DB,
		PoolSize: opts.PoolSize,
	})
	if _, err := client.Ping().Result(); err != nil {
		logrus.Errorf("Connected to Redis failed: %s", err.Error())
	}

	fmt.Println("Connected to Redis!")

	return client
}
