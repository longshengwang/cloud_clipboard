# 介绍
这是一个局域网内粘贴板共享工具软件(公网需要对安全方面做处理，暂时未做)

功能清单
- 客户端能自动找到服务端
- 服务端可以设定密码
- 客户端和服务端通信加密


# 如何使用
### 1. 从源码安装
##### 1.1. 编译代码
```
//build server
go build server.go  
//build client
go build client.go
```

##### 1.2. 执行应用
> 启用客户端时，不需要指定服务端的IP地址，因为用了UDP广播来查找服务端


### 2. 直接运行程序
[github release](https://github.com/longshengwang/cloud_clipboard/releases)


# 注意 [ IMPORTANT ]
在一个局域网下只能运行一个相同参数的服务端。

如果设置不同的参数（cw/sw），是可以同时运行多个的服务端的
```
Usage of ./server_macos:
  -auth string   
    	Server Auth Password. Cannot more than 32 char(256 bit) (default "cloud_clipboard_password")
  -cw string        
    	Client Hello Word (default "Hello, is my clipboard?")
  -discoveryPort int
    	Discovery Service Port (default 9266)
  -port int
    	Server Port (default 5166)
  -sw string
    	Server Hello Word (default "Hey, brother, you are at the home of clipboard server.")
```


