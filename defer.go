/**
 * @Author: DollarKiller
 * @Description: defer 对一些对象的简化
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 13:54 2019-10-28
 */
package mongo

import "github.com/dollarkillerx/mongo/mongo-driver/mongo"

type ResultPul struct {
	dbName         string
	collectionName string
	collection     *mongo.Collection
	tag            int
}

type ResultDbPul struct {
	dbName   string
	database *mongo.Database
	tag      int
}
