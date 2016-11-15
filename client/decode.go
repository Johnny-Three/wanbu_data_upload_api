package client

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

func SliceAtoi(strArr []string) ([]int, error) {
	// NOTE:  Read Arr as Slice as you like
	var str string // O
	var i int      // O
	var err error  // O

	iArr := make([]int, 0, len(strArr))
	for _, str = range strArr {
		i, err = strconv.Atoi(str)
		if err != nil {
			return nil, err // O
		}
		iArr = append(iArr, i)
	}
	return iArr, nil
}

func Decode(msg string) error {

	js, err := simplejson.NewJson([]byte(msg))
	if err != nil {
		errback := fmt.Sprintf("decode json error the error msg is %s", err.Error())
		return errors.New(errback)
	}

	var wd WalkDayData
	walkdays := []WalkDayData{}
	userwalkdata := User_walkdays_struct{}

	userid := js.Get("userid").MustInt()
	wd.Timestamp = js.Get("timestamp").MustInt64()
	devtype := js.Get("devinfo").Get("reqtype").MustInt()
	arr, _ := js.Get("walkdays").Array()

	for index := range arr {

		walkdate := js.Get("walkdays").GetIndex(index).Get("walkdate").MustInt64()
		wd.WalkDate = walkdate

		var err0 error
		walkhour := js.Get("walkdays").GetIndex(index).Get("walkhour").MustString()
		wd.Hourdata, err0 = SliceAtoi(strings.Split(walkhour, ","))
		if err0 == nil {

			if len(wd.Hourdata) != 24 {
				errback := fmt.Sprintf("uid %d walkdate %d get wrong hourdata %v format", userid, walkdate, wd.Hourdata)
				return errors.New(errback)
			}
		}

		wd.Daydata = js.Get("walkdays").GetIndex(index).Get("walktotal").MustInt()
		wd.Faststepnum = js.Get("walkdays").GetIndex(index).Get("fast").MustInt()
		wd.Effecitvestepnum = js.Get("walkdays").GetIndex(index).Get("effective").MustInt()
		srecipe := js.Get("walkdays").GetIndex(index).Get("recipe").MustString()
		irecipe, err1 := SliceAtoi(strings.Split(srecipe, ","))
		if err1 == nil {

			if len(irecipe) != 3 {
				errback := fmt.Sprintf("uid %d walkdate %d get wrong recipe %v format", userid, walkdate, irecipe)
				return errors.New(errback)
			}
		}
		//no problem .. then assign the chufang related value..
		wd.Chufangid = irecipe[0]
		wd.Chufangfinish = irecipe[1]
		wd.Chufangtotal = irecipe[2]

		//用户此次上传的数据消息存储在MAP中..
		walkdays = append(walkdays, wd)

	}

	userwalkdata.Uid = userid
	userwalkdata.Timestamp = wd.Timestamp
	userwalkdata.Devtype = devtype
	userwalkdata.Walkdays = walkdays

	User_walk_data_chan <- userwalkdata

	return nil
}

func init() {

	User_walk_data_chan = make(chan User_walkdays_struct, 16)
}
