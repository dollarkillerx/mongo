/**
 * @Author: DollarKiller
 * @Description: db相关操作
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 09:31 2019-10-29
 */
package mongo

import (
	"context"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo/options"
	"github.com/dollarkillerx/mongo/mongo-driver/x/bsonx"
)

type Database struct {
	db     *Db
	dbName string
}

func (d *Db) Database(dbName string) *Database {
	return &Database{
		db:     d,
		dbName: dbName,
	}
}

type DbCollection struct {
	database   *Database
	collection string
	opts       []*options.CollectionOptions
}

func (d *Database) Collection(name string, opts ...*options.CollectionOptions) *DbCollection {
	return &DbCollection{
		database:   d,
		collection: name,
		opts:       opts,
	}
}

func (d *DbCollection) Drop() error {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}

	return database.Collection(d.collection, d.opts...).Drop(context.TODO())
}

func (d *DbCollection) Indexes() mongo.IndexView {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)

	return collection.Indexes()
}

// 设置索引
func (d *DbCollection) CreatIndex(key string, initial int32) (idxRet string,err error) {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)
	// 设置索引
	idx := mongo.IndexModel{
		Keys:    bsonx.Doc{{key, bsonx.Int32(initial)}},
		Options: options.Index().SetUnique(true),
	}

	return collection.Indexes().CreateOne(context.TODO(), idx)
}

func (d *DbCollection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)

	return collection.InsertOne(ctx, document, opts...)
}

func (d *DbCollection) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)

	return collection.InsertMany(ctx, documents, opts...)
}

func (d *DbCollection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)

	return collection.CountDocuments(ctx, filter, opts...)
}

func (d *DbCollection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)

	return collection.FindOne(ctx, filter, opts...)
}

func (d *DbCollection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)

	return collection.Find(ctx, filter, opts...)
}

func (d *DbCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)

	return collection.UpdateOne(ctx, filter, update, opts...)
}

func (d *DbCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)

	return collection.UpdateMany(ctx, filter, update, opts...)
}

func (d *DbCollection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)

	return collection.DeleteOne(ctx, filter, opts...)
}

func (d *DbCollection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err == nil {
		defer d.database.db.pulDatabase(ResultDb)
	}
	collection := database.Collection(d.collection, d.opts...)

	return collection.DeleteMany(ctx, filter, opts...)
}

// 暴露出去
func (d *DbCollection) GetCollection() (*mongo.Collection, *ResultDbPul, error) {
	database, ResultDb, err := d.database.db.getDatabase(d.database.dbName)
	if err != nil {
		return nil, nil, err
	}
	collection := database.Collection(d.collection, d.opts...)
	return collection, ResultDb, nil
}

// 放回
func (d *DbCollection) PulCollection(data *ResultDbPul) error {
	return d.database.db.pulDatabase(data)
}
