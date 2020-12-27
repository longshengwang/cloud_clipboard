# Description
This is a clipboard shared at the same network segment.

It has the functions
- client finds the server automatically. 
- The server has the authentication
- The data encrypted


# How to use
### 1. run from code build
##### 1.1. build the code to exec app
```
//build server
go build server.go  
//build client
go build client.go
```

##### 1.2. exec the app
> No need assign the server ip for client(Client find the server by udp broadcast)


### 2. run from the pre-build app
[github release](https://github.com/longshengwang/cloud_clipboard/releases)


# Tips [ IMPORTANT ]
Only one server can be run at a segment which run the same param.

But multi server can be run with the different params - cw/sw.
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


