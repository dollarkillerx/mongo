/**
 * @Author: DollarKiller
 * @Description: main
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 20:15 2019-10-27
 */
package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dollarkillerx/mongo/clog"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo/options"
	"sync"
	"time"
)

// db
type Db struct {
	sync.RWMutex
	init                sync.Once
	uri                 string
	timeOut             time.Duration
	maxOpen             int
	notLimitedOpen      bool     // 不限制打开数量
	collectionPool      sync.Map // 存储collectionPool的
	collectionTemporary sync.Map // 存储临时对象
}

// 初始化db
func Open(uri string) (*Db, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Db{
		uri:     uri,
		timeOut: 800 * time.Millisecond,
		maxOpen: 2,
	}, nil
}

// 设置超时时间
func (d *Db) SetConnMaxLifetime(time time.Duration) {
	d.Lock()
	defer d.Unlock()
	d.timeOut = time
}

// new
func (d *Db) New(dbName, collectionName string) *Collection {
	return &Collection{
		db:         d,
		dbName:     dbName,
		collection: collectionName,
	}
}

// 设置最大打开数目
func (d *Db) SetMaxOpenConn(maxOpen int) {
	d.Lock()
	defer d.Unlock()
	d.maxOpen = maxOpen
}

// 设置不限制打开数量
func (d *Db) SetNotLimitedOpen(open bool) {
	d.Lock()
	defer d.Unlock()
	d.notLimitedOpen = open
}

// 链接检测数据库可用性
func (d *Db) Ping() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return err
	}
	return client.Ping(context.TODO(), nil)
}

/**
新的池化思路
看了下mongo-driver的源码  每一个collection都包含了一个db client的实例
return &Collection{
	client:         coll.client,
	db:             coll.db,
	name:           coll.name,
	readConcern:    coll.readConcern,
	writeConcern:   coll.writeConcern,
	readPreference: coll.readPreference,
	readSelector:   coll.readSelector,
	writeSelector:  coll.writeSelector,
	registry:       coll.registry,
}
新的思路就是 collection级的池化,每一个collection都单独拥有一个db client 实例
*/

// 初始化collection pool 内部调用使用
func (d *Db) initCollection(dbName, collectionName string) {
	key := Sha1Encode(dbName + collectionName)

	_, ok := d.collectionPool.Load(key)
	if !ok {
		// 如果对象池中没有 就new 一个
		poll := NewObjPoll(func() interface{} {
			client, e := mongo.Connect(context.Background(), options.Client().ApplyURI(d.uri))
			if e != nil {
				panic(e)
			}
			collection := client.Database(dbName).Collection(collectionName)

			return collection
		}, d.maxOpen)
		d.collectionPool.Store(key, poll)
	}

	_, ok = d.collectionTemporary.Load(key)
	if !ok {
		poll := sync.Pool{
			New: func() interface{} {
				client, e := mongo.Connect(context.Background(), options.Client().ApplyURI(d.uri))
				if e != nil {
					panic(e)
				}
				collection := client.Database(dbName).Collection(collectionName)

				return collection
			},
		}
		d.collectionTemporary.Store(key, poll)
	}
}

// 从对象池获取collection
func (d *Db) getCollectionPool(dbName, collectionName string) (*mongo.Collection, error) {
	key := Sha1Encode(dbName + collectionName)
	d.initCollection(dbName, collectionName)
	var collection *mongo.Collection
	value, ok := d.collectionPool.Load(key)
	if !ok {
		e := errors.New("not data")
		clog.PrintWa(e)
		return nil, e
	}
	obj, ok := value.(*ObjPool)
	if !ok {
		e := errors.New("转换类型失败")
		clog.PrintWa(e)
		return nil, e
	}
	getObj, e := obj.GetObj(d.timeOut)
	if e != nil {
		clog.PrintWa(e)
		return nil, e
	}
	collection, ok = getObj.(*mongo.Collection)
	if !ok {
		e := errors.New("转换类型失败")
		clog.PrintWa(e)
		return nil, e
	}
	return collection, nil
}

// 放回对象池
func (d *Db) pulCollectionPool(dbName, collectionName string, data *mongo.Collection) error {
	key := Sha1Encode(dbName + collectionName)
	value, ok := d.collectionPool.Load(key)
	if !ok {
		e := errors.New("not data")
		clog.PrintWa(e)
		return e
	}
	pol, ok := value.(*ObjPool)
	if !ok {
		e := errors.New("转换类型失败")
		clog.PrintWa(e)
		return e
	}
	return pol.Release(data)
}

