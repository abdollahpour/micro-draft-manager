package util

import (
	"context"
	"log"
	"net/url"

	"github.com/abdollahpour/almaniha-draft/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectMongoDB(conf config.Configuration) *mongo.Database {
	//context.WithTimeout(context.Background(), 30*time.Second)
	//_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	var err error
	if client == nil {
		clientOptions := options.Client().ApplyURI(conf.MongoUri)
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	u, err := url.Parse(conf.MongoUri)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(u.Path[1:])
}

func DisconnectMongoDB() {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
