package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/tantalor93/go-mongo-samples/seed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})

	dbUri, db, collection := getParameters()

	logContext := log.Fields{
		"db_uri":     *dbUri,
		"collection": *collection,
		"db":         *db,
	}

	log.WithField("db_uri", dbUri).Info("connecting to mongo DB.")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ctx = context.WithValue(ctx, "log", logContext)

	client := createMongoClient(dbUri, ctx)

	col := client.Database(*db).Collection(*collection)

	seed.SeedDb(col, ctx)

	pipeline := []bson.M{
		{"$project": bson.M{"pages": 1, "language": 1}},
		{"$match": bson.M{"pages": bson.M{"$gt": 200}}},
		{"$group": bson.M{"_id": "$language", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}},
	}

	works := aggregate(ctx, col, pipeline)

	log.Info(works)
}

type AggResult struct {
	ID    string `bson:"_id"`
	Count int    `bson:count`
}

func aggregate(ctx context.Context, col *mongo.Collection, pipeline []bson.M) []AggResult {
	var result = make([]AggResult, 0)
	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		panic(err)
	}
	for cursor.Next(ctx) {
		var res AggResult
		decodeErr := cursor.Decode(&res)
		if decodeErr != nil {
			panic(decodeErr)
		}
		result = append(result, res)
	}

	return result
}

func createMongoClient(dbUri *string, ctx context.Context) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + *dbUri))
	if err != nil {
		panic("error creating mongo client")
	}

	err = client.Connect(ctx)

	if err != nil {
		panic("error while connecting to mongo: " + err.Error())
	}

	return client
}

func getParameters() (*string, *string, *string) {
	dbUri := flag.String("db_uri", "localhost:27017", "URL of DB, defaults to localhost:27017")
	db := flag.String("db", "local", "name of mongo database, defaults to local")
	collection := flag.String("collection", "samples", "collection to be seeded")
	flag.Parse()
	return dbUri, db, collection
}
