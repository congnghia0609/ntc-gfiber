/**
 *
 * @author nghiatc
 * @since Jan 5, 2021
 */

package tag

import (
	"context"
	"fmt"
	"github.com/congnghia0609/ntc-gfiber/mdb"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Const
const (
	TableTag = "tag"
)

// Tag struct
type Tag struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"-"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// InsertTag Insert Tag
func InsertTag(tag Tag) error {
	collection := mdb.FiberDB.Collection(TableTag)
	insertResult, err := collection.InsertOne(context.Background(), tag)
	if err != nil {
		fmt.Println("Error Inserted tag with insertResult:", insertResult)
		log.Println(err)
		return err
	}
	return nil
}

// UpdateTag update Tag
func UpdateTag(tag Tag) (int64, error) {
	collection := mdb.FiberDB.Collection(TableTag)
	filter := bson.D{{"_id", tag.ID}}
	opts := options.Replace().SetUpsert(true)
	updateResult, err := collection.ReplaceOne(context.Background(), filter, tag, opts)
	if err != nil {
		fmt.Println("Error update tag with updateResult:", updateResult)
		log.Println(err)
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

// DeleteTag delete Tag
func DeleteTag(id primitive.ObjectID) (int64, error) {
	collection := mdb.FiberDB.Collection(TableTag)
	filter := bson.D{{"_id", id}}
	// specify the SetCollation option to provide a collation that will ignore case for string comparisons
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	deleteResult, err := collection.DeleteOne(context.Background(), filter, opts)
	if err != nil {
		fmt.Println("Error delete tag with deleteResult:", deleteResult)
		log.Println(err)
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

// GetTag get tag
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#pkg-overview
func GetTag(id primitive.ObjectID) *Tag {
	var tag Tag
	collection := mdb.FiberDB.Collection(TableTag)
	filter := bson.D{{"_id", id}}
	err := collection.FindOne(context.Background(), filter).Decode(&tag)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &tag
}

// GetAllTag get all tag
func GetAllTag() []Tag {
	var tags []Tag
	collection := mdb.FiberDB.Collection(TableTag)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var tag Tag
		err := cur.Decode(&tag)
		if err != nil {
			log.Println(err)
		}
		// do something with result...
		tags = append(tags, tag)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}
	return tags
}

// GetTotalTag get total tag
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#pkg-examples
func GetTotalTag() (int64, error) {
	collection := mdb.FiberDB.Collection(TableTag)
	filter := bson.D{}
	// specify the MaxTime option to limit the amount of time the operation can run on the server
	opts := options.Count().SetMaxTime(10 * time.Second)
	count, err := collection.CountDocuments(context.Background(), filter, opts)
	return count, err
}

// GetSlideTag get slide tag
func GetSlideTag(skip int64, limlit int64) []Tag {
	var tags []Tag
	collection := mdb.FiberDB.Collection(TableTag)
	filter := bson.D{}
	// specify the Sort option to sort the returned documents by _id in Descending order
	opts := options.Find().SetSkip(skip).SetLimit(limlit).SetSort(bson.D{{"_id", -1}})
	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var tag Tag
		err := cur.Decode(&tag)
		if err != nil {
			log.Println(err)
		}
		// do something with result...
		tags = append(tags, tag)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}
	return tags
}
