package main

import (
	"bao_server/db"
	"bao_server/dispatch"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	jipiaoHandle := dispatch.NewJipiaoHandle()
	serverMux := http.NewServeMux()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//注册关闭服务信号
	go func() {
		sig := <-sigs
		fmt.Println("signal stop server")
		fmt.Println(sig)
		db.Close()
		os.Exit(0)
	}()
	//添加url处理类
	serverMux.Handle("/jipiao/", jipiaoHandle)
	log.Println("server start")
	err := http.ListenAndServe(":8888", serverMux)
	if err != nil {
		log.Fatalln("server start fail")
	}
}
