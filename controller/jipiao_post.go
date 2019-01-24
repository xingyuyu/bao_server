package controller

import (
	"bao_server/db"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type JipiaoInfo struct {
	liaotianbaoID string
	selfCity      string
	selfArrive    string
	selfTime      uint64
	expectCity    string
	expectArrive  string
	expectTime    uint64
	weixinId      string
}

var jipiao_info JipiaoInfo
var result []byte

//处理机票提交信息
func HandlePost(r *http.Request) []byte {
	err := db.ConnDb()
	if err != nil {
		result = errHandle("server error", -1)
		return result
	}
	r.ParseForm()
	log.Println("start handle post")
	log.Println("Postform=", r.PostForm)
	jipiao_info.liaotianbaoID = r.PostForm.Get("liaotianbao_id")
	if jipiao_info.liaotianbaoID == "" {
		result = errHandle("聊天包id为空", 2008)
		return result
	}
	city, err := cityHandle(&(r.PostForm), "self_city")
	if err != nil {
		result = errHandle("聊天包id为空", 2008)
		return result
	} else {
		jipiao_info.selfCity = *city
	}

	var cityArrive, erArrive = cityHandle(&(r.PostForm), "self_arrive")
	if erArrive != nil {
		result = errHandle("目的地填写有误", 2008)
		return result
	} else {
		jipiao_info.selfArrive = *cityArrive
	}

	var cityExpect, errExpect = cityHandle(&(r.PostForm), "expect_city")
	if errExpect != nil {
		jipiao_info.expectCity = ""
	} else {
		jipiao_info.expectCity = *cityExpect
	}

	var cityEA, errEA = cityHandle(&(r.PostForm), "expect_arrive")
	if errEA != nil {
		jipiao_info.expectArrive = ""
	} else {
		jipiao_info.expectArrive = *cityEA
	}
	var selfTime, errST = timeHandle(&(r.PostForm), "self_time")
	if errST != nil {
		result = errHandle("出发时间填写错误", 2008)
		return result
	} else {
		jipiao_info.selfTime = selfTime
	}

	var expectTime, errET = timeHandle(&(r.PostForm), "expect_time")
	if errET != nil {
		jipiao_info.expectTime = 0
	} else {
		jipiao_info.expectTime = expectTime
	}
	jipiao_info.weixinId = r.PostForm.Get("weixin_id")
	log.Println("weixin_id=", jipiao_info.weixinId)
	log.Println("jipiao_info=", jipiao_info)
	sql := constructInserSql(&jipiao_info)
	log.Println("sql=", sql)
	affectNum, err := db.Insert(&sql)
	log.Println("affect num=", affectNum)
	if err == nil {
		result = getSuccMsg()
	} else {
		log.Println("sql inser fail", err.Error())
		result = errHandle("server err", -1)
	}
	return result
}

func cityHandle(postForm *url.Values, filedName string) (*string, error) {
	var cityName = strings.Trim(postForm.Get(filedName), " ")
	if cityName == "" {
		return &cityName, errors.New("empty")
	}
	cityName = strings.TrimSuffix(cityName, "市")
	return &cityName, nil
}

func timeHandle(postForm *url.Values, filedName string) (uint64, error) {
	var time = postForm.Get(filedName)
	parseTime, err := strconv.ParseUint(time, 10, 64)
	if err != nil {
		return 0, errors.New("parse error")
	}
	return parseTime, nil
}

func constructInserSql(info *JipiaoInfo) string {
	timestamp := time.Now().Unix()
	sql := fmt.Sprintf("insert into jipiao_exchange(liaotianbao_id, weixin_id,self_city, self_arrive,self_time,expect_city,expect_arrive,expect_time,update_time,status) values('%s', '%s','%s', '%s', %d, '%s', '%s', %d, %d, 0) ON DUPLICATE KEY UPDATE liaotianbao_id='%s', weixin_id='%s', self_city='%s',self_arrive='%s',self_time=%d,expect_city='%s',expect_arrive='%s',expect_time=%d,update_time=%d;",
		info.liaotianbaoID,
		info.weixinId,
		info.selfCity,
		info.selfArrive,
		info.selfTime,
		info.expectCity,
		info.expectArrive,
		info.expectTime,
		timestamp,
		info.liaotianbaoID,
		info.weixinId,
		info.selfCity,
		info.selfArrive,
		info.selfTime,
		info.expectCity,
		info.expectArrive,
		info.expectTime,
		timestamp)
	return sql
}

type PostRetMsg struct {
	Msg   string `json:"msg"`
	Errno int    `json:"errno"`
}

func errHandle(msg string, errNo int) []byte {
	err_msg := &PostRetMsg{Msg: msg, Errno: errNo}
	log.Println("err msg=", err_msg)
	b, err := json.Marshal(err_msg)
	if err != nil {
		log.Println("err json encode error")
	}
	log.Println("err json=", string(b))
	return b
}

func getSuccMsg() []byte {
	msg := &PostRetMsg{Msg: "success", Errno: 0}
	b, err := json.Marshal(msg)
	if err != nil {
		log.Println("success msg json encode error")
	}
	log.Println("json=", string(b))
	return b
}
