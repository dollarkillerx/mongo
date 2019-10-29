# mongo
mongo数据库连接池 适合高并发场景

采用双池设计,保证每次都能获取到
### 入门
- [安装](#安装)
- [快速开始](#快速开始)

### 安装
``` 
go get github.com/dollarkillerx/mongo
```

### 快速开始
``` 
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
	db.SetMaxOpenConn(1) // 设置最大打开数量
	db.SetConnMaxLifetime(400 * time.Millisecond) // 设置超时时间

	database := db.Database("okp") // 设置数据库
	collection := database.Collection("BOOK") // 设置collection

	// 清空数据库和集合
	e = collection.Drop()
```

内部做了池化 效率提升

你可以使用我们提供的方法,也可以在池中获取  记得要放回蛤

### 对外暴露的接口 用完要记得放回
``` 
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
```

### 简单CURL 操作
使用本插件提供的方法  都是自动从池中获取 和 放回
``` 
package main

import (
	"context"
	"github.com/dollarkillerx/mongo"
	"github.com/dollarkillerx/mongo/clog"
	"github.com/dollarkillerx/mongo/mongo-driver/bson"
	"github.com/dollarkillerx/mongo/mongo-driver/bson/primitive"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo/options"
	"log"
	"time"
)

type Book struct {
	Id       primitive.ObjectID `bson:"_id"`
	Name     string
	Category string
	Weight   int
	Author   AuthorInfo
}

type AuthorInfo struct {
	Name    string
	Country string
}

const (
	categoryComputer = "计算机"
	categorySciFi    = "科幻"
	countryChina     = "中国"
	countryAmerica   = "美国"
)

var (
	books = []interface{}{
		&Book{
			Id:       primitive.NewObjectID(),
			Name:     "深入理解计算机操作系统",
			Category: categoryComputer,
			Weight:   1,
			Author: AuthorInfo{
				Name:    "兰德尔 E.布莱斯特",
				Country: countryAmerica,
			},
		},
		&Book{
			Id:       primitive.NewObjectID(),
			Name:     "深入理解Linux内核",
			Category: categoryComputer,
			Weight:   1,
			Author: AuthorInfo{
				Name:    "博韦，西斯特",
				Country: countryAmerica,
			},
		},
		&Book{
			Id:       primitive.NewObjectID(),
			Name:     "三体",
			Category: categorySciFi,
			Weight:   1,
			Author: AuthorInfo{
				Name:    "刘慈欣",
				Country: countryChina,
			},
		},
	}
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

	// 清空数据库和集合
	e = collection.Drop()

	idxRet, e := collection.CreatIndex("name", 1)
	if e != nil {
		panic(e)
	}
	log.Println("Collection.Indexes().CreateOne:", idxRet)

	// 插入一条数据
	insertOneResult, err := collection.InsertOne(context.TODO(), books[0])
	if err != nil {
		clog.PrintWa(err)
		log.Fatal()
	}

	log.Println("Collection.InsertOne:", insertOneResult.InsertedID)

	// 插入多条数据
	insertManyResult, err := collection.InsertMany(context.TODO(), books[1:])
	if err != nil {
		clog.PrintWa(err)
		log.Fatal(err)
	}
	log.Println("collection.InsertMany:", insertManyResult.InsertedIDs)

	// 获取数据总数
	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		clog.PrintWa(err)
		log.Fatal(count)
	}
	log.Println("collection.CountDocuments:", count)

	// 查询单挑数据
	var one Book
	e = collection.FindOne(context.Background(), bson.M{"name": "三体"}).Decode(&one)
	if e != nil {
		clog.PrintWa(err)
		log.Fatal(err)
	}
	log.Println("collection.FindOne: ", one)

	// 查询多条数据 方式一
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		clog.PrintWa(err)
		log.Fatal(err)
	}
	if err := cur.Err(); err != nil {
		clog.PrintWa(err)
		log.Fatal(err)
	}
	var all []*Book
	err = cur.All(context.TODO(), &all)
	if err != nil {
		clog.PrintWa(err)
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	log.Println("collection.Find curl.All", all)
	for _, one := range all {
		log.Println(one)
	}

	// 查询多条数据 方式二
	cur, err = collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var b Book
		if err = cur.Decode(&b); err != nil {
			log.Fatal(err)
		}
		log.Println("collection.Find cur.Next:", b)
	}
	cur.Close(context.TODO())

	// 模糊查询
	cur, err = collection.Find(context.TODO(), bson.M{"name": primitive.Regex{Pattern: "深入"}})
	if err != nil {
		log.Fatal(err)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var b Book
		if err = cur.Decode(&b); err != nil {
			log.Fatal(err)
		}
		log.Println("collection.Find name=primitive.Regex{深入}: ", b)
	}
	cur.Close(context.TODO())

	// 二级结构体查询
	// 二级结构体查询
	cur, err = collection.Find(context.Background(), bson.M{"author.country": countryChina})
	// cur, err = collection.Find(context.Background(), bson.D{bson.E{"author.country", countryChina}})
	if err != nil {
		log.Fatal(err)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.Background()) {
		var b Book
		if err = cur.Decode(&b); err != nil {
			log.Fatal(err)
		}
		log.Println("collection.Find author.country=", countryChina, ":", b)
	}
	cur.Close(context.Background())

	// 修改一条数据
	b1 := books[0].(*Book)
	b1.Weight = 2
	update := bson.M{"$set": b1}
	updateResult, err := collection.UpdateOne(context.Background(), bson.M{"name": b1.Name}, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection.UpdateOne:", updateResult)

	// 修改一条数据，如果不存在则插入
	new := &Book{
		Id:       primitive.NewObjectID(),
		Name:     "球状闪电",
		Category: categorySciFi,
		Author: AuthorInfo{
			Name:    "刘慈欣",
			Country: countryChina,
		},
	}
	update = bson.M{"$set": new}
	updateOpts := options.Update().SetUpsert(true)
	updateResult, err = collection.UpdateOne(context.Background(), bson.M{"_id": new.Id}, update, updateOpts)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection.UpdateOne:", updateResult)

	// 删除一条数据
	deleteResult, err := collection.DeleteOne(context.Background(), bson.M{"_id": new.Id})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection.DeleteOne:", deleteResult)
}
```