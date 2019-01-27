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
	Event        Cdata    `xml:"Event"`
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
	action      string
	info        JipiaoInfo
	common_info CommonInfo
}

func HandleSearch(reqBody *string) []byte {
	weixinReceive := parseReqParam(*reqBody)
	//回复关注信息
	if weixinReceive.MsgType.Value == "event" && weixinReceive.Event.Value == "subscribe" {
		msg := getSubEvent()
		resMsg := contructWeiXinResponse(weixinReceive.FromUserName.Value, weixinReceive.ToUserName.Value, *msg)
		xmlResult := formatWeixinResponse(resMsg)
		return []byte(*xmlResult)
	}
	parseResult, err := parseUserSemantics(weixinReceive.Content.Value, weixinReceive.FromUserName.Value)
	var msg string
	var searchInfo []JipiaoInfo
	var commonInfo []CommonInfo
	var meWurenjiInfo []CommonInfo
	var meSongshuInfo []CommonInfo
	var meDianziyanInfo []CommonInfo

	if err == nil {
		switch parseResult.action {
		case "do_me_search":
			// log.Println("do_me_search,weixinId=", parseResult.info.weixinId)
			// myInfo := getMeInfoData(parseResult.info.weixinId)
			// myCommonInfo := getMeCommonInfoData(parseResult.info.weixinId)
			// if myInfo == nil {
			// 	msg = "没查询到您相关的活动信息"
			// } else {
			// 	searchInfo = getAllInfoData(myInfo.selfCity, myInfo.selfArrive, myInfo.selfTime)
			// 	meWurenjiInfo = getWurenjiInfo(myCommonInfo[0].selfAttr, myCommonInfo[0].expectAttr)
			// 	meSongshuInfo = getSongshuInfo(myCommonInfo[1].selfAttr, myCommonInfo[0].expectAttr)
			// 	meDianziyanInfo = getDianziyanInfo(myCommonInfo[1].selfAttr, myCommonInfo[0].expectAttr)
			// }
			break
		case "do_expect_city_search":
			log.Println("getExpectCityData,info=", parseResult)
			searchInfo = getExpectCityData(parseResult.info.selfCity)
			break
		case "do_expect_time":
			log.Println("getExpectTimeData info=", parseResult)
			searchInfo = getExpectTimeData(parseResult.info.selfTime)
			break
		case "do_city_search":
			log.Println("do_city_search info=", parseResult)
			searchInfo = getExpectCityAndArriveData(parseResult.info.selfCity, parseResult.info.selfArrive)
			break
		case "do_all_search":
			log.Println("do_all_search info=", parseResult)
			searchInfo = getAllInfoData(parseResult.info.selfCity, parseResult.info.selfArrive, parseResult.info.selfTime)
			break
		case "do_wurenji_search":
			log.Println("do_wurenji_search=", parseResult)
			commonInfo = getWurenjiInfo(parseResult.common_info.selfAttr, parseResult.common_info.expectAttr)
			break
		case "do_dianziyan_search":
			log.Println("do_dianziyan_search=", parseResult)
			commonInfo = getDianziyanInfo(parseResult.common_info.selfAttr, parseResult.common_info.expectAttr)
			break
		case "do_songshu_search":
			log.Println("do_songshu_search=", parseResult)
			commonInfo = getSongshuInfo(parseResult.common_info.selfAttr, parseResult.common_info.expectAttr)
			break
		}
	} else {
		msg = "请正确提交信息"
	}
	if parseResult.action == "do_wurenji_search" {
		if len(commonInfo) == 0 {
			msg = "sorry，没有查询到与您相关的无人机信息"
		} else {
			msg = formatCommonData(commonInfo)
		}
	} else if parseResult.action == "do_dianziyan_search" {
		if len(commonInfo) == 0 {
			msg = "sorry，没有查询到与您相关的电子烟信息"
		} else {
			msg = formatCommonData(commonInfo)
		}
	} else if parseResult.action == "do_songshu_search" {
		if len(commonInfo) == 0 {
			msg = "sorry，没有查询到与您相关的三只松鼠信息"
		} else {
			msg = formatCommonData(commonInfo)
		}
	} else if parseResult.action == "do_me_search" {
		jipiaoMsg := fomatData(searchInfo)
		wurenjiMsg := formatCommonData(meWurenjiInfo)
		dianziyanMsg := formatCommonData(meDianziyanInfo)
		songshuMsg := formatCommonData(meSongshuInfo)
		msgFormat := "为您查询到得机票信息如下:\n%s\n为您查询到的无人机信息如下:\n%s\n为你查询到的电子烟信息如下:\n%s\n为您查询到的三只松鼠信息如下:\n%s\n"
		msg = fmt.Sprintf(msgFormat, jipiaoMsg, wurenjiMsg, dianziyanMsg, songshuMsg)
	} else if parseResult.action == "do_err" {
		errFormat := "您输入了我无法识别的指令哟，请重新输入。我们的规则如下\n%s"
		msg = fmt.Sprintln(errFormat, getSubEvent())
	} else {
		if len(searchInfo) == 0 {
			msg = "sorry,没查到相关航班信息"
		} else {
			msg = fomatData(searchInfo)
		}
	}
	resMsg := contructWeiXinResponse(weixinReceive.FromUserName.Value, weixinReceive.ToUserName.Value, msg)
	xmlResult := formatWeixinResponse(resMsg)
	return []byte(*xmlResult)
}

