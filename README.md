# bird_qq_bot

因群里复读怪太多，突发奇想速冲的一个 QQ 机器人。 

项目基于 [Logiase/MiraiGo-Template: A template for MiraiGo](https://github.com/Logiase/MiraiGo-Template) 二次开发。

## 快速开始
把 [application.example.yaml](./application.example.yaml) 复制一份，并重命名为 `application.yaml`，
填写账号密码，运行即可。

## 现有功能
### 群复读撤回
检测群消息是否和最近 10 条内重复（可设定次数），若有，撤回。

### 外卖 roll 点
群内发送「外卖」，即可获得一个随机数，范围为 1-100。寝室拿外卖神器。  

## 未来计划

1. 更多功能
2. CI/CD