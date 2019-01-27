package token

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type WeixinAccessTokenRes struct {
	Token   string `json:"access_token"`
	Expires int64  `json:"expires_in"`
}

type AccessToken struct {
	Token   string
	Expires int64
	ReqTime int64
}

var token AccessToken
var WRLock sync.RWMutex

func GetToken() string {
	WRLock.RLock()
	defer WRLock.RUnlock()
	return token.Token
}

func getWeixinAccessToken() {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=xxx&secret=xxx"
	res, err := http.Get(url)
	if err == nil {
		if res.StatusCode == http.StatusOK {
			body, _ := ioutil.ReadAll(res.Body)
			log.Println("access token body=", string(body))
			var parseResult WeixinAccessTokenRes
			err := json.Unmarshal(body, &parseResult)
			if err == nil {
				log.Println("parse access token info=", parseResult)
				//update data
				WRLock.Lock()
				defer WRLock.Unlock()
				token.Token = parseResult.Token
				log.Println("new access token=", token.Token)
				token.Expires = parseResult.Expires
				token.ReqTime = time.Now().Unix()
			} else {
				log.Fatalln("req weixin access token res json parse error")
			}

		} else {
			log.Fatalln("req weixin access token res status not ok")
		}
	} else {
		log.Fatalln("req weixin access token error")
	}
}

func parseAccessToken(jsonData []byte) (*WeixinAccessTokenRes, error) {
	var parseResult WeixinAccessTokenRes
	err := json.Unmarshal(jsonData, &parseResult)
	if err != nil {
		return nil, errors.New("json parse access token error")
	}
	log.Println("parse access token =", parseResult)
	return &parseResult, nil
}

func startGetAccessThread() {
	for {
		now := time.Now().Unix()
		//提前5s获取token
		if now >= token.ReqTime+token.Expires {
			getWeixinAccessToken()
		} else {
			log.Println("wait")
		}
		time.Sleep(1)
	}
}

func InitGetAccessToken() {
	go startGetAccessThread()
}
