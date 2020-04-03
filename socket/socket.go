package socket

import (
	"acln.ro/zerocopy"
	"github.com/Code-Fight/golog"
	"io"
	"net"
	"reflect"
)
// tcp data 拆包 并进行连接分发
// reaminData 没成功解包 剩下的数据 参与到下一次的解包中,
// isContuine 一直等待连接，不会关闭下游过来的连接，如果设置为false，并且 isSucc也是false 就会关闭下游的连接
// isSucc 是否成功分析完data，并分配到服务器，通过server返回
// server 返回的 上游服务器连接ip：port 格式
var Unpack func(data []byte) (reaminData []byte,isContuine bool,isSucc bool,server string)


func HandleConnection(conn net.Conn) {

	//声明一个临时缓冲区，用来存储被截断的数据
	var tmpBuffer []byte



	

	//是否已经连接
	isConnected:=false

	// isContuine 一直等待连接，不会关闭下游过来的连接，如果设置为false，并且 isSucc也是false 就会关闭下游的连接
	// isSucc 是否成功分析完data，并分配到服务器，通过server返回
	// server 返回的 上游服务器连接ip：port 格式
	isContuine :=false
	isSucc :=false
	server :=""

	firstDataLen:=0

	buffer := make([]byte, 4096)

	//尝试接收数据
	for {

		if isConnected{

			connClient, connClientErr := net.Dial("tcp", server)

			if connClientErr != nil {
				conn.Close()
				log.Error("connClientErr:",server)
				return
			}
			//把之前的数据给发过去
			connClient.Write(buffer[:firstDataLen])
			//建立0拷贝通道


			 go zerocopy.Transfer(connClient, conn)
			 _,zerocopyError:= zerocopy.Transfer(conn, connClient)
			if zerocopyError!=nil{
				isConnected=false
				return
			}



		}else {
			// 如果没有建立链接，需要先去创建连接
			var err error
			firstDataLen, err = conn.Read(buffer)
			if err != nil {
				if err == io.EOF {
					//对端关闭 对端发送了 FIN过来 请求关闭 这里注意TCP的半关闭
					log.Infof("Client Closed!")

				}
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					log.Errorf("TimeOut close client: %s:",opErr.Addr.String())
				}
				log.Error(conn.RemoteAddr().String(), " connection error: ", err, reflect.TypeOf(err))



				return
			}
			// 从缓冲区中读取数据 并尝试解包
			tmpBuffer, isContuine ,isSucc ,server = Unpack(append(tmpBuffer, buffer[:firstDataLen]...))


			if isSucc&&len(server)>0{
				isConnected = true
			}

			if !isSucc{
				if !isContuine{
					conn.Close()
					return
				}
			}
		}




	}

}




