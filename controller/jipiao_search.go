package controller

//微信发送过来消息体
type WeiXinReceiveMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   uint64
	MsgType      string
	Content      string
	MsgId        string
}

//发送给微信消息体
type WeiXinResponseMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   uint64
	MsgType      string
	Content      string
}

//用户发送语义解析结构
type SemanticsResult struct {
	action string
	info   JipiaoInfo
}

// func HandleSearch(r *http.Request) []byte {

// }

// func getExpectTimeData(expectTime uint16) []JipiaoInfo {

// }

// func getExpectCityData(city string) []JipiaoInfo {

// }

// func getExpectArriveData(city string) []JipiaoInfo {

// }

// func getExpectCityAndArriveData(city string, arraive string) []JipiaoInfo {

// }

// func getAllInfoData(city string, arraive string, time uint32) []JipiaoInfo {

// }

// func fomatData([]JipiaoInfo) []byte {

// }

// func parseReqParam(content string) *WeiXinMsg {

// }

func parseUserSemantics(semantics string) {

}
