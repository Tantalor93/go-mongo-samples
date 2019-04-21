package seed

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
)

func SeedDb(col *mongo.Collection, ctx context.Context) {

	purgeDbCollection(col, ctx)

	bytes, err := ioutil.ReadFile("seed.json")
	var documents []interface{}
	err = json.Unmarshal(bytes, &documents)
	if err != nil {
		panic("error while parsing seed file: " + err.Error())
	}
	log.WithFields(getLogContext(ctx)).Info("seeding collection")
	result, err := col.InsertMany(ctx, documents)
	if err != nil {
		panic("error while seeding mongo database: " + err.Error())
	}
	log.WithField("seeded_document_ids", result.InsertedIDs).
		WithFields(getLogContext(ctx)).
		Info("seeding of DB successful")
}

func purgeDbCollection(col *mongo.Collection, ctx context.Context) {
	deleteResult, err := col.DeleteMany(ctx, bson.D{}, nil)
	if err != nil {
		panic("error cleaning mongo DB")
	}
	log.WithField("deleted_count", deleteResult.DeletedCount).
		WithFields(getLogContext(ctx)).
		Info("purged collection")
}


func getLogContext(ctx context.Context) log.Fields {
	return ctx.Value("log").(log.Fields)
}

