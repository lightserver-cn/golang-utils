package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoOptions struct {
	Addr       string        `json:"addr"`        // 地址
	Port       string        `json:"port"`        // 端口
	Username   string        `json:"username"`    // 用户名
	Password   string        `json:"password"`    // 密码
	Database   string        `json:"database"`    // 数据库
	TimeOut    time.Duration `json:"time_out"`    // 过期时间 int64 ns * 1000 * 1000 = s
	AuthSource string        `json:"auth_source"` // 验证权限数据库
}

func NewMongo(opts *MongoOptions) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, &options.ClientOptions{
		AppName: nil,
		Auth: &options.Credential{
			AuthMechanism:           "",
			AuthMechanismProperties: nil,
			AuthSource:              opts.AuthSource,
			Username:                opts.Username,
			Password:                opts.Password,
			PasswordSet:             true,
		},
		ConnectTimeout: &opts.TimeOut,
		Hosts:          []string{opts.Addr + ":" + opts.Port},
	})

	if err != nil {
		logrus.Errorf("Connected to MongoDB failed: %s", err.Error())
	}

	if err = mongoClient.Ping(context.TODO(), nil); err != nil {
		logrus.Errorf("Ping MongoDB client failed: %s", err.Error())
	}

	fmt.Println("Connected to MongoDB!")

	return mongoClient.Database(opts.Database)
}
