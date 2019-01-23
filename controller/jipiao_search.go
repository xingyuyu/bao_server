package controller

import (
	"bao_server/db"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

//微信发送过来消息体
type WeiXinReceiveMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   Cdata    `xml:"ToUserName"`
	FromUserName Cdata    `xml:"FromUserName"`
	CreateTime   Cdata    `xml:"CreateTime"`
	MsgType      Cdata    `xml:"MsgType"`
	Content      Cdata    `xml:"Content"`
	MsgId        Cdata    `xml:"MsgId"`
}

type Cdata struct {
	Value string `xml:",cdata"`
}

//发送给微信消息体
type WeiXinResponseMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   Cdata    `xml:"ToUserName"`
	FromUserName Cdata    `xml:"FromUserName"`
	MsgType      Cdata    `xml:"MsgType"`
	Content      Cdata    `xml:"Content"`
	CreateTime   uint64
}

//用户发送语义解析结构
type SemanticsResult struct {
	action string
	info   JipiaoInfo
}

// func HandleSearch(r *http.Request) []byte {

// }

func handleDbRowData(rows *sql.Rows) []JipiaoInfo {
	infos := make([]JipiaoInfo, 0)
	for rows.Next() {
		var item JipiaoInfo
		rows.Scan(&item.liaotianbaoID, &item.selfCity, &item.selfArrive, &item.selfTime)
		infos = append(infos, item)
	}
	return infos
}

func getExpectTimeData(expectTime uint64) []JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where self_time = %d"
	sql := fmt.Sprintf(sqlFormat, expectTime)
	log.Println("search by expect time sql=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbRowData(rows)
	}
	return nil
}

func getExpectCityData(city string) []JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where self_city = %s"
	sql := fmt.Sprintf(sqlFormat, city)
	log.Println("search by expect city sql=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbRowData(rows)
	}
	return nil
}

func getExpectArriveData(city string) []JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where self_city = %s"
	sql := fmt.Sprintf(sqlFormat, city)
	log.Println("search by expect city sql=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbRowData(rows)
	}
	return nil
}

func getExpectCityAndArriveData(city string, arraive string) []JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where self_city = %s and self_arrive = %s"
	sql := fmt.Sprintf(sqlFormat, city, arraive)
	log.Println("search by getExpectCityAndArriveData=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbRowData(rows)
	}
	return nil
}

func getAllInfoData(city string, arraive string, time uint64) []JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where self_city = %s and self_arrive = %s and self_time=%d"
	sql := fmt.Sprintf(sqlFormat, city, arraive, time)
	log.Println("search by getAllInfoData=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbRowData(rows)
	}
	return nil
}

//格式化后格式如下
// 聊天宝账户:xxxx
// 出发城市:xx
//到达城市:xx
// 出发时间: xx
// --------------
func fomatData(infos []JipiaoInfo) string {
	format := "聊天宝账户:%s 出发城市:%s 到达城市:%s 出发时间:%s"
	resulSlice := make([]string, len(infos))
	for _, item := range infos {
		tempStr := fmt.Sprintf(format, item.liaotianbaoID, item.selfCity, item.selfArrive, formatTime(item.selfTime))
		resulSlice = append(resulSlice, tempStr)
	}
	result := strings.Join(resulSlice, "\n")
	return result
}

func parseReqParam(content string) *WeiXinReceiveMsg {
	var receiveMsg WeiXinReceiveMsg
	err := xml.Unmarshal([]byte(content), &receiveMsg)
	if err != nil {
		log.Println("weixin req xml parse fail")
	}
	return &receiveMsg
}

func parseUserSemantics(semantics string) (*SemanticsResult, error) {
	splitArr := strings.Split(semantics, " ")
	var result SemanticsResult
	if len(semantics) == 0 {
		return &result, errors.New("查不到您输入的信息")
	}
	if len(splitArr) > 3 {
		return &result, errors.New("查不到您输入的信息")
	} else {
		if len(splitArr) == 1 && (semantics == "我的" || semantics == "我") {
			result.action = "do_me_search"
			result.info.selfCity = splitArr[0]
		} else if len(splitArr) == 1 && (strings.Contains(semantics, "-")) {
			timeStamp, err := parseTime(semantics)
			if err != nil {
				return &result, errors.New("您输入的查询时间不对")
			}
			result.action = "do_expect_city"
			result.info.selfTime = uint64(timeStamp)
		} else if len(splitArr) == 1 {
			result.action = "do_expect_city_search"
			result.info.selfCity = splitArr[0]
		} else if len(splitArr) == 2 {
			result.action = "do_city_search"
			result.info.selfCity = splitArr[0]
			result.info.selfArrive = splitArr[1]
		} else if len(splitArr) == 3 {
			result.action = "do_all_search"
			result.info.selfCity = splitArr[0]
			result.info.selfArrive = splitArr[1]
			timeStamp, err := parseTime(splitArr[2])
			if err != nil {
				return &result, errors.New("您输入的查询时间不对")
			}
			result.info.selfTime = uint64(timeStamp)
		}
		return &result, nil
	}
}

func formatTime(timestamp uint64) string {
	tmObj := time.Unix(int64(timestamp), 0)
	result := tmObj.Format("01-01")
	return result
}

func formatWeixinResponse(weixinRes *WeiXinResponseMsg) string {
	output, err := xml.MarshalIndent(*weixinRes, "", "")
	if err != nil {
		log.Println("encode xml fail")
		return ""
	}
	return string(output)
}

func parseTime(inputTime string) (int64, error) {
	format := "2019-%s 00:00:00"
	timeStr := fmt.Sprintf(format, inputTime)
	loc, _ := time.LoadLocation("Local")
	result, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	if err != nil {
		return 0, errors.New("输入时间解析失败")
	} else {
		return result.Unix(), nil
	}
}
