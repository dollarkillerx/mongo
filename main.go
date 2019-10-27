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

type Db struct {
	sync.RWMutex
	db      *mongo.Client
	host    string
	timeOut time.Duration
}

func Open(host string) (*Db, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, &options.ClientOptions{Hosts: []string{host}})
	if err != nil {
		return nil, err
	}
	return &Db{host: host, db: client, timeOut: 10 * time.Second}, nil
}

func (d *Db) SetTimeOut(time time.Duration) {
	d.timeOut = time
}

func (d *Db) Ping() error {
	d.Lock()
	defer d.Unlock()
	err := d.db.Ping(context.TODO(), nil)
	return err
}
