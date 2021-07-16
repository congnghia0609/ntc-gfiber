/**
 *
 * @author nghiatc
 * @since Jan 4, 2021
 */

package mdb

import (
	"context"
	"fmt"
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
	ID  string `bson:"_id" json:"id"`
	Seq int64  `bson:"seq" json:"seq"`
}

func (id NId) String() string {
	return fmt.Sprintf("NId{Id: %v, Seq: %v}", id.ID, id.Seq)
}

// Next generate a auto increment version ID for the given key
func Next(id string) (int64, error) {
	var nid NId

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$inc", bson.D{{"seq", int64(1)}}}}
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{ReturnDocument: &after, Upsert: &upsert}

	collection := FiberDB.Collection(TableNId)
	err := collection.FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(&nid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil
		} else {
			log.Println(err)
			return 0, err
		}
	}
	// log.Println("nid:", nid)

	return nid.Seq, nil
}

// ResetID reset id gen to value
func ResetID(id string, value int64) (int64, error) {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"seq", value}}}}

	collection := FiberDB.Collection(TableNId)
	result, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return 0, err
	}
	// log.Println("nid:", nid)

	return result.MatchedCount, nil
}
