package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	database *mongo.Database
}

type MongoOptions struct {
	Addr       string         `json:"addr"`
	Port       string         `json:"port"`
	Username   string         `json:"username"`
	Password   string         `json:"password"`
	Database   string         `json:"database"`
	TimeOut    *time.Duration `json:"time_out"`
	AuthSource string         `json:"auth_source"`
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
		ConnectTimeout: opts.TimeOut,
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
