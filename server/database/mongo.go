package database

import (
	"context"
	"hedgehog-hids-server/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func Conn(url string,db string) (* mongo.Database,error){
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(url)
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	dbclient := client.Database(db)
	log.Info("Connected to MongoDB success!")
	return dbclient,nil
}

