package main

import (
	"fmt"
	"time"
	client "wanbu_data_upload_api/client"
	. "wanbu_data_upload_api/logs"
	_ "wanbu_data_upload_api/routers"

	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	nsq "github.com/bitly/go-nsq"
	_ "github.com/go-sql-driver/mysql"
)

func nowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func convertTimeToString(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

func init() {

	orm.RegisterDataBase("default", "mysql", "heyu:wanbu12#$@tcp(101.201.141.236:3306)/wanbu")
}

var consumer *nsq.Consumer
var nsqadress = "192.168.20.249:4161"

func main() {

	//对接NSQ，消费上传消息
	consumer, err := client.NewConsummer("base_data_upload", "walkdatet1")
	if err != nil {
		panic(err)
	}

	//Consumer运行，消费消息..
	go func(consumer *nsq.Consumer) {

		err := client.ConsumerRun(consumer, "base_data_upload", nsqadress)
		if err != nil {
			panic(err)
		}
	}(consumer)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()

	//正常流程
	for {
		select {

		case m := <-client.User_walk_data_chan:
			if v, err := client.AddWanbuDataUploadRecord(&m); err != nil {
				Logger.Critical(errors.New(fmt.Sprintf("AddWanbuDataUploadRecord error happens[%s]", err.Error())))
			} else {
				Logger.Infof("Userid:[%d]上传了[%d]天的数据，总步数为[%d]，上传时间为[%s]", v.Userid, v.Walkdays, v.Walksteps, convertTimeToString(v.Uploadtime))
			}

		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}
