package dispatch

import (
	"bao_server/controller"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type JipiaoServer struct{}

var handle JipiaoServer

func NewJipiaoHandle() *JipiaoServer {
	return &handle
}

func (*JipiaoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("handle,uri=", r.URL.Path)
	log.Println("handle,uri=", strings.Trim(r.URL.Path, " "))
	if strings.Trim(r.URL.Path, " ") == "/api/jipiao/post_info" {
		ret := controller.HandlePost(r)
		fmt.Fprint(w, string(ret))
	} else if strings.Trim(r.URL.Path, " ") == "/api/jipiao/search" {
		log.Println("lalalalalalalalallalal")
		ret := controller.JipiaoHandleSearch(r)
		fmt.Fprint(w, string(ret))
	} else if r.URL.Path == "/api/jipiao/update_status" {

	} else {
		log.Println("not found ptah=", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}
}
