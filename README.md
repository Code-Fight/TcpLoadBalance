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



