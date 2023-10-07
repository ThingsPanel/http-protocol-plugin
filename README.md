# http-procotol-plugin

## 功能

独立程序，通过http接收数据转发至thingspanel的mqtt对应的主题。

设备是mqtt以外的协议除了可以从规则引擎接入，也可以开发协议插件服务接入到ThingsPanel

## 插件如何注册到平台

手动注册

1. 点击 `应用管理`->`接入协议`->`注册插件`
2. 填入插件信息
   **名称**：创建设备时，会显示在选择协议下拉框中
   **设备类型**：必填，选直连设备
   **接入地址**：插件服务的ip地址和端口（设备对接,作为平台中的提示信息）
   **HTTP服务地址**：插件服务的ip地址和端口（必填，供平台后端调用）
   **插件订阅主题前缀**： （必填）plugin/http/

## 结构图

![结构图](./images/协议插件.png)

![时序图](images/时序图.png)

## 插件表单

`./form_config.json` （表单规则详情请参考modbus-protocol-plugin案例） ThingsPanel前端通过 `/api/form/config`接口获取表单配置，生成子设备的表单，用户填写的表单数据会出现在ThingsPanel提供的 `/api/plugin/device/config`接口返回的数据中。

## 交换数据相关

设备post发送json数据至插件

api/device/AccessToken/attributes

（api前面需链接为http协议插件所部署地址，例如http://127.0.0.1:9988/api/device/AccessToken/attributes;AccessToken为在ThingsPanel平台，添加设备绑定插件时获得的，或自定义的值）

转发至mqtt的device/attributes 主题

json数据格式：

```
{
    "temp": 18.5,
    "low": 40,
    ...
}
```

响应体：

```
{
  "code": 200,//200：成功，404：失败
  "ts": //时间戳 微秒
}
```

插件发送数据至ThingsPanel平台

1.直连设备

mqtt用户：root （使用thingspanel-go配置文件中的用户名和密码）

发布主题：device/attributes

报文规范：{"token":device_token,"values":{key:value...}}

或自定义报文：{"token":device_token,"values":自定义报文}

token：设备AccessToken或子设备AccessToken

2.在线离线通知

mqtt用户：root （使用thingspanel-go配置文件中的用户名和密码）

发布主题：device/status

报文规范：{"accessToken":accessToken,"values":{"status":status}}

accessToken:设备或网关连接时送来的密钥

status: "0"-离线 "1"-上线

或自定义报文：自定义报文
