package dispatch

import (
	"bao_server/controller"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type CommonServer struct{}

var common_handle CommonServer

func NewCommonHandle() *CommonServer {
	return &common_handle
}

func (*CommonServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("handle,uri=", r.URL.Path)
	if strings.Trim(r.URL.Path, " ") == "/api/common/post_info" {
		ret := controller.CommonHandlePost(r)
		fmt.Fprint(w, string(ret))
	} else if  strings.Trim(r.URL.Path, " ") == "/api/common/search" {
		ret := controller.CommonHandleSearch(r)
		fmt.Fprint(w, string(ret))
	} else if r.URL.Path == "/api/common/update_status" {

	} else {
		log.Println("not found ptah=", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}
}
