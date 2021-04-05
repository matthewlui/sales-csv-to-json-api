package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SalesRecordDB interface {
	InsertMany([]interface{}) (*mongo.InsertManyResult, error)
}

type salesRecordDB struct {
	collection *mongo.Collection
}

func NewSalesRecordDB(dbEndPoint string, db string, collection string) SalesRecordDB {
	clientOptions := options.Client().ApplyURI(dbEndPoint)
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	return &salesRecordDB{collection: client.Database(db).Collection(collection)}
}

func (md *salesRecordDB) InsertMany(documents []interface{}) (*mongo.InsertManyResult, error) {
	insertManyResult, err := md.collection.InsertMany(context.TODO(), documents)
	return insertManyResult, err
}
