package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//reponse返回结果
type Response struct {
	Code string      `json:"code"`
	Ts   string      `json:"ts"`
	Data interface{} `json:"data"`
}

func Result(code string, msg string, data interface{}, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		code,
		msg,
		data,
	})
}
func (r *Response) OK(data interface{}, ctx *gin.Context) {
	Result("200", "success", data, ctx)
}
func (r *Response) Failed(data interface{}, ctx *gin.Context) {
	Result("404", "error", data, ctx)
}
func (r *Response) FailWithDetailed(data interface{}, msg string, ctx *gin.Context) {
	Result("401", msg, data, ctx)
}

func (r *Response) OkWithData(data interface{}, msg string, ctx *gin.Context) {
	Result("200", "success", data, ctx)
}
