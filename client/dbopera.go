package client

import (
	"time"
	"wanbu_data_upload_api/models"

	"github.com/astaxie/beego/orm"
)

func AddWanbuDataUploadRecord(m *User_walkdays_struct) (x *models.WanbuDataUploadRecord, err error) {

	var walksteps, walkdays int

	for _, v := range m.Walkdays {
		walksteps += v.Daydata
		walkdays++
	}

	x = &models.WanbuDataUploadRecord{
		Userid:     m.Uid,
		Walkdays:   int8(walkdays),
		Walksteps:  int32(walksteps),
		Uploadtime: m.Timestamp,
		Timestamp:  time.Now().Unix(),
	}
	o := orm.NewOrm()
	_, err = o.Insert(x)
	return x, err
}
