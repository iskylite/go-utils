
# go-utils

my personal golang utils modules

## ilog

日志包，参考标准包log和[goframe](https://github.com/gogf/gf)框架的内部包glog进行实现的，只要是定制成个人喜欢的日志格式，减少配置，提示参考glog的日志分段方法设计实现了日志文件大小分段。

本日志包不考虑性能，毕竟要考虑性能可以使用zap啊。

具体使用方法可以参考例子文件[example.go](./ilog/_example/example.go)

## itool

主要是一些工具函数的封装，会持续添加

## iparse

主要是热加载toml日志解析包，更新中...