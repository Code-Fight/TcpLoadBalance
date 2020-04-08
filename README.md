# TcpLoadBalance
根据tcp数据包内容，进行自定义的负载均衡。
![](doc/1.jpg)
##由来
因为一个业务场景的需要，需要对tcp进行分流，并且要尽可能快的响应，但是数据包还比较大。
用现有的nginx或者openresty之类的也可以实现，但是基础的功能都是按照tcp的src来进行balance。
由于我们的业务场景特殊，需要根据数据包内容进行分流，当然用现有的中间件也可以进行分流，像前面提到的nginx和openresty
但是，都需要自己去一些基于Lua的开发，由于tcp都是长连接，那么既需要进行数据包的拆包balance，又需要状态的维持
从openresty github 中有看到别人就类似的场景咨询作者，但是作者不太推荐这种应用场景

##要做什么
1.首先最基础的，肯定是先建立一个tcp server来接收下游的client

2.根据下游client的data来进行数据解包和balance

3.与后端建立链接，并保持链接

上面提到的是最基础的实现

**除此之外，我们还需要稍微扩展一下**

1.后端服务健康监测

2.注册发现

3.故障转移

再之外还有很多，暂时先不提，通过上面的几部分，基本上可以达到应用的基本需求。

##注意点
1.在client 结束balance后应该尽可能的减少syscall，采用linux的 zero copy来尽可能的减少

2.提供注册发现的api，通过要注意可能造成client的迁移

##使用




---
# TcpLoadBalance
Custom TCP load balancing, TCP parsing and forwarding

##About the TcpLoadBalance
![](doc/1.jpg)

##How to use
Open the 'main.go' ,then edit the 'InitTcpServer' method
Refer to the code below
```go
func InitTcpServer(port string) {
	.....

	socket.Unpack=Unpack

	.....
}
```
Modfiy the 'socket.Unpack=Unpack',you should define the 'Unpack' method
and refer to your procotol.

Refer to the code below
```go
func Unpack(buffer []byte) (remainData []byte,isContuine bool,isSucc bool,server string) {
        ......
	return buffer,false,true,"192.168.2.115:12048"
}
```
###'Unpack' return params:

**reaminData**  : if the method got a your define 'protocol data' or 'buffer data' can't provide your need data,you should return the 'remain data'

**isContuine**  : if the method got a your define 'protocol data' ,return false,then will stop exec this method  

**isSucc**      :if the method got a your define 'protocol data',return true,and you should return a server ip 

**server**      :if 'isSucc' is true,this param will provide a server ip for to connect

### Summary
About the 'Unpack' method your should do 2 part work
1.Got a your need 'protocol data' by the 'buffer data'
2.Assign backend server ip by your 'protocol data'

##TODO
benchmark



