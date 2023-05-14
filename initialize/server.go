package initialize

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type server interface {
	ListenAndServe() error
}

func initServer(addr string, router *gin.Engine) server {
	return &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
