package action

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hedgehog-hids-server/common"
	"hedgehog-hids-server/log"
	"time"
)
type Student struct {
	Name string
	Age int
}

// ComputerInfoSave 保存client信息
func ComputerInfoSave(info common.ComputerInfo) {
	c := common.DB.Collection ("client")
	info.Uptime = time.Now()
	update := bson.M{"$set":info}
	updateOpts := options.Update().SetUpsert(true)
	_, err :=c.UpdateOne(context.TODO(),bson.M{"ip": info.IP}, update,updateOpts)
	if err != nil {
		log.Error(err)
	}else {
		log.Infof("ComputerInfo of %s update success",info.IP)
	}
}
