package initialize

import (
	"http-procotol-plugin/global"
	"http-procotol-plugin/routers"
	"log"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	//gin.SetMode(gin.ReleaseMode) //开启生产模式
	r := gin.Default()
	routers.RegisterRouter(r)
	s := initServer(global.Conf.Addr, r)
	log.Printf("启动http服务，地址：%s\n", global.Conf.Addr)
	log.Fatal(s.ListenAndServe().Error())
}
