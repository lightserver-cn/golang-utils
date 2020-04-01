package datasource

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type RedisOptions struct {
	Network            string        `json:"network"`              // 网络类型 tcp 或者是 unix
	Addr               string        `json:"addr"`                 // ip
	Port               string        `json:"port"`                 // 端口
	Password           string        `json:"password"`             // 密码
	DB                 int           `json:"db"`                   // 数据库
	MaxRetries         int           `json:"max_retries"`          // 放弃连接前的最大重试次数
	MinRetryBackoff    time.Duration `json:"min_retry_backoff"`    // 每次重试直接的最小回退  int64 ns * 1000 * 1000 = s
	MaxRetryBackoff    time.Duration `json:"max_retry_backoff"`    // 每次重试直接的最大回退
	DialTimeout        time.Duration `json:"dial_timeout"`         // 建立新连接的拨号超时
	ReadTimeout        time.Duration `json:"read_timeout"`         // 套接字读取超时。如果达到，命令将失败，超时而不是阻塞
	WriteTimeout       time.Duration `json:"write_timeout"`        // 套接字写入超时。如果达到，命令将失败，超时而不是阻塞
	PoolSize           int           `json:"pool_size"`            // 套接字连接的最大数目
	MinIdleConns       int           `json:"min_idle_conns"`       // 最小空闲连接数，新连接是慢的
	MaxConnAge         time.Duration `json:"max_conn_age"`         // 客户端退出（关闭）连接的连接期限
	PoolTimeout        time.Duration `json:"pool_timeout"`         // 在返回错误之前，所有连接正忙，等待的时间
	IdleTimeout        time.Duration `json:"idle_timeout"`         // 客户端关闭空闲连接的时间。应该小雨服务器的超时
	IdleCheckFrequency time.Duration `json:"idle_check_frequency"` // 空闲连接收割器进行空闲检查的频率
}

func NewRedis(opts *RedisOptions) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Network:            opts.Network,
		Addr:               opts.Addr + ":" + opts.Port,
		Password:           opts.Password,
		DB:                 opts.DB,
		MaxRetries:         opts.MaxRetries,
		MinRetryBackoff:    opts.MinRetryBackoff,
		MaxRetryBackoff:    opts.MaxRetryBackoff,
		DialTimeout:        opts.DialTimeout,
		ReadTimeout:        opts.ReadTimeout,
		WriteTimeout:       opts.WriteTimeout,
		PoolSize:           opts.PoolSize,
		MinIdleConns:       opts.PoolSize,
		MaxConnAge:         opts.MaxConnAge,
		PoolTimeout:        opts.PoolTimeout,
		IdleTimeout:        opts.IdleTimeout,
		IdleCheckFrequency: opts.IdleCheckFrequency,
	})
	if _, err := client.Ping().Result(); err != nil {
		logrus.Errorf("Connected to Redis failed: %s", err.Error())
	}

	fmt.Println("Connected to Redis!")

	return client
}
