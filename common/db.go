package common

import (
	"SuperCrawler/conf"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"os"
)


func GetMgoCli() *mongo.Client{
	uri := fmt.Sprintf("mongodb://%s:%s", conf.Cfg.MongoDB.Host, conf.Cfg.MongoDB.Port)
	var clientOptions *options.ClientOptions
	if os.Getenv("DEBUG") == "1" {
		clientOptions = options.Client().ApplyURI(uri).SetWriteConcern(writeconcern.New(writeconcern.W(0)))
	}else{
		credential := options.Credential{
			AuthMechanism: conf.Cfg.MongoDB.AuthMechanism,
			AuthSource: conf.Cfg.MongoDB.AuthSource,
			Username: conf.Cfg.MongoDB.Username,
			Password: conf.Cfg.MongoDB.Password,
		}
		clientOptions = options.Client().ApplyURI(uri).SetAuth(credential).SetWriteConcern(writeconcern.New(writeconcern.W(0)))
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Error(fmt.Sprintf("[GetMongoCli] %s", err))
	}

	return client
}
