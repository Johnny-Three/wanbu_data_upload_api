package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type WanbuDataUploadRecord struct {
	Id       int   `orm:"column(id);auto"`
	Touserid int   `orm:"column(touserid)"`
	Daynum   int32 `orm:"column(daynum)"`
	Stepnum  int32 `orm:"column(stepnum)"`
	Dateline int64 `orm:"column(dateline)"`
	Devtype  int8  `orm:"column(devtype)"`
}

func (t *WanbuDataUploadRecord) TableName() string {
	return "wanbu_pm_personalupload"
}

func init() {
	orm.RegisterModel(new(WanbuDataUploadRecord))
}

// AddWanbuDataUploadRecord insert a new WanbuDataUploadRecord into database and returns
// last inserted Id on success.
func AddWanbuDataUploadRecord(m *WanbuDataUploadRecord) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetWanbuDataUploadRecordById retrieves WanbuDataUploadRecord by Id. Returns error if
// Id doesn't exist
func GetWanbuDataUploadRecordById(id int, min, max int64) (v int, err error) {

	//var rs orm.RawSeter
	o := orm.NewOrm()
	var totalItem int = 0
	var qs string
	if min == 0 && max == 0 {
		qs = fmt.Sprintf("SELECT count(*) FROM wanbu_pm_personalupload where touserid=%d", id)
	} else {
		qs = fmt.Sprintf("SELECT count(*) FROM wanbu_pm_personalupload where touserid=%d and dateline>=%d and dateline <=%d", id, min, max)
	}
	//总条数,总页数
	o.Raw(qs).QueryRow(&totalItem) //获取总条数
	return totalItem, nil
}

// GetAllWanbuDataUploadRecord retrieves all WanbuDataUploadRecord matches certain condition. Returns empty list if
// no records exist
func GetAllWanbuDataUploadRecord(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(WanbuDataUploadRecord))
	// query k=v
	for k, v := range query {

		if k == "dateline" {

			//如果查找到dateline这个特殊query，查某个时间段的时间，必须符合格式，否则返回格式错误
			ts := strings.Split(v, "-")
			if len(ts) != 2 {
				return nil, errors.New("query dateline 格式错误")
			}
			//从字符串转为时间戳，第一个参数是格式，第二个是要转换的时间字符串
			t1, err := time.Parse("20060102", ts[0])
			if err != nil {
				return nil, err
			}
			//从字符串转为时间戳，第一个参数是格式，第二个是要转换的时间字符串
			t2, err := time.Parse("20060102", ts[1])
			if err != nil {
				return nil, err
			}
			qs = qs.Filter("dateline__gte", t1.Unix()) // dateline >= t1
			qs = qs.Filter("dateline__lte", t2.Unix()) // dateline <= t2

		} else {

			// rewrite dot-notation to Object__Attribute
			k = strings.Replace(k, ".", "__", -1)
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []WanbuDataUploadRecord
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateWanbuDataUploadRecord updates WanbuDataUploadRecord by Id and returns error if
// the record to be updated doesn't exist
func UpdateWanbuDataUploadRecordById(m *WanbuDataUploadRecord) (err error) {
	o := orm.NewOrm()
	v := WanbuDataUploadRecord{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteWanbuDataUploadRecord deletes WanbuDataUploadRecord by Id and returns error if
// the record to be deleted doesn't exist
func DeleteWanbuDataUploadRecord(id int) (err error) {
	o := orm.NewOrm()
	v := WanbuDataUploadRecord{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&WanbuDataUploadRecord{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
