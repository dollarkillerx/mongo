# mongo
更人性化的mongo操作  (伪ORM)

### 入门
- [安装](#安装)
- [快速开始](#快速开始)

### 安装
``` 
go get github.com/dollarkillerx/mongo
```

### 快速开始
``` 
package main

import (
	"github.com/dollarkillerx/mongo"
	"log"
)

func main() {
	db, e := mongo.Open("127.0.0.1:27017")
	if e != nil {
		panic(e)
	}
	e = db.Ping()
	if e != nil {
		panic(e)
	}
	log.Println("链接成功 200ok!")
}
```

内部做了池化 效率提升

你可以使用我们提供的方法,也可以在池中获取  记得要放回蛤

### 对外暴露的接口 用完要记得放回
``` 
import (
	"github.com/dollarkillerx/mongo"
	"log"
	"time"
)

func main() {
	db, e := mongo.Open("mongodb://127.0.0.1:27017")
	if e != nil {
		panic(e)
	}
	e = db.Ping()
	if e != nil {
		panic(e)
	}
	log.Println("链接成功 200ok!")

	db.SetConnMaxLifetime(300 * time.Millisecond) // 设置超时时间  默认800 * time.Millisecond
	db.SetMaxOpenConn(10)                         // 设置最大链接数 默认2

	// 获取一个collection的链接   (本插件会池化collection,采用双池设计保证每次都能获取到)
	collection, resultPul, e := db.GetCollection("mongo1", "Coc")
	if e != nil {
		panic(e)
	} else {
		log.Println("获取成功")
	}
	collection = collection

	e = db.PulCollection(resultPul)
	if e != nil {
		panic(e)
	} else {
		log.Println("放回成功")
	}
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
	log.SetFlags(log.Llongfile | log.LstdFlags)
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
	db.SetMaxOpenConn(5)
	db.SetConnMaxLifetime(400 * time.Millisecond)

	// 创建读取
	collection := db.New("test", "os")
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

}
```