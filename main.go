package main

import (
	"TcpLoadBalance/balance"
	"TcpLoadBalance/linuxRoute"
	"TcpLoadBalance/socket"
	"TcpLoadBalance/units"
	"fmt"
	log "github.com/Code-Fight/golog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
)

//定义backend
//如果要使用别的balance算法，继承Balancer，然后实现即可
//参考HashBalancer
var bl balance.Balancer = &balance.HashBalancer{}
var route *linuxRoute.Route
func main() {

	//注入版本信息
	VersionInit()

	//config
	units.ConfigInit()
	port:=units.GetPort()

	//balance初始化
	bl.BalanceInit()

	//初始化LinuxRoute
	route = linuxRoute.NewRoute()

	//启动server
	go InitTcpServer(port)

	//debug
	httpErr:= http.ListenAndServe("0.0.0.0:16061", nil)
	if httpErr!=nil{
		fmt.Println("HTTP Server ERR:",httpErr.Error())
	}
}

func CheckError(err error) {
	if err != nil {
		log.Errorf("Start Error: %s", err.Error())
		os.Exit(1)
	}
}

func InitTcpServer(port string) {
	//初始化 server
	netListen, err := net.Listen("tcp", "0.0.0.0:"+port)
	CheckError(err)
	defer netListen.Close()
	fmt.Print("TCP Server Running on Port:",port+"\r\n")

	socket.Unpack=Unpack

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		log.Debug(conn.RemoteAddr().String(), " tcp connect success")
		go socket.HandleConnection(conn,route)
	}
}

//解包
func Unpack(buffer []byte) (remainData []byte,isContuine bool,isSucc bool,server string) {
	//return buffer,false,true,"192.168.2.115:12048"

	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {

		// 第一个字符是
		// * 标记
		if buffer[i] == 0x2a {
			if length-i < 7 {
				//可能被拆包了，头部太短，等待下一次处理
				break
			}
			dataLen := int(units.BytesToUint32(buffer[i+1:i+5]))
			messageLength := dataLen + 7
			if messageLength > length-i {
				//数据区长度不够  无法解析 等待下一次
				break
			}
			data := buffer[i : i+messageLength]

			serverIp:=AssignBackend(data)

			return data,false,true,serverIp




			i = messageLength+i-1

		}

	}
	// 如果首位不是* 并且 一直没找到* 那么只能丢弃掉
	if i == length {
		//log.Printf("数据丢弃:%x",buffer)
		return make([]byte, 0),true,false,""

	}
	return buffer[i:],true,false,""
}

//根据protocol data进行分配后端服务器
func AssignBackend(data []byte) string {


	//根据数据包的19-21位进行hash
	//这里可以根据自己的实际需要进行选择
	return bl.GetNode(string(data[19:]))
}

