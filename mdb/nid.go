/**
 *
 * @author nghiatc
 * @since Jan 4, 2021
 */

package mdb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Const
const (
	TableNId = "nlidgen"
)

// NId is long id gen
type NId struct {
	ID  string `bson:"_id" json:"_id"`
	Seq int64  `bson:"seq" json:"seq"`
}

// Next generate a auto increment version ID for the given key
func Next(id string) (int64, error) {
	var nid NId
	client := GetClient()
	defer Close(client)
	collection := client.Database(DbName).Collection(TableNId)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$inc", bson.D{{"seq", 1}}}}

	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&nid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil
		} else {
			log.Println(err)
			return 0, err
		}
	}
	// log.Println("nid:", nid)

	return nid.Seq + 1, nil
}

// ResetID reset id gen to value
func ResetID(id string, value int64) (int64, error) {
	client := GetClient()
	defer Close(client)
	collection := client.Database(DbName).Collection(TableNId)

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"seq", value}}}}

	result, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return 0, err
	}
	// log.Println("nid:", nid)

	return result.MatchedCount, nil
}
