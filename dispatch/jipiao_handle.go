package dispatch

import (
	"bao_server/controller"
	"fmt"
	"log"
	"net/http"
)

type JipiaoServer struct{}

var handle JipiaoServer

func NewJipiaoHandle() *JipiaoServer {
	return &handle
}

func (*JipiaoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("handle,uri=", r.URL.Path)
	if r.URL.Path == "/jipiao/post_info" {
		ret := controller.HandlePost(r)
		fmt.Fprint(w, string(ret))
	} else if r.URL.Path == "jipiao/search" {

	} else if r.URL.Path == "jipiao/update_status" {

	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
