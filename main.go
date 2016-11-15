package main

import (
	"fmt"
	"os"
	"time"
	client "wanbu_data_upload_api/client"

	config "github.com/msbranco/goconfig"

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

	cf, _ := config.ReadConfigFile("../etc/config.ini")
	rdip1, _ := cf.GetString("DBCONN1", "IP")
	rdusr1, _ := cf.GetString("DBCONN1", "USERID")
	rdpwd1, _ := cf.GetString("DBCONN1", "USERPWD")
	rdname1, _ := cf.GetString("DBCONN1", "DBNAME")
	dbadress := rdusr1 + ":" + rdpwd1 + "@tcp(" + rdip1 + ")/" + rdname1
	orm.RegisterDataBase("default", "mysql", dbadress)

	nsqip, _ := cf.GetString("NSQ", "IP")
	nsqport, _ := cf.GetString("NSQ", "PORT")
	nsqadress = nsqip + ":" + nsqport
}

var consumer *nsq.Consumer
var nsqadress = "192.168.20.248:4161"
var version = "1.0.0PR1"

func main() {

	args := os.Args

	if len(args) == 2 && (args[1] == "-v") {

		fmt.Println("看好了兄弟，现在的版本是【", version, "】，可别弄错了")
		os.Exit(0)
	}

	//对接NSQ，消费上传消息
	consumer, err := client.NewConsummer("base_data_upload", "wanbudatauplaodapi")
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

	//正常流程
	go func() {
		for {
			select {

			case m := <-client.User_walk_data_chan:
				if v, err := client.AddWanbuDataUploadRecord(&m); err != nil {
					Logger.Critical(errors.New(fmt.Sprintf("AddWanbuDataUploadRecord error happens[%s]", err.Error())))
				} else {
					Logger.Infof("Userid:[%d]上传了[%d]天的数据，总步数为[%d]，上传时间为[%s]", v.Touserid, v.Daynum, v.Stepnum, convertTimeToString(v.Dateline))
				}

			default:
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()

}
