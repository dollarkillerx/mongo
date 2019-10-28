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