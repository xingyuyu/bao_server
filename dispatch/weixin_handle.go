package dispatch

import (
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
		} else if r.Method == "POST" {
			r.ParseMultipartForm(1024)
			body, _ := ioutil.ReadAll(r.Body)
			log.Println("weixin post msg=", string(body))
		}
	}
}
