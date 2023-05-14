package controller

import (
	"http-procotol-plugin/service"
	"http-procotol-plugin/utils"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

//相关视图，跳转至services里对应服务
type TpController struct{}

var Response = utils.Uutils.Response

//获取表单配置
func (w *TpController) Config(ctx *gin.Context) {
	if data, err := service.TpSer.Config(); err != nil || data == "" {
		Response.Failed("获取失败!", ctx)
	} else {
		Response.OK(data, ctx)
	}
}

//更新设备
func (w *TpController) UpdateDevice(ctx *gin.Context) {
	var device utils.Device
	_ = ctx.ShouldBindJSON(&device)
	if err := service.TpSer.UpdateDevice(device); err != nil {
		Response.Failed("更新失败!", ctx)
	} else {
		Response.OK("更新成功!", ctx)
	}
}

//新增设备
func (w *TpController) AddDevice(ctx *gin.Context) {
	var device utils.Device
	err := ctx.ShouldBindJSON(&device)
	if err != nil {
		Response.Failed("设备无效!", ctx)
	}
	if err := service.TpSer.AddDevice(device); err != nil {
		Response.Failed("新增失败!", ctx)
	} else {
		Response.OK("新增成功!", ctx)
	}
}

//删除设备
func (w *TpController) DeleteDevice(ctx *gin.Context) {
	var id utils.Device
	_ = ctx.ShouldBindJSON(&id)
	if err := service.TpSer.DeleteDevice(id); err != nil {
		Response.Failed("删除失败!", ctx)
	} else {
		Response.OK("删除成功!", ctx)
	}

}

//接收属性
func (w *TpController) Attributes(ctx *gin.Context) {
	accesstoken := ctx.Param("accesstoken")
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	if err := service.TpSer.Attributes(accesstoken, body); err != nil {
		Response.Failed("发送失败!", ctx)
	} else {
		Response.OK("发送成功!", ctx)
	}

}

//接收事件数据
func (w *TpController) Event(ctx *gin.Context) {
	accesstoken := ctx.Param("accesstoken")
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	if err := service.TpSer.Event(accesstoken, body); err != nil {
		Response.Failed("发送失败!", ctx)
	} else {
		Response.OK("发送成功!", ctx)
	}
}

//接收命令响应数据
func (w *TpController) CommandReply(ctx *gin.Context) {
	accesstoken := ctx.Param("accesstoken")
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	if err := service.TpSer.CommandReply(accesstoken, body); err != nil {
		Response.Failed("发送失败!", ctx)
	} else {
		Response.OK("发送成功!", ctx)
	}
}