// 从临时对象池获取
func (d *Db) getCollectionTemporary(dbName, collectionName string) (*mongo.Collection, error) {
	key := Sha1Encode(dbName + collectionName)
	d.initCollection(dbName, collectionName)
	var collection *mongo.Collection
	value, ok := d.collectionTemporary.Load(key)
	if !ok {
		e := errors.New("not data")
		clog.PrintWa(e)
		return nil, e
	}

	obj, ok := value.(*sync.Pool)
	if !ok {
		e := errors.New("转换类型失败")
		clog.PrintWa(e)
		return nil, e
	}
	getObj := obj.Get()
	collection, ok = getObj.(*mongo.Collection)
	if !ok {
		e := errors.New("转换类型失败")
		clog.PrintWa(e)
		return nil, e
	}
	return collection, nil
}

// 放回临时对象池
func (d *Db) pulCollectionTemporary(dbName, collectionName string, data *mongo.Collection) error {
	key := Sha1Encode(dbName + collectionName)
	value, ok := d.collectionTemporary.Load(key)
	if !ok {
		e := errors.New("not data")
		clog.PrintWa(e)
		return e
	}
	pol, ok := value.(*sync.Pool)
	if !ok {
		e := errors.New("转换类型失败")
		clog.PrintWa(e)
		return e
	}
	pol.Put(data)
	return nil
}

// 获取的进一步封装
// resultPul 放回
func (d *Db) getCollection(dbName, collectionName string) (collection *mongo.Collection, resultPul *ResultPul, err error) {
	collection, err = d.getCollectionPool(dbName, collectionName)
	if err == nil {
		resultPul = &ResultPul{
			dbName:         dbName,
			collectionName: collectionName,
			collection:     collection,
			tag:            1,
		}
		return collection, resultPul, err
	} else {
		// 向临时对象池中获取
		collection, err = d.getCollectionTemporary(dbName, collectionName)
		resultPul = &ResultPul{
			dbName:         dbName,
			collectionName: collectionName,
			collection:     collection,
			tag:            2,
		}
		return collection, resultPul, err
	}
}

// 放回的进一步封装
func (d *Db) pulCollection(resultPul *ResultPul) error {
	switch resultPul.tag {
	case 1:
		err := d.pulCollectionPool(resultPul.dbName, resultPul.collectionName, resultPul.collection)
		return err
	case 2:
		err := d.pulCollectionTemporary(resultPul.dbName, resultPul.collectionName, resultPul.collection)
		return err
	default:
		err := fmt.Errorf("tag 参数错误")
		clog.PrintWa(err)
		return err
	}
}

// 暴露出去的 获取和放回
func (d *Db) GetCollection(dbName, collectionName string) (collection *mongo.Collection, resultPul *ResultPul, err error) {
	return d.getCollection(dbName, collectionName)
}

func (d *Db) PulCollection(resultPul *ResultPul) error {
	return d.pulCollection(resultPul)
}

