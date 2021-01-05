/**
 *
 * @author nghiatc
 * @since Jan 3, 2021
 */

package mdb

import (
	"context"
	"log"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Const
const (
	DbName = "fiberdb"
)

// GetClient return mongo client
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Connect
// mongodb://localhost:27017
// mongodb://localhost:27017,localhost:27018/?replicaSet=replset
// mongodb://user:password@localhost:27017
func GetClient() *mongo.Client {
	c := nconf.GetConfig()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(c.GetString("mongodb.uri")))
	if err != nil {
		log.Println(err)
		return nil
	}
	return client
}

// Close disconnect mongo client
func Close(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Println(err)
	}
}
