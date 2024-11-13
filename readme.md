### 本项目概述
    本项目是基于github.com/k-si/bilibili_live项目二开完成
### 使用教程
[点此访问使用教程](https://www.yuque.com/yuqueyonghu3xsgin/igligh/rpr4oslh4nwt2pwv?singleDoc#)
### 功能列表
- [x] 弹幕欢迎（普通用户，舰长用户）
    - [x] 指定用户欢迎
    - [x] 荣耀等级欢迎
    - [x] 分时段欢迎词
    - [x] 自定义欢迎词
    - [x] 欢迎黑名单（模糊匹配，精准匹配）
- [x] 关注答谢
- [x] 分享感谢
- [x] 下播公告
- [x] 礼物答谢
    - [x] 自定义礼物答谢延迟
    - [x] 盲盒统计
- [x] AI机器人
    - [x] 两个ai机器人接口
    - [x] `@帮助`查看机器人触发关键字
- [x] PK相关
    - [x] PK输出对手信息
    - [x] 对方用户串门识别
- [x] 秒级定时弹幕
- [x] 自动更新（GUI）
- [x] 弹幕指令控制
### Docker
```bash
# 创建配置文件
mkdir -p <your path>/etc
wget -O <your path>/etc/bilidanmaku-api.yaml  https://github.moeyy.xyz/https://raw.githubusercontent.com/xbclub/BilibiliDanmuRobot/master/etc/bilidanmaku-api.yaml
# 启动容器
docker run -itd --name bilibilidanmurobot --restart=always -v <your path>:/app/data xbclub/bilibilidanmurobot:latest
# 扫码登录
docker logs bilibilidanmurobot
```
###  构建环境
 * golang version 1.21+
### 构建指令
```
 # 获取构建依赖
 go mod download
 # 构建软件
 go build cli/bilidanmaku.go
```
#### linux使用musl工具链编译musl全静态可执行文件指令
##### 以下以x64为例
```bash
# 指定编译工具链，可通过https://musl.cc/下载
export CC=/home/user/下载/x86_64-linux-musl-native/bin/gcc
export CXX=/home/user/下载/x86_64-linux-musl-native/bin/g++
export CFLAGS="-I/home/user/下载/x86_64-linux-musl-native/include"
export LDFLAGS="-I/home/user/下载/x86_64-linux-musl-native/lib -Wl,-Bstatic"

# 禁用 Cgo，避免引入动态库
export CGO_ENABLED=0                                      
export GOOS=linux
# 指定编译架构
export GOARCH=amd64

 # 获取构建依赖
 go mod download
 # 构建软件
 go build cli/bilidanmaku.go
```
### 鸣谢
- github.com/k-si/bilibili_live