func handleDbRowData(rows *sql.Rows) []JipiaoInfo {
	infos := make([]JipiaoInfo, 0)
	if rows != nil {
		for rows.Next() {
			var item JipiaoInfo
			rows.Scan(&item.liaotianbaoID, &item.selfCity, &item.selfArrive, &item.selfTime)
			infos = append(infos, item)
		}
	}
	return infos
}

func handleDbCommonRowData(rows *sql.Rows) []CommonInfo {
	infos := make([]CommonInfo, 0)
	if rows != nil {
		for rows.Next() {
			var item CommonInfo
			rows.Scan(&item.liaotianbaoID, &item.weixinId, &item.selfAttr, &item.expectAttr)
			infos = append(infos, item)
		}
	}
	return infos
}

func getExpectTimeData(expectTime uint64) []JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where self_time = %d order by update_time limit 20"
	sql := fmt.Sprintf(sqlFormat, expectTime)
	log.Println("search by expect time sql=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbRowData(rows)
	}
	return nil
}

func getExpectCityData(city string) []JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where self_city = '%s' and self_time > 0 order by update_time limit 20 "
	sql := fmt.Sprintf(sqlFormat, city)
	log.Println("search by expect city sql=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbRowData(rows)
	}
	return nil
}

func getExpectArriveData(city string) []JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where self_city = '%s' order by update_time"
	sql := fmt.Sprintf(sqlFormat, city)
	log.Println("search by expect city sql=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbRowData(rows)
	}
	return nil
}

func getExpectCityAndArriveData(city string, arraive string) []JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where self_city = '%s' and self_arrive = '%s' order by update_time limit 20"
	sql := fmt.Sprintf(sqlFormat, city, arraive)
	log.Println("search by getExpectCityAndArriveData=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbRowData(rows)
	}
	return nil
}

func getAllInfoData(city string, arraive string, time uint64) []JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where self_city = '%s' and self_arrive = '%s' and self_time=%d order by update_time limit 20"
	sql := fmt.Sprintf(sqlFormat, city, arraive, time)
	log.Println("search by getAllInfoData=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbRowData(rows)
	}
	return nil
}

func getMeInfoData(weixinID string) *JipiaoInfo {
	sqlFormat := "select liaotianbao_id,self_city,self_arrive,self_time from jipiao_exchange where weixin_id = '%s'"
	sql := fmt.Sprintf(sqlFormat, weixinID)
	log.Println("search by getMeInfoData=", sql)
	rows, err := db.Select(&sql)
	if !rows.Next() || err != nil {
		return nil
	}
	var item JipiaoInfo
	rows.Scan(&item.liaotianbaoID, &item.selfCity, &item.selfArrive, &item.selfTime)
	log.Println("result=", result)
	return &item
}

func getMeCommonInfoData(weixinID string) map[int]CommonInfo {
	sqlFormat := "select liaotianbao_id,weixin_id,self_attr,expect_attr,huodong_type from common_exchange where weixin_id = '%s'"
	sql := fmt.Sprintf(sqlFormat, weixinID)
	log.Println("search by getMeCommonData=", sql)
	rows, err := db.Select(&sql)
	if !rows.Next() || err != nil {
		return nil
	}
	var resultMap map[int]CommonInfo
	for rows.Next() {
		var item CommonInfo
		rows.Scan(&item.liaotianbaoID, &item.weixinId, &item.selfAttr, &item.expectAttr, &item.huodongType)
		resultMap[item.huodongType] = item
	}
	return resultMap
}

func getWurenjiInfo(self string, expect string) []CommonInfo {
	sqlFormat := "select liaotianbao_id,weixin_id,self_attr,expect_attr from common_exchange where huodong_type = 0 and self_attr = '%s' and expect_attr = '%s' limit 20;"
	sql := fmt.Sprintf(sqlFormat, expect, self)
	log.Println("search by getMeInfoData=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbCommonRowData(rows)
	}
	return nil
}

func getSongshuInfo(self string, expect string) []CommonInfo {
	sqlFormat := "select liaotianbao_id,weixin_id,self_attr,expect_attr from common_exchange where huodong_type = 1 and self_attr = '%s' and expect_attr = '%s' limit 20;"
	sql := fmt.Sprintf(sqlFormat, expect, self)
	log.Println("search by getMeInfoData=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbCommonRowData(rows)
	}
	return nil
}

