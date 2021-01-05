/**
 *
 * @author nghiatc
 * @since Jan 3, 2021
 */

package post

import (
	"context"
	"fmt"
	"log"
	"ntc-gfiber/mdb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Const
const (
	TablePost = "post"
)

// Post struct
type Post struct {
	// ID        primitive.ObjectID `bson:"_id" json:"id"`
	ID        int64     `bson:"_id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Body      string    `bson:"body" json:"body"`
	CreatedAt time.Time `bson:"created_at" json:"-"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// InsertPost Insert Post
func InsertPost(post Post) error {
	client := mdb.GetClient()
	defer mdb.Close(client)
	collection := client.Database(mdb.DbName).Collection(TablePost)
	insertResult, err := collection.InsertOne(context.Background(), post)
	if err != nil {
		fmt.Println("Error Inserted post with insertResult:", insertResult)
		log.Println(err)
		return err
	}
	return nil
}

// UpdatePost update Post
func UpdatePost(post Post) (int64, error) {
	client := mdb.GetClient()
	defer mdb.Close(client)
	collection := client.Database(mdb.DbName).Collection(TablePost)
	filter := bson.D{{"_id", post.ID}}
	opts := options.Replace().SetUpsert(true)
	updateResult, err := collection.ReplaceOne(context.Background(), filter, post, opts)
	if err != nil {
		fmt.Println("Error update post with updateResult:", updateResult)
		log.Println(err)
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

// DeletePost delete Post
func DeletePost(id int64) (int64, error) {
	client := mdb.GetClient()
	defer mdb.Close(client)
	collection := client.Database(mdb.DbName).Collection(TablePost)
	filter := bson.D{{"_id", id}}
	// specify the SetCollation option to provide a collation that will ignore case for string comparisons
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	deleteResult, err := collection.DeleteOne(context.Background(), filter, opts)
	if err != nil {
		fmt.Println("Error delete post with deleteResult:", deleteResult)
		log.Println(err)
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

// GetPost get post
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#pkg-overview
func GetPost(id int64) Post {
	var post Post
	client := mdb.GetClient()
	defer mdb.Close(client)
	collection := client.Database(mdb.DbName).Collection(TablePost)
	filter := bson.D{{"_id", id}}
	err := collection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		log.Println(err)
	}
	return post
}

// GetAllPost get all post
func GetAllPost() []Post {
	var posts []Post
	client := mdb.GetClient()
	defer mdb.Close(client)
	collection := client.Database(mdb.DbName).Collection(TablePost)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var post Post
		err := cur.Decode(&post)
		if err != nil {
			log.Println(err)
		}
		// do something with result...
		posts = append(posts, post)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}
	return posts
}

// GetTotalPost get total post
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#pkg-examples
func GetTotalPost() (int64, error) {
	client := mdb.GetClient()
	defer mdb.Close(client)
	collection := client.Database(mdb.DbName).Collection(TablePost)
	filter := bson.D{}
	// specify the MaxTime option to limit the amount of time the operation can run on the server
	opts := options.Count().SetMaxTime(10 * time.Second)
	count, err := collection.CountDocuments(context.Background(), filter, opts)
	return count, err
}

// GetSlidePost get slide post
func GetSlidePost(skip int64, limlit int64) []Post {
	var posts []Post
	client := mdb.GetClient()
	defer mdb.Close(client)
	collection := client.Database(mdb.DbName).Collection(TablePost)
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
		var post Post
		err := cur.Decode(&post)
		if err != nil {
			log.Println(err)
		}
		// do something with result...
		posts = append(posts, post)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}
	return posts
}
