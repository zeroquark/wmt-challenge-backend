package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
	ctx    context.Context
}

func GetDB() DB {
	uri := DefaultCredentials.GetURI()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	return DB{client, ctx}
}

func (db *DB) GetDefaultDatabase() *mongo.Database {
	return db.Client.Database(DefaultCredentials.database)
}

func (db *DB) GetDefaultCollection() *mongo.Collection {
	return db.GetDefaultDatabase().Collection(DefaultCredentials.collection)
}

func (db *DB) GetAllRecords() {
	collection := db.GetDefaultCollection()
	cursor, err := collection.Find(db.ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(db.ctx) {
		var product bson.M
		if err = cursor.Decode(&product); err != nil {
			log.Fatal(err)
		}
		fmt.Println(product)
	}
}

func (db *DB) FindById(id int) (*Product, error) {
	var result Product
	filter := bson.D{{Key: "id", Value: id}}

	err := db.GetDefaultCollection().FindOne(context.TODO(), filter).Decode(&result)
	if err != nil && result.Brand == "" {
		return nil, err
	}

	return &result, nil
}
func (db *DB) FindByBrand(brand string) ([]*Product, error) {
	// Collection of found documents (Product)
	var results []*Product
	// Set filter and find options
	filter := bson.D{{Key: "brand", Value: primitive.Regex{Pattern: brand}}}
	// findOptions := options.Find().SetLimit(3)
	// Set the cursor
	// cur, err := db.GetDefaultCollection().Find(context.TODO(), filter, findOptions)
	cur, err := db.GetDefaultCollection().Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var p Product
		err := cur.Decode(&p)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &p)
	}

	if err := cur.Err(); err != nil {
		log.Fatal()
		return results, err
	}
	cur.Close(context.TODO())

	return results, nil
}

func (db *DB) FindByDescription(description string) ([]*Product, error) {
	// Collection of found documents (Product)
	var results []*Product
	// Set filter and find options
	filter := bson.D{{Key: "description", Value: primitive.Regex{Pattern: description}}}
	// findOptions := options.Find().SetLimit(3)
	// Set the cursor
	// cur, err := db.GetDefaultCollection().Find(context.TODO(), filter, findOptions)
	cur, err := db.GetDefaultCollection().Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var p Product
		err := cur.Decode(&p)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &p)
	}

	if err := cur.Err(); err != nil {
		log.Fatal()
		return results, err
	}
	cur.Close(context.TODO())

	return results, nil
}

func (db *DB) Close() {
	if err := db.Client.Disconnect(db.ctx); err != nil {
		panic(err)
	}
}
