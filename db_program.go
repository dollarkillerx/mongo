/**
 * @Author: DollarKiller
 * @Description: 常用解决方案 只池化到DB
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 17:31 2019-10-28
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
)

// 初始化db池
func (d *Db) initDbPool(dbName string) {
	key := Sha1Encode(dbName)

	_, ok := d.dbPool.Load(key)
	// 如果不存在 再执行此处
	if !ok {
		pool := NewObjPoll(func() interface{} {
			client, e := mongo.Connect(context.TODO(), options.Client().ApplyURI(d.uri))
			if e != nil {
				panic(e)
			}
			return client.Database(dbName)
		}, d.maxOpen)
		d.dbPool.Store(key, pool)
	}

	_, ok = d.dbTemporary.Load(key)
	if !ok {
		pool := &sync.Pool{
			New: func() interface{} {
				client, e := mongo.Connect(context.TODO(), options.Client().ApplyURI(d.uri))
				if e != nil {
					panic(e)
				}
				return client.Database(dbName)
			},
		}
		d.dbTemporary.Store(key, pool)
	}
}

// 从对象池获取collection
func (d *Db) getDbPool(dbName string) (*mongo.Database, error) {
	// 初始化pool
	d.initDbPool(dbName)
	key := Sha1Encode(dbName)

	value, ok := d.dbPool.Load(key)
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
	database, ok := getObj.(*mongo.Database)
	if !ok {
		e := errors.New("转换类型失败")
		clog.PrintWa(e)
		return nil, e
	}
	return database, nil
}

// 放回对象池
func (d *Db) pulDbPool(dbName string, db *mongo.Database) error {
	// 初始化pool
	d.initDbPool(dbName)
	key := Sha1Encode(dbName)

	value, ok := d.dbPool.Load(key)
	if !ok {
		e := errors.New("not data")
		clog.PrintWa(e)
		return e
	}
	obj, ok := value.(*ObjPool)
	if !ok {
		e := errors.New("转换类型失败")
		clog.PrintWa(e)
		return e
	}

	return obj.Release(db)
}

// 从临时对象池中获取
func (d *Db) getDbTemporary(dbName string) (*mongo.Database, error) {
	// 初始化pool
	d.initDbPool(dbName)
	key := Sha1Encode(dbName)

	value, ok := d.dbTemporary.Load(key)
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
	get := obj.Get()
	database, ok := get.(*mongo.Database)
	if !ok {
		e := errors.New("转换类型失败")
		clog.PrintWa(e)
		return nil, e
	}
	return database, nil
}

// 放回临时对象池
func (d *Db) pulDbTemporary(dbName string, db *mongo.Database) error {
	// 初始化pool
	d.initDbPool(dbName)
	key := Sha1Encode(dbName)

	value, ok := d.dbTemporary.Load(key)
	if !ok {
		e := errors.New("not data")
		clog.PrintWa(e)
		return e
	}
	obj, ok := value.(*sync.Pool)
	if !ok {
		e := errors.New("转换类型失败")
		clog.PrintWa(e)
		return e
	}

	obj.Put(db)
	return nil
}

// 进一步封装
func (d *Db) getDatabase(dbName string) (database *mongo.Database, ResultDb ResultDbPul, err error) {
	data, err := d.getDbPool(dbName)
	if err == nil {
		ResultDb = ResultDbPul{
			dbName:   dbName,
			database: data,
			tag:      1,
		}
		return data, ResultDb, err
	} else {
		temporary, err := d.getDbTemporary(dbName)
		ResultDb = ResultDbPul{
			dbName:   dbName,
			database: data,
			tag:      2,
		}
		return temporary, ResultDb, err
	}
}

// 放回的进一步封装
func (d *Db) pulDatabase(pul ResultDbPul) error {
	switch pul.tag {
	case 1:
		return d.pulDbPool(pul.dbName, pul.database)
	case 2:
		return d.pulDbTemporary(pul.dbName, pul.database)
	default:
		err := fmt.Errorf("tag 参数错误")
		clog.PrintWa(err)
		return err
	}
}

// 暴露
func (d *Db) GetDatabase(dbName string) (database *mongo.Database, ResultDb ResultDbPul, err error) {
	return d.getDatabase(dbName)
}

func (d *Db) PulDatabase(pul ResultDbPul) error {
	return d.pulDatabase(pul)
}
