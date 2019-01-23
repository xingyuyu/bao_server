package controller

import (
	"bao_server/db"
	"log"
	"testing"
	"time"
)

func Test_InitDB(t *testing.T) {
	if db.InitDB() { //try a unit test on function
		t.Log("connect to db success") // 如果不是如预期的那么就报错
	} else {
		t.Error("connect to db success") //记录一些你期望记录的信息
	}
}

func Test_SearchByExpectTime(t *testing.T) {
	//getExpectTimeData(122355666)
	//tm := time.Now().Unix()
	var tm1 int64
	tm1 = 1548253623
	tmObj := time.Unix(tm1, 0)
	format := tmObj.Format("01-01")
	log.Println(format)
}

func Test_formatWeixinResponse(t *testing.T) {
	var test WeiXinResponseMsg
	var content Cdata
	content.Value = "聊天包ID:xxx 出发城市: 北京\n聊天包ID:cccc 出发城市: 北京"
	test.Content = content
	test.CreateTime = 1323234
	test.FromUserName = Cdata{"dddd"}
	test.MsgType = Cdata{"text"}
	test.ToUserName = Cdata{"sssss"}
	xml := formatWeixinResponse(&test)
	log.Println(xml)
}

func Test_ParseReq(t *testing.T) {
	xml := "<xml><ToUserName><![CDATA[toUser]]></ToUserName>  <FromUserName><![CDATA[fromUser]]></FromUserName> <CreateTime>1348831860</CreateTime><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[this is a test]]></Content><MsgId>1234567890123456</MsgId></xml>"
	result := parseReqParam(xml)
	log.Println(*result)
}

func Test_Time(t *testing.T) {
	res, err := parseTime("02-3")
	if err != nil {
		log.Println("xxx")

	}
	log.Println("time", res)
}

func Test_parseUserSemantics(t *testing.T) {
	result, err := parseUserSemantics("南京")
	if err == nil {
		log.Println("resul", result)
	}
	result1, err1 := parseUserSemantics("南京 北京")
	if err1 == nil {
		log.Println("resul", result1)
	}
	result2, err2 := parseUserSemantics("02-03")
	if err2 == nil {
		log.Println("resul", result2)
	}
	result3, err3 := parseUserSemantics("南京 北京 04-05")
	if err3 == nil {
		log.Println("resul", result3)
	}

	result4, err4 := parseUserSemantics("南京 北京 04-99")
	if err4 == nil {
		log.Println("resul", result4)
	} else {
		log.Println("error=", err4.Error())
	}
}
