/**
 * @Author: DollarKiller
 * @Description:
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 20:33 2019-10-27
 */
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
