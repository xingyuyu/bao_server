package dispatch

import (
	"bao_server/controller"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type WeixinServer struct{}

var wxHandle WeixinServer

func NewWeixinHandle() *WeixinServer {
	return &wxHandle
}

func (*WeixinServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("weixin handle,uri=", r.URL.Path)
	// 处理微信验证消息
	if strings.Trim(r.URL.Path, " ") == "/" {
		r.ParseForm()
		if r.Method == "GET" {
			fmt.Fprint(w, r.Form.Get("echostr"))
		} else if r.Method == "POST" { // 微信消息处理
			r.ParseMultipartForm(1024)
			body, _ := ioutil.ReadAll(r.Body)
			bodyStr := string(body)
			xmlRet := controller.HandleSearch(&bodyStr)
			log.Println("weixin post msg=", string(body))
			sendResult, err := fmt.Fprint(w, xmlRet)
			if err != nil {
				log.Println("send to weixin fail")
			} else {
				log.Println("end request,sendResult=", sendResult)
			}
		}
	}
}
