package infrastructure

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Find(ctx context.Context, filters interface{}) (cur *mongo.Cursor, err error)
	FindOne(ctx context.Context, filters interface{}, options *options.FindOneOptions) (res *mongo.SingleResult, err error)
	InsertOne(ctx context.Context, data interface{}) (res *mongo.InsertOneResult, err error)
	UpdateOne(filters interface{}, options *options.UpdateOptions, data interface{}) (res *mongo.UpdateResult, err error)
	ReplaceOne(filters interface{}, options *options.ReplaceOptions, data interface{}) (res *mongo.UpdateResult, err error)
}

type repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewRepo(collection string) *repository {
	repository := new(repository)
	repository.client = NewDatastore()
	repository.collection = repository.client.Database("mainframe").Collection(collection)
	return repository
}

func (r *repository) FindOne(ctx context.Context, filters interface{}, options *options.FindOneOptions) (res *mongo.SingleResult, err error) {
	cur := r.collection.FindOne(context.TODO(), filters, options)
	if cur.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	if cur.Err() != nil {
		panic(cur.Err())
	}
	return cur, nil
}

func (r *repository) Find(
	ctx context.Context,
	filters interface{},
	options *options.FindOptions) (cur *mongo.Cursor, err error) {
	cur, err = r.collection.Find(context.TODO(), filters, options)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, cur.Err()
	}
	return cur, nil
}

func (r *repository) InsertOne(ctx context.Context, data interface{}) (res *mongo.InsertOneResult, err error) {
	res, err = r.collection.InsertOne(context.TODO(), data)
	if err != nil {
		panic(err)
	}
	return res, err
}

func (r *repository) UpdateOne(filters interface{}, options *options.UpdateOptions, data interface{}) (res *mongo.UpdateResult, err error) {
	res, err = r.collection.UpdateOne(context.TODO(), filters, data, options)
	if err != nil {
		panic(err)
	}
	return res, err
}

func (r *repository) ReplaceOne(filters interface{}, options *options.ReplaceOptions, data interface{}) (res *mongo.UpdateResult, err error) {
	res, err = r.collection.ReplaceOne(context.TODO(), filters, data, options)
	if err != nil {
		panic(err)
	}
	return res, err
}
