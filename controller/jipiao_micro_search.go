package controller

import (
	"bao_server/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var jipiao_search_result []byte

type JipiaoRet struct {
	LiaotianbaoId string `json:"liaotianbao_id"`
	SelfCity string `json:"self_city"`
	SelfArrive string `json:"self_arrive"`
	SelfTime uint64 `json:"self_time"`
	ExpectCity string `json:"expect_city"`
	ExpectArrive string `json:"expect_arrive"`
	ExpectTime uint64 `json:"expect_time"`
	WeixinId string `json:"weixin_id"`
}

type JipiaoRetMsg struct {
	Msg   string `json:"msg"`
	Errno int    `json:"errno"`
	JipiaoList []JipiaoRet `json:"jipiao_list"`
}

//处理机票提交信息
func JipiaoHandleSearch(r *http.Request) []byte {
	err := db.ConnDb()
	if err != nil {
		jipiao_search_result = errHandle("server error", -1)
		return jipiao_search_result
	}
	vars := r.URL.Query();
	weixin_id, ok := vars["weixin_id"]
	log.Println("start handle get jipiao list")
	log.Println("URL=", r.URL)
		
	var sql string;
	if !ok {
		var tmp string;
		tmp = ""
		sql = constructQuerySql(&tmp)
	} else {
    		fmt.Printf("param weixin_id a value is [%s]\n", weixin_id);
  		sql = constructQuerySql(&weixin_id[0])
	}

	log.Println("sql=", sql)
        rows, err := db.Select(&sql)

	if err == nil {
		var jipiao_ret_msg JipiaoRetMsg
		jipiao_ret_msg.Msg = "success"
		jipiao_ret_msg.Errno = 0
		for rows.Next() {
			var jipiao_ret JipiaoRet
			err = rows.Scan(&jipiao_ret.LiaotianbaoId, &jipiao_ret.WeixinId, &jipiao_ret.SelfCity, &jipiao_ret.SelfArrive, &jipiao_ret.SelfTime, &jipiao_ret.ExpectCity, &jipiao_ret.ExpectArrive, &jipiao_ret.ExpectTime)
			log.Println("sql=", err)
			log.Println("sql=", jipiao_ret)
			jipiao_ret_msg.JipiaoList = append(jipiao_ret_msg.JipiaoList, jipiao_ret)
		}
		log.Println("jipiao_result=", jipiao_ret_msg)
		b, err := json.Marshal(jipiao_ret_msg)
		if err != nil {
			log.Println("success json encode error")
		}
		log.Println("json=", string(b))
		jipiao_search_result = b
	} else {
		log.Println("sql select fail", err.Error())
		jipiao_search_result = errHandle("server err", -1)
	}
	return jipiao_search_result
}


func constructQuerySql(weixin *string) string {
	var sql string;
	if *weixin == "" {
		sql = "select liaotianbao_id, weixin_id, self_city, self_arrive, self_time, expect_city, expect_arrive, expect_time from jipiao_exchange order by update_time limit 20"
	} else {
		sql = fmt.Sprintf("select liaotianbao_id, weixin_id, self_city, self_arrive, self_time, expect_city, expect_arrive, expect_time from jipiao_exchange where weixin_id = '%s';", *weixin)
	}
	return sql
}

func getListSuccMsg() []byte {
	msg := &PostRetMsg{Msg: "success", Errno: 0}
	b, err := json.Marshal(msg)
	if err != nil {
		log.Println("success msg json encode error")
	}
	log.Println("json=", string(b))
	return b
}
