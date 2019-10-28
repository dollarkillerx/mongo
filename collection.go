/**
 * @Author: DollarKiller
 * @Description: collection
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 22:15 2019-10-27
 */
package mongo

import (
	"context"
	"github.com/dollarkillerx/mongo/clog"
	"github.com/dollarkillerx/mongo/mongo-driver/bson"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo/options"
	"github.com/dollarkillerx/mongo/mongo-driver/x/bsonx"
)

type Collection struct {
	db         *Db
	dbName     string // 数据库名称
	collection string // 文档
}

// 清空数据库
func (c *Collection) Drop() error {
	collection, resultPul, err := c.db.getCollection(c.dbName, c.collection)
	if err != nil {
		return err
	} else {
		defer c.db.pulCollection(resultPul)
	}
	return collection.Drop(context.TODO())
}

// 设置索引
func (c *Collection) CreatIndex(key string, initial int32) (string, error) {
	collection, resultPul, err := c.db.getCollection(c.dbName, c.collection)
	if err != nil {
		return "", err
	} else {
		defer c.db.pulCollection(resultPul)
	}

	// 设置索引
	idx := mongo.IndexModel{
		Keys:    bsonx.Doc{{key, bsonx.Int32(initial)}},
		Options: options.Index().SetUnique(true),
	}

	return collection.Indexes().CreateOne(context.TODO(), idx)
}

// 插入一条数据
func (c *Collection) InsertOne(data interface{}) (*mongo.InsertOneResult, error) {
	collection, resultPul, err := c.db.getCollection(c.dbName, c.collection)
	if err != nil {
		return nil, err
	} else {
		defer c.db.pulCollection(resultPul)
	}
	return collection.InsertOne(context.TODO(), data)
}

// 插件多条数据
func (c *Collection) InsertMany(data []interface{}) (*mongo.InsertManyResult, error) {
	collection, resultPul, err := c.db.getCollection(c.dbName, c.collection)
	if err != nil {
		return nil, err
	} else {
		defer c.db.pulCollection(resultPul)
	}
	return collection.InsertMany(context.TODO(), data)
}

// 获取总数
func (c *Collection) CountDocuments() (int64, error) {
	collection, resultPul, err := c.db.getCollection(c.dbName, c.collection)
	if err != nil {
		return 0, err
	} else {
		defer c.db.pulCollection(resultPul)
	}
	return collection.CountDocuments(context.TODO(), bson.D{})
}

// 查询单条数据
func (c *Collection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	collection, resultPul, err := c.db.getCollection(c.dbName, c.collection)
	if err != nil {
		clog.PrintWa(err)
		return nil
	} else {
		defer c.db.pulCollection(resultPul)
	}
	return collection.FindOne(ctx, filter, opts...)
}

// 查询多条数据 方式一
func (c *Collection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	collection, resultPul, err := c.db.getCollection(c.dbName, c.collection)
	if err != nil {
		clog.PrintWa(err)
		return nil, err
	} else {
		defer c.db.pulCollection(resultPul)
	}

	return collection.Find(ctx, filter, opts...)
}

// 修改数据
func (c *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	collection, resultPul, err := c.db.getCollection(c.dbName, c.collection)
	if err != nil {
		clog.PrintWa(err)
		return nil, err
	} else {
		defer c.db.pulCollection(resultPul)
	}

	return collection.UpdateOne(ctx, filter, update, opts...)
}

// 删除数据
func (c *Collection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	collection, resultPul, err := c.db.getCollection(c.dbName, c.collection)
	if err != nil {
		clog.PrintWa(err)
		return nil, err
	} else {
		defer c.db.pulCollection(resultPul)
	}

	return collection.DeleteOne(ctx, filter, opts...)
}
