/**
 * @Author: DollarKiller
 * @Description:
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 10:29 2019-10-29
 */
package main

import (
	"context"
	"github.com/dollarkillerx/mongo"
	"github.com/dollarkillerx/mongo/clog"
	"log"
	"time"
)

func main() {
	uri := "mongodb://127.0.0.1:27017"
	db, e := mongo.Open(uri)
	if e != nil {
		panic(e)
	}

	e = db.Ping()
	if e != nil {
		panic(e)
	}
	clog.Println("MongoDb 链接成功")

	// 配置
	db.SetMaxOpenConn(1)
	db.SetConnMaxLifetime(400 * time.Millisecond)

	database := db.Database("okp")
	collection := database.Collection("BOOK")

	getCollection, pul, e := collection.GetCollection()
	if e == nil {
		defer func() {
			err := collection.PulCollection(pul)
			if err == nil {
				log.Println("放回成功")
			}
		}()
	}
	e = getCollection.Drop(context.TODO())
	if e == nil {
		log.Println("数据库清空成功200ok")
	}
}
