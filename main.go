package main

import (
	"bao_server/db"
	"bao_server/dispatch"
	"bao_server/mylog"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	jipiaoHandle := dispatch.NewJipiaoHandle()
	wexinHandle := dispatch.NewWeixinHandle()
	serverMux := http.NewServeMux()
	err := mylog.InitLog()
	if err != nil {
		fmt.Println("init log fail")
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//注册关闭服务信号
	go func() {
		sig := <-sigs
		fmt.Println("signal stop server")
		fmt.Println(sig)
		db.Close()
		mylog.CloseLog()
		os.Exit(0)
	}()
	err = db.ConnDb()
	if err != nil {
		log.Panicln("db init fail")
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	//添加url处理类
	serverMux.Handle("/", wexinHandle)
	serverMux.Handle("/api/jipiao/", jipiaoHandle)
	log.Println("server start")
	err = http.ListenAndServe(":8888", serverMux)
	if err != nil {
		log.Fatalln("server start fail")
	}
}
