package controller

import (
	"bao_server/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var common_search_result []byte

type CommonRet struct {
	LiaotianbaoId string `json:"liaotianbao_id"`
	SelfAttr string `json:"self_attr`
	ExpectAttr string `json:"expect_attr"`
	WeixinId string `json:"weixin_id"`
}

type CommonRetMsg struct {
	Msg   string `json:"msg"`
	Errno int    `json:"errno"`
	CommonList []CommonRet `json:"common_list"`
}

//处理机票提交信息
func CommonHandleSearch(r *http.Request) []byte {
	err := db.ConnDb()
	if err != nil {
		common_search_result = errHandle("server error", -1)
		return common_search_result
	}
	vars := r.URL.Query();
	weixin_id, weixin_id_ok := vars["weixin_id"]
	huodong_type, type_ok := vars["huodong_type"]
	log.Println("start handle get common list")
	log.Println("URL=", r.URL)
		
	var sql string;
	if !type_ok {
		common_search_result = errHandle("活动类型为空", 2008)
		return common_search_result
	}
	huodong_type_id, err := strconv.Atoi(huodong_type[0])

	if !weixin_id_ok {
		var tmp string;
		tmp = ""
		sql = constructCommonQuerySql(&tmp, huodong_type_id)
	} else {
    		fmt.Printf("param weixin_id a value is [%s] huodong_id a value is [%d]\n", weixin_id, huodong_type_id);
  		sql = constructCommonQuerySql(&weixin_id[0], huodong_type_id)
	}

	log.Println("sql=", sql)
        rows, err := db.Select(&sql)

	if err == nil {
		var common_ret_msg CommonRetMsg
		common_ret_msg.Msg = "success"
		common_ret_msg.Errno = 0
		for rows.Next() {
			var common_ret CommonRet
			err = rows.Scan(&common_ret.LiaotianbaoId, &common_ret.WeixinId, &common_ret.SelfAttr, &common_ret.ExpectAttr)
			log.Println("sql=", common_ret)
			common_ret_msg.CommonList = append(common_ret_msg.CommonList, common_ret)
		}
		log.Println("common_result=", common_ret_msg)
		b, err := json.Marshal(common_ret_msg)
		if err != nil {
			log.Println("success json encode error")
		}
		log.Println("json=", string(b))
		common_search_result = b
	} else {
		log.Println("sql select fail", err.Error())
		common_search_result = errHandle("server err", -1)
	}
	return common_search_result
}


func constructCommonQuerySql(weixin *string, huodong_type int) string {
	var sql string;
	if *weixin == "" {
		sql = fmt.Sprintf("select liaotianbao_id, weixin_id, self_attr, expect_attr from common_exchange where huodong_type = %d order by update_time limit 20", huodong_type)
	} else {
		sql = fmt.Sprintf("select liaotianbao_id, weixin_id, self_attr, expect_attr from common_exchange where weixin_id = '%s' and huodong_type = %d;", *weixin, huodong_type)
	}
	return sql
}
