package client

var User_walk_data_chan chan User_walkdays_struct

type User_walkdays_struct struct {
	Uid       int
	Timestamp int64
	Devtype   int
	Walkdays  []WalkDayData
}

type WalkDayData struct {
	Daydata          int
	Hourdata         []int
	Chufangid        int
	Chufangfinish    int
	Chufangtotal     int
	Faststepnum      int
	Effecitvestepnum int
	WalkDate         int64
	Timestamp        int64
}
