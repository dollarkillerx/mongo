/**
 * @Author: DollarKiller
 * @Description:
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 16:21 2019-10-28
 */
package main

import (
	"context"
	"fmt"
	"github.com/dollarkillerx/mongo"
	"github.com/dollarkillerx/mongo/clog"
	"github.com/dollarkillerx/mongo/mongo-driver/bson"
	"github.com/dollarkillerx/mongo/mongo-driver/bson/primitive"
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

	// 创建读取
	collection := db.NewCollection("test", "os")
	// 空气数据库
	e = collection.Drop()
	if e != nil {
		log.Fatalln(e)
	}

	// 设置索引
	idxRet, e := collection.CreatIndex("name", 1)
	if e != nil {
		log.Fatalln(e)
	}
	log.Println("Collection.Indexes().CreateOne:", idxRet)

	// 插入一条数据
	insertOneResult, err := collection.InsertOne(books[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Collection.InsertOne:", insertOneResult.InsertedID)

	// 插入多条数据
	insertManyResult, err := collection.InsertMany(books[1:])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection.InsertMany:", insertManyResult.InsertedIDs)

	// 获取数据总数
	count, err := collection.CountDocuments()
	if err != nil {
		log.Fatal(count)
	}
	log.Println("collection.CountDocuments:", count)

	// 查询单挑数据
	var one Book
	e = collection.FindOne(context.TODO(), bson.M{"name": "三体"}).Decode(&one)
	if e != nil {
		log.Fatal(err)
	}
	log.Println("collection.FindOne: ", one)

	// 查询多条数据 方式一
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	var all []*Book
	err = cur.All(context.TODO(), &all)
	if err != nil {
		log.Fatal(err)
	}
	e = cur.Close(context.TODO())
	if e != nil {
		panic(e)
	}

	log.Println("collection.Find curl.All", all)
	for _, one := range all {
		fmt.Println(one)
	}

	//// 查询多条数据 方式二
	//cur, err = collection.Find(context.Background(), bson.D{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if err := cur.Err(); err != nil {
	//	log.Fatal(err)
	//}
	//for cur.Next(context.TODO()) {
	//	var b Book
	//	if err = cur.Decode(&b); err != nil {
	//		log.Fatal(err)
	//	}
	//	clog.Println("collection.Find cur.Next:")
	//	clog.Println(b)
	//}
	//cur.Close(context.TODO())
}
