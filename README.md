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