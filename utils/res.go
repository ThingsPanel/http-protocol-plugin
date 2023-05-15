package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var Response = new(Resp)

//reponse返回结果
type Resp struct {
	Code string `json:"code"`
	Ts   int64  `json:"ts"`
}

func Result(code string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Resp{
		code,
		time.Now().UnixMicro(),
	})
}
func (r *Resp) OK(ctx *gin.Context) {
	Result("200", ctx)
}
func (r *Resp) Failed(ctx *gin.Context) {
	Result("404", ctx)
}
