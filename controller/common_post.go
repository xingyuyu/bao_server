package controller

import (
	"bao_server/db"
	"fmt"
	"log"
	"net/http"
	"time"
	"strconv"
)

type CommonInfo struct {
	liaotianbaoID string
	selfAttr      string
	expectAttr    string
	weixinId      string
	huodongType   int
}

var common_info CommonInfo
var common_result []byte

//处理通用提交信息
func CommonHandlePost(r *http.Request) []byte {
	err := db.ConnDb()
	if err != nil {
		common_result = errHandle("server error", -1)
		return common_result
	}
	r.ParseForm()
	log.Println("start handle post")
	log.Println("Postform=", r.PostForm)
	common_info.liaotianbaoID = r.PostForm.Get("liaotianbao_id")
	if common_info.liaotianbaoID == "" {
		common_result = errHandle("聊天包id为空", 2008)
		return common_result
	}

	common_info.selfAttr = r.PostForm.Get("self_attr")
	if common_info.selfAttr == "" {
		common_result = errHandle("本人卡片属性为空", 2008)
		return common_result
	}

	common_info.expectAttr = r.PostForm.Get("expect_attr")
	if common_info.expectAttr == "" {
		common_result = errHandle("期望卡片属性为空", 2008)
		return common_result
	}

	common_info.huodongType, err = strconv.Atoi(r.PostForm.Get("huodong_type"))
	common_info.weixinId = r.PostForm.Get("weixin_id")
	log.Println("weixin_id=", common_info.weixinId)
	log.Println("common_info=", common_info)
	sql := constructInsertSql(&common_info)
	log.Println("sql=", sql)
	affectNum, err := db.Insert(&sql)
	log.Println("affect num=", affectNum)
	if err == nil {
		common_result = getSuccMsg()
	} else {
		log.Println("sql inser fail", err.Error())
		common_result = errHandle("server err", -1)
	}
	return common_result
}

func constructInsertSql(info *CommonInfo) string {
	unix_time := time.Now().UnixNano()
	timestamp := unix_time / 1000000000
	sql := fmt.Sprintf("insert into common_exchange(liaotianbao_id, weixin_id, self_attr, expect_attr, update_time, status, huodong_type) values('%s', '%s','%s', '%s', %d, 0, %d);",
		info.liaotianbaoID, 
		info.weixinId, 
		info.selfAttr, 
		info.expectAttr, 
		timestamp,
		info.huodongType)
	return sql
}

