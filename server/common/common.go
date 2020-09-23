package common

import (
	"flag"
	"go.mongodb.org/mongo-driver/mongo"
	"hedgehog-hids-server/database"
	"hedgehog-hids-server/log"
	"os"
)

var (
	// DB 数据库连接池
	DB      *mongo.Database
	db		*string
	mongodb *string
	Config  serverConfig
	LocalIp string
	err     error
)

func init(){
	mongodb = flag.String("url","mongodb://localhost:27017","mongodb mongodb://admin:123456@localhost/")
	db = flag.String("dbname","agent","mongodb database")
	flag.Parse()
	//if len(os.Args) <= 2 {
	//	flag.PrintDefaults()
	//	os.Exit(1)
	//}

	DB,err = database.Conn(*mongodb,*db)
	if err != nil{
		log.Error(err)
		flag.PrintDefaults()
		os.Exit(1)
	}


}