// 一下是过气代码 重构
//// 返回Db 如果你想做更多自义定的事情
//func (d *Db) Db() *mongo.Client {
//	return d.getDbByTemporary()
//}
//
//// 从对象次获取数据
//func (d *Db) getDbByPool(timeOut time.Duration) (*mongo.Client, error) {
//	// 初始化对象池
//	d.initPool()
//	obj, e := d.dbPool.GetObj(timeOut)
//	if e != nil {
//		return nil, e
//	}
//	return obj.(*mongo.Client), nil
//}
//
//// 放回对象池
//func (d *Db) putDbByPool(mongo interface{}) error {
//	return d.dbPool.Release(mongo)
//}
//
//// 从临时对象池获取
//func (d *Db) getDbByTemporary() *mongo.Client {
//	// 初始化对象池
//	return d.dbTemporary.Get().(*mongo.Client)
//}
//
//// 放回临时对象池
//func (d *Db) putDbByTemporary(mongo interface{}) {
//	d.dbTemporary.Put(mongo)
//}
//
//// collection 返回
//func (d *Db) NewCollection(dbName, collection string) *Collection {
//	return &Collection{dbName: dbName, collection: collection, db: d}
//}
//
//// 初始化 对象池
//func (d *Db) initPool() {
//	d.Lock()
//	openNum := d.maxOpen
//	d.Unlock()
//	d.init.Do(func() {
//		// 初始化对对象池
//		d.dbPool = NewObjPoll(func() interface{} {
//			ctx, _ := context.WithTimeout(context.Background(), d.timeOut)
//			client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
//			if err != nil {
//				panic(err)
//			}
//			return client
//		}, openNum)
//		// 初始化临时对象次
//		d.dbTemporary = &sync.Pool{
//			New: func() interface{} {
//				ctx, _ := context.WithTimeout(context.Background(), d.timeOut)
//				client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
//				if err != nil {
//					panic(err)
//				}
//				return client
//			},
//		}
//	})
//}
// 初始化 collection 池
//func (c *Db) initCollectionPool(dbName, collectionName string) {
//	c.Lock()
//	openNum := c.maxOpen
//	c.Unlock()
//	// 创建key
//	key := Sha1Encode(dbName + collectionName)
//	// 初始化对象池
//	// 查询collectionPool中key是否存在
//	_, ok := c.collectionPool.Load(key)
//	if !ok {
//		// 如果不存在 就 创建
//		poll := NewObjPoll(func() interface{} {
//			var client *mongo.Client
//			var err error
//			client, err = c.getDbByPool(time.Millisecond * 200)
//			if err != nil {
//				client = c.getDbByTemporary()
//				defer c.putDbByTemporary(client)
//			} else {
//				defer c.putDbByPool(client)
//			}
//
//			collection := client.Database(dbName).Collection(collectionName)
//			return collection
//		}, openNum)
//
//		c.collectionPool.Store(key, poll)
//	}
//	// 初始化临时对象池
//	_, b := c.collectionTemporary.Load(key)
//	if !b {
//		temporary := &sync.Pool{
//			New: func() interface{} {
//				var client *mongo.Client
//				var err error
//				client, err = c.getDbByPool(time.Millisecond * 200)
//				if err != nil {
//					client = c.getDbByTemporary()
//					defer c.putDbByTemporary(client)
//				} else {
//					defer c.putDbByPool(client)
//				}
//
//				collection := client.Database(dbName).Collection(collectionName)
//				return collection
//			},
//		}
//		c.collectionTemporary.Store(key, temporary)
//	}
//}
//
//// 获取collection pool 内容
//func (d *Db) getCollectionPool(dbName, collectionName string) (*mongo.Collection, error) {
//	key := Sha1Encode(dbName + collectionName)
//	// 初始化pool
//	d.initCollectionPool(dbName, collectionName)
//	value, ok := d.collectionPool.Load(key)
//	if !ok {
//		err := errors.New("不存在")
//		clog.PrintWa(err)
//		return nil, err
//	}
//	collectionPool := value.(*ObjPool)
//	collectionInterface, err := collectionPool.GetObj(d.timeOut)
//	if err != nil {
//		// 如果pool里面为空
//		return nil, err
//	}
//	return collectionInterface.(*mongo.Collection), nil
//}
//
//// 放回collection pool
//func (d *Db) pulCollectionPool(dbName, collectionName string, collection *mongo.Collection) error {
//	key := Sha1Encode(dbName + collectionName)
//
//	value, ok := d.collectionPool.Load(key)
//	if !ok {
//		err := errors.New("不存在")
//		clog.PrintWa(err)
//		return err
//	}
//	collectionPool := value.(*ObjPool)
//	return collectionPool.Release(collection)
//}
//
//// 从临时对象池中获取
//func (d *Db) getCollectionTemporary(dbName, collectionName string) (*mongo.Collection, error) {
//	key := Sha1Encode(dbName + collectionName)
//	// 初始化pool
//	d.initCollectionPool(dbName, collectionName)
//	value, ok := d.collectionTemporary.Load(key)
//	if !ok {
//		err := errors.New("不存在")
//		clog.PrintWa(err)
//		return nil, err
//	}
//	collectionPool := value.(*sync.Pool)
//	collection := collectionPool.Get()
//	return collection.(*mongo.Collection), nil
//}
//
////放回临时对象池
//func (d *Db) putCollectionTemporary(dbName, collectionName string, data *mongo.Collection) error {
//	key := Sha1Encode(dbName + collectionName)
//	// 初始化pool
//	d.initCollectionPool(dbName, collectionName)
//	value, ok := d.collectionTemporary.Load(key)
//	if !ok {
//		err := errors.New("不存在")
//		clog.PrintWa(err)
//		return err
//	}
//	collectionPool := value.(*sync.Pool)
//	collectionPool.Put(data)
//	return nil
//}
