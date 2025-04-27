package RTSP

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)


//TCP客户端
type SocketClientSelf struct {
	FirstType  string						//通用：该socket的类型，对于推流端和播放端都是OPTIONS rtsp，对于get请求，则是get，对于post则是post，对于普通的tcp则是tcp
	conn net.Conn							//通用：
	Close bool								//通用：
	Connstr string							//通用：
	LastRecvBuf   []byte                  	//通用：上次接受的内容

	ConnectTime     int64					//连接时间，单位秒
	DataNum         int64					//数据量
	//下面的变量都是RTSP的推流端和播放端用到的
	Session  string							//播放端和推流端：session信息
	Ready    bool							//播放端和推流端：RTSP通信是否已经完成
	Flag     string							//播放端和推流端：客户端标志，Push或Play
	Channel  string							//播放端和推流端：对应的rtsp通道

	HasSend   bool							//播放端：是否已经向播放端推送服务器缓存的RTP信息

	Players map[string]*SocketClientSelf 	//推流端：该频道播放者的信息
	PlayersLock sync.RWMutex				//推流端：读写锁，Players读写的时候，要加锁
	RecordStart bool						//推流端：当推流端开始推流之后，会启动一个协程，用来给该频道的播放端推送RTP消息
	Sdp      string							//推流端：媒体SDP信息
	MyChan chan []byte						//推流端：每接受到推流端的一个RTP消息，则向
}

type RTSPServer struct {
	ServerPort string
	name string
	MyPush *Push
	Pushers map[string]*SocketClientSelf
	PushersLock sync.RWMutex
	FrameBuffer  int
}
func (this *RTSPServer) Init(name string,Port string,FrameBuffer int){
	this.name=name
	this.Pushers=make(map[string]*SocketClientSelf)
	this.ServerPort=Port
	this.FrameBuffer=FrameBuffer
	if this.FrameBuffer<0 {
		this.FrameBuffer=0
	}
	this.MyPush=new(Push)
	this.MyPush.Init(this)
}

func (this *RTSPServer) ConnHandler(conn net.Conn) {
	//1.conn是否有效
	if conn == nil {
		log.Panic("无效的 socket 连接")
	}
	mystruct :=new(SocketClientSelf)
	mystruct.conn=conn
	mystruct.Ready=false
	mystruct.Connstr=conn.RemoteAddr().String()
	mystruct.Close=false
	mystruct.HasSend=false
	mystruct.FirstType=""
	mystruct.Session=GetRandomString(9)
	mystruct.MyChan=make(chan []byte,10000)
	mystruct.Players=make(map[string]*SocketClientSelf)
	mystruct.RecordStart=false
	mystruct.ConnectTime=time.Now().Unix()
	mystruct.DataNum=0
	//2.新建网络数据流存储结构
	buf := make([]byte, 2048)
	//3.循环读取网络数据流
	for {
		if mystruct.Close{
			break
		}
		cnt, err := conn.Read(buf)
		if cnt == 0 || err != nil {
			this.Close(conn,mystruct)
			break
		}
		mystruct.DataNum=mystruct.DataNum+int64(cnt)
		allbytes := BytesCombine(mystruct.LastRecvBuf, buf[0:cnt])
		mystruct.LastRecvBuf=make([]byte,0)
		this.ToInterface(conn,allbytes,mystruct)
	}
}

func (this *RTSPServer) Close(conn net.Conn,mystruct *SocketClientSelf){
	conn.Close()
	mystruct.Close = true
	if mystruct.Flag=="Play" {
		this.PushersLock.Lock()
		for k,v:=range this.Pushers{
			if k==mystruct.Channel{
				v.PlayersLock.Lock()
				_,ok:=v.Players[mystruct.Connstr]
				if ok{
					delete(v.Players,mystruct.Connstr)
				}
				v.PlayersLock.Unlock()
			}
		}
		this.PushersLock.Unlock()
		fmt.Println("客户端",mystruct.Connstr,"停止播放频道",mystruct.Channel)
	}

	if mystruct.Flag=="Push" {
		mystruct.MyChan<-[]byte("0")
		this.PushersLock.Lock()
		_,ok:=this.Pushers[mystruct.Channel]
		if ok{
			delete(this.Pushers,mystruct.Channel)
		}
		this.PushersLock.Unlock()
		fmt.Println("客户端",mystruct.Connstr,"停止推送频道",mystruct.Channel)
	}
}


//处理接收到的接口消息
func (this *RTSPServer) ToInterface(conn net.Conn,buf []byte,mystruct *SocketClientSelf){
	this.MyPush.Msg(conn,buf,mystruct)
}
func (this *RTSPServer) SendData(conn net.Conn,data []byte,mystruct *SocketClientSelf){
	_,e:=conn.Write(data)
	if e!=nil{
		this.Close(conn,mystruct)
	}
}
//开启ServerSocket
func (this *RTSPServer) Run() {
	//1.监听端口
	cServer, err := net.Listen("tcp", ":"+this.ServerPort)
	if err != nil {
		Log("开启服务器失败，端口号",this.ServerPort,err)
		return
	}
	Log("开启服务器成功，端口号",this.ServerPort)
	for {
		//2.接收来自 client 的连接,会阻塞
		conn, err0 := cServer.Accept()
		if err0 != nil {
			fmt.Println("连接出错")
			continue
		}

		//并发模式 接收来自客户端的连接请求，一个连接 建立一个 conn，服务器资源有可能耗尽 BIO模式
		go this.ConnHandler(conn)
	}
}
