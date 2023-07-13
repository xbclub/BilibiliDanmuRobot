### 本项目概述
    本项目是基于github.com/k-si/bilibili_live项目二开完成
### 使用教程
[点此访问使用教程](https://www.yuque.com/yuqueyonghu3xsgin/igligh/rpr4oslh4nwt2pwv?singleDoc#)
### 功能列表
- [x] 弹幕欢迎（普通用户，舰长用户）
    - [x] 指定用户欢迎
    - [x] 分时段欢迎词
    - [x] 自定义欢迎词
    - [x] 欢迎黑名单（模糊匹配，精准匹配）
- [x] 关注答谢
- [x] 礼物答谢
    - [x] 自定义礼物答谢延迟
- [x] AI机器人
    - [x] 两个ai机器人接口
    - [x] `@帮助`查看机器人触发关键字
- [x] PK相关
    - [x] PK输出对手信息
    - [x] 对方用户串门识别
- [x] 秒级定时弹幕
- [ ] 自动更新（GUI）
- [ ] 弹幕指令控制

###  构建环境
 * golang version 1.20+
### 构建指令
```
 # 获取构建依赖
 go mod download
 # 构建软件
 go build cli/bilidanmaku.go
```
### 鸣谢
- github.com/k-si/bilibili_live
