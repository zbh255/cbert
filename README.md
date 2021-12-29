# cbert
> `cbert`是一个`Go`编写的、简单实现的配置管理平台，它被设计成一个无额外依赖的，轻量级的应用，只需要简单的配置即可运行起来
> `cbert`目前只完成了基本功能，现阶段并不适用于生产环境
> `cbert`目前提供了`Java`客户端: [java-cbert](https://github.com/zbh255/java-cbert),`Go`客户端还在规划中
---
#### 它为什么被设计出来？
为了提供一个在多机环境下管理配置文件和一些自定义错误信息的平台，我自己参与的一些项目也使用了该平台。
####  它是怎么工作的？
`cbert server`与`cbert client`的交互是基于`tcp`自举的二进制协议，了解协议的具体设计规范可以查看这一篇文档: (ProtocolDesign.md)

---
#### 简单地配置，使它可以工作
> 以下面的配置文件为例，了解一些概念可以使你更好的适应自己的需求来调节程序的配置
> - `connection.addr`是服务器的地址
> - `status.addr`是将统计信息输出到的目的地地址
> - `source.dir`是存储用户新增文件的根目录
> - `source.user_confid`是用户配置文件的路径
```json
{

	"source": {

		"dir": "./",

		"user_config": "./user_config.json",

		"request_num": 500,

		"user_num": 100

	},

	"connection": {

		"timeout":500,
		"addr":"127.0.0.1:8080"
	},

	"status": {

		"dest": "console",

		"addr": "127.0.0.1:8080"

	},

	"log": {

		"access_log": "./access.log",

		"error_log": "./error.log"

	}

}
```
#### 编译与运行
> 可以自行的拉取本代码库进行编译或者使用对应平台编译好的文件
```shell
cd ./server && go build main.go -o cbert
mv ./cbert ./source
```
run
```shell
./cbert --config ./project_config
```
