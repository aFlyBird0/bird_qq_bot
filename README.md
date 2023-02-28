# bird_qq_bot

因群里复读怪太多，突发奇想速冲的一个 QQ 机器人。 

项目基于 [Logiase/MiraiGo-Template: A template for MiraiGo](https://github.com/Logiase/MiraiGo-Template) 二次开发。

## 快速开始
把 [application.example.yaml](./application.example.yaml) 复制一份，并重命名为 `application.yaml`，

两种登录方式：

1. 扫码登录，`bot.login-method` 填 `qrcode`。然后程序会在启动时，同时在终端中打印出二维码与在当前目录生成 `qrcode.png` 文件，用手机扫码登录即可。
2. 账号密码登录，`bot.login-method` 填 `common`，并填写 `bot.account` 和 `bot.password`。

注：很可能会触发QQ风控，比如密码登录会提示需要滑块验证，这时候就会自动转成扫码登录，同时扫码登录要求扫码的手机和运行QQ机器人的服务器在同一网络环境在。

## 现有功能
### 群复读撤回
检测群消息是否和最近 10 条内重复（可设定次数），若有，撤回。

### 外卖 roll 点
群内发送「外卖」，即可获得一个随机数，范围为 1-100。寝室拿外卖神器。  

### 考研分数段统计
根据群名片统计各分数段人数，支持多群统计、自定义规则（需要简单改代码）。

统计结果展示方式：
* 图片：把统计结果转成图片后发到群里
* 文字（webserver）
  * localServer：在部署机器人的服务器上启动一个 web 服务，把统计结果以网址的形式发到群里
  * remoteServer：把统计信息推送到远程服务器上（见Remote-Webserver 部署），把统计结果以网址的形式发到群里

详细配置见配置文件。

#### Remote-Webserver 部署（可选）

本地平台编译：
```go
go build -o webserver sidecar/webserver/main.go
./webserver -port 8090
```


macOS 平台交叉编译二进制：
```go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o webserver sidecar/webserver/main.go
./webserver -port 8090
```

todo: 直接用 CI 和 goreleaser 编译二进制文件

### 土味情话
群内发送「宝贝」，即可获得一个随机土味情话。单身狗过节神器。

### 随机@
群内发送「开枪」，随机@一位幸运儿。

### 天行API
需要申请[天行api](https://www.tianapi.com/apiview/142), 模块默认未开启。
群内发送指令，即可被@回复消息  
以下是支持的指令：
* 「早安」
* 「晚安」
* 「舔狗日记」
* 「健康小贴士」

## 未来计划

1. 更多功能
2. CI/CD
