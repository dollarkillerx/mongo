/**
 * @Author: DollarKiller
 * @Description: collection
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 22:15 2019-10-27
 */
package mongo

import (
	"context"
	"github.com/dollarkillerx/mongo/mongo-driver/bson"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo/options"
	"github.com/dollarkillerx/mongo/mongo-driver/x/bsonx"
	"sync"
	"time"
)

type Collection struct {
	db             *Db
	dbName         string // 数据库名称
	collection     string // 文档
	collectionPool *sync.Pool
}

// 清空数据库
func (c *Collection) Drop() error {
	var client *mongo.Client
	var err error
	client, err = c.db.getDbByPool(time.Millisecond * 200)
	if err != nil {
		client = c.db.getDbByTemporary()
		defer c.db.putDbByTemporary(client)
	} else {
		defer c.db.putDbByPool(client)
	}

	collection := client.Database(c.dbName).Collection(c.collection)
	return collection.Drop(context.TODO())
}

// 设置索引
func (c *Collection) CreatIndex(key string, initial int32) (string, error) {
	var client *mongo.Client
	var err error
	client, err = c.db.getDbByPool(time.Millisecond * 200)
	if err != nil {
		client = c.db.getDbByTemporary()
		defer c.db.putDbByTemporary(client)
	} else {
		defer c.db.putDbByPool(client)
	}

	collection := client.Database(c.dbName).Collection(c.collection)

	// 设置索引
	idx := mongo.IndexModel{
		Keys:    bsonx.Doc{{key, bsonx.Int32(initial)}},
		Options: options.Index().SetUnique(true),
	}

	return collection.Indexes().CreateOne(context.TODO(), idx)
}

// 插入一条数据
func (c *Collection) InsertOne(data interface{}) (*mongo.InsertOneResult, error) {
	var client *mongo.Client
	var err error
	client, err = c.db.getDbByPool(time.Millisecond * 200)
	if err != nil {
		client = c.db.getDbByTemporary()
		defer c.db.putDbByTemporary(client)
	} else {
		defer c.db.putDbByPool(client)
	}

	collection := client.Database(c.dbName).Collection(c.collection)
	return collection.InsertOne(context.TODO(), data)
}

// 插件多条数据
func (c *Collection) InsertMany(data []interface{}) (*mongo.InsertManyResult, error) {
	var client *mongo.Client
	var err error
	client, err = c.db.getDbByPool(time.Millisecond * 200)
	if err != nil {
		client = c.db.getDbByTemporary()
		defer c.db.putDbByTemporary(client)
	} else {
		defer c.db.putDbByPool(client)
	}

	collection := client.Database(c.dbName).Collection(c.collection)
	return collection.InsertMany(context.TODO(), data)
}

// 获取总数
func (c *Collection) CountDocuments() (int64, error) {
	var client *mongo.Client
	var err error
	client, err = c.db.getDbByPool(time.Millisecond * 200)
	if err != nil {
		client = c.db.getDbByTemporary()
		defer c.db.putDbByTemporary(client)
	} else {
		defer c.db.putDbByPool(client)
	}

	collection := client.Database(c.dbName).Collection(c.collection)
	return collection.CountDocuments(context.TODO(), bson.D{})
}

// 查询单条数据
func (c *Collection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	var client *mongo.Client
	var err error
	client, err = c.db.getDbByPool(time.Millisecond * 200)
	if err != nil {
		client = c.db.getDbByTemporary()
		defer c.db.putDbByTemporary(client)
	} else {
		defer c.db.putDbByPool(client)
	}

	collection := client.Database(c.dbName).Collection(c.collection)
	return collection.FindOne(ctx, filter, opts...)
}

//