func getDianziyanInfo(self string, expect string) []CommonInfo {
	sqlFormat := "select liaotianbao_id,weixin_id,self_attr,expect_attr from common_exchange where huodong_type = 2 and self_attr = '%s' and expect_attr = '%s' limit 20;"
	sql := fmt.Sprintf(sqlFormat, expect, self)
	log.Println("search by getMeInfoData=", sql)
	rows, err := db.Select(&sql)
	if err == nil {
		return handleDbCommonRowData(rows)
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
	result = strings.Trim(result, "\n")
	return result
}

func formatCommonData(infos []CommonInfo) string {
	format := "聊天宝账户:%s"
	resulSlice := make([]string, len(infos))
	for _, item := range infos {
		tempStr := fmt.Sprintf(format, item.liaotianbaoID)
		resulSlice = append(resulSlice, tempStr)
	}
	result := strings.Join(resulSlice, "\n")
	result = strings.Trim(result, "\n")
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

func parseUserSemantics(semantics string, weixinID string) (*SemanticsResult, error) {
	splitArr := strings.Split(semantics, " ")
	var result SemanticsResult
	if len(semantics) == 0 {
		return &result, errors.New("查不到您输入的信息")
	}
	if len(splitArr) > 3 {
		return &result, errors.New("查不到您输入的信息")
	} else {
		if len(splitArr) == 1 && (strings.Contains(semantics, "-")) {
			timeStamp, err := ParseTime(semantics)
			if err != nil {
				return &result, errors.New("您输入的查询时间不对")
			}
			result.action = "do_expect_time"
			result.info.selfTime = uint64(timeStamp)
		} else if len(splitArr) == 1 {
			result.action = "do_expect_city_search"
			result.info.selfCity = splitArr[0]
		} else if len(splitArr) == 3 && splitArr[0] == "无人机" {
			result.action = "do_wurenji_search"
			result.common_info.selfAttr = splitArr[1]
			result.common_info.expectAttr = splitArr[2]
			result.common_info.huodongType = 0
		} else if len(splitArr) == 3 && splitArr[0] == "电子烟" {
			result.action = "do_dianziyan_search"
			result.common_info.selfAttr = splitArr[1]
			result.common_info.expectAttr = splitArr[2]
			result.common_info.huodongType = 2
		} else if len(splitArr) == 3 && splitArr[0] == "三只松鼠" {
			result.action = "do_songshu_search"
			result.common_info.selfAttr = splitArr[1]
			result.common_info.expectAttr = splitArr[2]
			result.common_info.huodongType = 1
		} else if len(splitArr) == 3 {
			result.action = "do_all_search"
			result.info.selfCity = splitArr[0]
			result.info.selfArrive = splitArr[1]
			timeStamp, err := ParseTime(splitArr[2])
			if err != nil {
				return &result, errors.New("您输入的查询时间不对")
			}
			result.info.selfTime = uint64(timeStamp)
		} else if len(splitArr) == 2 {
			result.action = "do_city_search"
			result.info.selfCity = splitArr[0]
			result.info.selfArrive = splitArr[1]
		} else {
			result.action = "do_err"
		}
		return &result, nil
	}
}

func formatTime(timestamp uint64) string {
	if timestamp == 0 {
		return "暂无填写出发时间"
	}
	tmObj := time.Unix(int64(timestamp), 0)
	result := tmObj.Format("2006-01-02 15:04:05")
	tempArr := strings.Split(result, " ")
	return tempArr[0]
}

func formatWeixinResponse(weixinRes *WeiXinResponseMsg) *string {
	output, err := xml.MarshalIndent(*weixinRes, "", "")
	retStr := ""
	if err != nil {
		log.Println("encode xml fail")
		return &retStr
	}
	retStr = string(output)
	return &retStr
}

func contructWeiXinResponse(to string, from string, msg string) *WeiXinResponseMsg {
	var weixinResMsg WeiXinResponseMsg
	toCData := &weixinResMsg.ToUserName
	toCData.Value = to
	fromCData := &weixinResMsg.FromUserName
	fromCData.Value = from
	msgCData := &weixinResMsg.Content
	msgCData.Value = msg
	msgTypeCData := &weixinResMsg.MsgType
	msgTypeCData.Value = "text"
	weixinResMsg.CreateTime = uint64(time.Now().Unix())
	return &weixinResMsg
}

func ParseTime(inputTime string) (int64, error) {
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

func getSubEvent() *string {
	str := "亲，您终于来了！新春“送换夺”就差您一个了！\n点击“参与活动”，尽快填写“送换夺卡券”信息，就能最快找到您的聊天宝小伙伴！\n\n登记完毕，在输入框中输入以下规则，就可以查询到匹配合适的聊天宝ID，完成新春大挑战。\n机票活动 输入“出发地 到达地 日期”，如“北京 上海 02-01”，可以进行查询。\n无人机活动 输入“无人机 您持有的省份 想交换的省份”，如“无人机 浙江 北京”，可以进行查询。\n三只松鼠活动 输入“坚果 您的尾号 想交换的尾号”，如“坚果 1 9”，可以进行查询。\n电子烟活动 输入“电子烟 您的数字 想交换的数字”，如“电子烟 10 6”，可以进行查询。"
	return &str
}
