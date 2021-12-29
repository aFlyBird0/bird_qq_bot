# bird_qq_bot

因群里复读怪太多，突发奇想速冲的一个 QQ 机器人。 

项目基于 [Mrs4s/MiraiGo: qq协议的golang实现, 移植于mirai](https://github.com/Mrs4s/MiraiGo) 二次开发。

## 现有功能
### 防复读
检测群消息是否和最近 10 条内重复（可设定次数），若有，撤回。

### 外卖 roll 点
群内发送「外卖」，即可获得一个随机数，范围为 1-100。寝室拿外卖神器。  