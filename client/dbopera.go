package client

import (
	"wanbu_data_upload_api/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func AddWanbuDataUploadRecord(m *User_walkdays_struct) (x *models.WanbuDataUploadRecord, err error) {

	var walksteps, walkdays int

	for _, v := range m.Walkdays {
		walksteps += v.Daydata
		walkdays++
	}

	x = &models.WanbuDataUploadRecord{
		Touserid: m.Uid,
		Daynum:   int32(walkdays),
		Stepnum:  int32(walksteps),
		Dateline: m.Timestamp,
		Devtype:  int8(m.Devtype),
	}
	o := orm.NewOrm()
	_, err = o.Insert(x)

	beego.Trace("userid", x.Touserid, "insert ok")

	return x, err
}
