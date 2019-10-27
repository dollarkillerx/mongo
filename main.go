/**
 * @Author: DollarKiller
 * @Description: main
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 20:15 2019-10-27
 */
package mongo

import (
	"context"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo/options"
	"sync"
	"time"
)

// db
type Db struct {
	sync.RWMutex
	init           sync.Once
	db             *mongo.Client
	uri            string
	timeOut        time.Duration
	maxOpen        int
	dbPool         *ObjPool   // db对象池
	dbTemporary    *sync.Pool // db备用临时对象池
	notLimitedOpen bool       // 不限制打开数量
}

// 初始化db
func Open(uri string) (*Db, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Db{
		uri:     uri,
		db:      client,
		timeOut: 10 * time.Second,
		maxOpen: 2,
	}, nil
}

// 设置超时时间
func (d *Db) SetConnMaxLifetime(time time.Duration) {
	d.Lock()
	defer d.Unlock()
	d.timeOut = time
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
	d.Lock()
	defer d.Unlock()
	err := d.db.Ping(context.TODO(), nil)
	return err
}

// 返回Db 如果你想做更多自义定的事情
func (d *Db) Db() *mongo.Client {
	return d.getDbByTemporary()
}

// 从对象次获取数据
func (d *Db) getDbByPool(timeOut time.Duration) (*mongo.Client, error) {
	// 初始化对象池
	d.initPool()
	obj, e := d.dbPool.GetObj(timeOut)
	if e != nil {
		return nil, e
	}
	return obj.(*mongo.Client), nil
}

// 放回对象池
func (d *Db) putDbByPool(mongo interface{}) error {
	return d.dbPool.Release(mongo)
}

// 从临时对象池获取
func (d *Db) getDbByTemporary() *mongo.Client {
	// 初始化对象池
	d.initPool()
	return d.dbTemporary.Get().(*mongo.Client)
}

// 放回临时对象池
func (d *Db) putDbByTemporary(mongo interface{}) {
	d.dbTemporary.Put(mongo)
}

// collection 返回
func (d *Db) NewCollection(dbName, collection string) *Collection {
	return &Collection{dbName: dbName, collection: collection, db: d}
}

// 初始化 对象池
func (d *Db) initPool() {
	d.init.Do(func() {
		// 初始化对对象池
		d.dbPool = NewObjPoll(func() interface{} {
			ctx, _ := context.WithTimeout(context.Background(), d.timeOut)
			client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
			if err != nil {
				panic(err)
			}
			return client
		}, d.maxOpen)
		// 初始化临时对象次
		d.dbTemporary = &sync.Pool{
			New: func() interface{} {
				ctx, _ := context.WithTimeout(context.Background(), d.timeOut)
				client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
				if err != nil {
					panic(err)
				}
				return client
			},
		}
	})
}
