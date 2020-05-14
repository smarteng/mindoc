#!/usr/bin/env bash
# 获取当前时间
BuildTime=`date +'%Y-%m-%d %H:%M:%S'`
# 获取 Go 的版本
BuildGoVersion=`go version`

# 将以上变量序列化至 LDFlags 变量中
LDFlags=" \
    -X 'github.com/smarteng/mindoc/conf.BUILD_TIME=${BuildTime}' \
    -X 'github.com/smarteng/mindoc/conf.GO_VERSION=${BuildGoVersion}' \
"

go build -ldflags "-w" -ldflags "$LDFlags" 

echo 'build done.'