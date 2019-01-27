package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const urlFormat = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"

type UnionID struct {
	ID string `xml:unionid`
}

func GetUnionID(open_id string, access_token string) string {
	url := fmt.Sprintf(urlFormat, access_token, open_id)
	res, err := http.Get(url)
	var body []byte
	if err == nil {
		body, _ = ioutil.ReadAll(res.Body)
		var unionID UnionID
		err := json.Unmarshal(body, unionID)
		if err == nil {
			return unionID.ID
		} else {
			log.Println("parse union id json error=", string(body))
		}
	} else {
		log.Fatalln("get union id fail")
	}
	return ""
}
