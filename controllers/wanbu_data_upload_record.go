package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
	"wanbu_data_upload_api/models"

	"github.com/astaxie/beego"
)

// WanbuDataUploadRecordController oprations for WanbuDataUploadRecord
type WanbuDataUploadRecordController struct {
	beego.Controller
}

// URLMapping ...
func (c *WanbuDataUploadRecordController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create WanbuDataUploadRecord
// @Param	body		body 	models.WanbuDataUploadRecord	true		"body for WanbuDataUploadRecord content"
// @Success 201 {int} models.WanbuDataUploadRecord
// @Failure 403 body is empty
// @router / [post]
func (c *WanbuDataUploadRecordController) Post() {
	var v models.WanbuDataUploadRecord
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddWanbuDataUploadRecord(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get WanbuDataUploadRecord by id
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Success 200 {object} models.WanbuDataUploadRecord
// @Failure 403
// @router /count [get]
func (c *WanbuDataUploadRecordController) GetOne() {

	var userid int
	var min, max int64

	var query = make(map[string]string)
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	for k, v := range query {

		if k == "dateline" {

			//如果查找到dateline这个特殊query，查某个时间段的时间，必须符合格式，否则返回格式错误
			ts := strings.Split(v, "-")
			if len(ts) != 2 {
				c.Data["json"] = errors.New("query dateline 格式错误")
			}
			//从字符串转为时间戳，第一个参数是格式，第二个是要转换的时间字符串
			t1, err := time.ParseInLocation("20060102", ts[0], time.Local)
			if err != nil {
				c.Data["json"] = err.Error()
				c.ServeJSON()
				break
			}
			//从字符串转为时间戳，第一个参数是格式，第二个是要转换的时间字符串
			t2, err := time.ParseInLocation("20060102", ts[1], time.Local)
			if err != nil {
				c.Data["json"] = err.Error()
				c.ServeJSON()
				break
			}

			min = t1.Unix()
			max = t2.Unix()
		}
		if k == "touserid" {
			userid, _ = strconv.Atoi(v)
		}
	}

	v, err := models.GetWanbuDataUploadRecordById(userid, min, max)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get WanbuDataUploadRecord
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.WanbuDataUploadRecord
// @Failure 403
// @router / [get]
func (c *WanbuDataUploadRecordController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllWanbuDataUploadRecord(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the WanbuDataUploadRecord
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.WanbuDataUploadRecord	true		"body for WanbuDataUploadRecord content"
// @Success 200 {object} models.WanbuDataUploadRecord
// @Failure 403 :id is not int
// @router /:id [put]
func (c *WanbuDataUploadRecordController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.WanbuDataUploadRecord{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateWanbuDataUploadRecordById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the WanbuDataUploadRecord
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *WanbuDataUploadRecordController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteWanbuDataUploadRecord(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
