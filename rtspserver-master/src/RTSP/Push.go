package RTSP

import (
	"fmt"
	"net"
)
type Push struct {
	MyServer *RTSPServer
}
func (this *Push) Init(cur *RTSPServer){
	this.MyServer=cur
}
func(this *Push)Msg(conn net.Conn,buf []byte,mystruct *SocketClientSelf)(int){
	if len(buf)<12{
		return 0
	}
	if mystruct.FirstType==""{
		mystruct.FirstType=string(buf[0:12])
	}
	if mystruct.FirstType!="OPTIONS rtsp" {
		return 0
	}
	if len(buf)<10{
		return 0
	}
	ff:=this.Data(conn,buf,mystruct)
	if ff==0{
		this.Options(conn,buf,mystruct)
		this.ANNOUNCE(conn,buf,mystruct)
		this.SETUP(conn,buf,mystruct)
		this.RECORD(conn,buf,mystruct)
		this.PLAY(conn,buf,mystruct)
		this.DESCRIBE(conn,buf,mystruct)
	}

	return 1
}
//处理Options请求
/*
OPTIONS rtsp://192.168.1.201:5545/2_1 RTSP/1.0
CSeq: 1
User-Agent: Lavf58.37.100
*/
func(this *Push)Options(conn net.Conn,msgbytes []byte,mystruct *SocketClientSelf)(int){
	strbuf:=string(msgbytes)
	if strbuf[0:7]!="OPTIONS" {
		return 0
	}
	mystruct.LastRecvBuf=msgbytes
	post:= RegFind(strbuf,"^OPTIONS rtsp://[^:]+?:[\\d]+/([\\w\\W]+?) RTSP/([0-9.]+)[\\s]+CSeq:[\\s]*([0-9]+)[\\w\\W]+[\\s]+")
	if len(post)!=4{
		return 0
	}


	mystruct.Channel=post[1]
	sendstr:=fmt.Sprintf("RTSP/1.0 200 OK\nCSeq: %s\nSession: %s\nPublic: DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE, OPTIONS, ANNOUNCE, RECORD\n\n",post[3],mystruct.Session)
	this.MyServer.SendData(conn,[]byte(sendstr),mystruct)
	mystruct.LastRecvBuf=make([]byte,0)
	return 1
}

//处理ANNOUNCE请求
/*
ANNOUNCE rtsp://192.168.1.201:5545/2_1 RTSP/1.0
Content-Type: application/sdp
CSeq: 2
User-Agent: Lavf58.37.100
Session: ZTnZLWlGg
Content-Length: 296

v=0
o=- 0 0 IN IP4 127.0.0.1
s=No Name
c=IN IP4 192.168.1.201
t=0 0
a=tool:libavformat 58.37.100
m=video 0 RTP/AVP 96
a=rtpmap:96 H264/90000
a=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z00AKp2oHgCJ+WbgICAoAAADAAgAAAMBlCA=,aO48gA==; profile-level-id=4D002A
a=control:streamid=0
*/
func(this *Push)ANNOUNCE(conn net.Conn,msgbytes []byte,mystruct *SocketClientSelf)(int){
	strbuf:=string(msgbytes)
	if strbuf[0:8]!="ANNOUNCE" {
		return 0
	}

	mystruct.LastRecvBuf=msgbytes
	post:= RegFind(strbuf,"^ANNOUNCE [\\w\\W]+?CSeq:[\\s]*([0-9]+)[\\w\\W]*[\\s]+Content-Length:[\\s]*([0-9]+)[^\\r^\\n]*?([\\w\\W]*)")
	if len(post)!=4{
		return 0
	}
	ln:=Str2Int(post[2])

	if len(post[3])<=4{
		mystruct.LastRecvBuf=msgbytes
		return 0
	}
	post[3]=post[3][4:]
	if ln>len(post[3]){
		mystruct.LastRecvBuf=msgbytes
		return 0
	}

	if ln<len(post[3]){
		mystruct.LastRecvBuf=[]byte(post[3][ln:])
	}
	mystruct.Sdp=post[3][0:ln]

	mystruct.Flag="Push"
	this.MyServer.PushersLock.Lock()
	_,ok:=this.MyServer.Pushers[mystruct.Channel]
	if ok{
		this.MyServer.PushersLock.Unlock()
		return 0
	}
	this.MyServer.Pushers[mystruct.Channel]=mystruct
	this.MyServer.PushersLock.Unlock()
	sendstr:=fmt.Sprintf("RTSP/1.0 200 OK\nCSeq: %s\nSession: %s\n\n",post[1],mystruct.Session)

	this.MyServer.SendData(conn,[]byte(sendstr),mystruct)
	mystruct.LastRecvBuf=make([]byte,0)
	return 1
}
//处理SETUP请求
/*
SETUP rtsp://192.168.1.201:5545/2_1/streamid=0 RTSP/1.0
Transport: RTP/AVP/TCP;unicast;interleaved=0-1;mode=record
CSeq: 3
User-Agent: Lavf58.37.100
Session: ZTnZLWlGg


RTSP/1.0 200 OK
CSeq: 3
Session: ZTnZLWlGg
Transport: RTP/AVP/TCP;unicast;interleaved=0-1;mode=record


*/
func(this *Push)SETUP(conn net.Conn,msgbytes []byte,mystruct *SocketClientSelf)(int){
	strbuf:=string(msgbytes)
	if strbuf[0:5]!="SETUP" {
		return 0
	}
	mystruct.LastRecvBuf=msgbytes
	post:= RegFind(strbuf,"^SETUP [\\w\\W]+? RTSP/[0-9.]+[\\s]+Transport:[\\s]*([^\\n]+?)[\\s]+CSeq:[\\s]*([0-9]+)[\\w\\W]+")
	if len(post)!=3{
		return 0
	}

	sendstr:=fmt.Sprintf("RTSP/1.0 200 OK\nCSeq: %s\nSession: %s\nTransport: %s\n\n",post[2],mystruct.Session,post[1])
	this.MyServer.SendData(conn,[]byte(sendstr),mystruct)
	mystruct.LastRecvBuf=make([]byte,0)
	return 1
}


//处理SETUP请求
/*
RECORD rtsp://192.168.1.201:5545/2_1 RTSP/1.0
Transport: RTP/AVP/TCP;unicast;interleaved=0-1;mode=record
CSeq: 3
User-Agent: Lavf58.37.100
Session: ZTnZLWlGg

RTSP/1.0 200 OK
CSeq: 3
Session: ZTnZLWlGg
Transport: RTP/AVP/TCP;unicast;interleaved=0-1;mode=record


*/
func(this *Push)RECORD(conn net.Conn,msgbytes []byte,mystruct *SocketClientSelf)(int){
	strbuf:=string(msgbytes)
	if strbuf[0:6]!="RECORD" {
		return 0
	}
	mystruct.LastRecvBuf=msgbytes
	post:= RegFind(strbuf,"^RECORD [\\w\\W]+? RTSP/[0-9.]+[\\w\\W]+?CSeq:[\\s]*([0-9]+)")
	if len(post)!=2{
		return 0
	}
	mystruct.Ready=true
	sendstr:=fmt.Sprintf("RTSP/1.0 200 OK\nSession: %s\nCSeq: %s\n\n",mystruct.Session,post[1])
	this.MyServer.SendData(conn,[]byte(sendstr),mystruct)

	if mystruct.RecordStart==false{
		mystruct.RecordStart=true
		go this.PushData(mystruct)
	}
	mystruct.LastRecvBuf=make([]byte,0)
	fmt.Println("客户端",mystruct.Connstr,"开始推送频道",mystruct.Channel)
	return 1
}

func(this *Push) PushData(mystruct *SocketClientSelf){
	TempData:=make([][]byte,0)
	for{
		d:=<-mystruct.MyChan

		if len(d)==1{
			break
		}
		if this.MyServer.FrameBuffer>0{
			if len(TempData)==this.MyServer.FrameBuffer{
				TempData = append(TempData[:0], TempData[1:]...)
			}
			TempData = append(TempData, d)
		}

		Removes:=make(map[string]string)
		Removes_Flag:=false
		mystruct.PlayersLock.RLock()
		for _,v:=range mystruct.Players{
			if  this.MyServer.FrameBuffer>0 && mystruct.HasSend==false{
				mystruct.HasSend=true
				for _,dd:=range TempData{
					_,e:=v.conn.Write(dd)
					if e!=nil{
						Removes[v.Connstr]=v.Connstr
						Removes_Flag=true
						break
					}
				}
			}
			_,e:=v.conn.Write(d)
			if e!=nil{
				Removes[v.Connstr]=v.Connstr
				Removes_Flag=true
			}
		}

		mystruct.PlayersLock.RUnlock()

		if Removes_Flag{
			mystruct.PlayersLock.Lock()
			for key,_:=range Removes{
				_,okk:=mystruct.Players[key]
				if okk{
					delete(mystruct.Players,key)
				}
			}
			mystruct.PlayersLock.Unlock()
		}
	}

}

//处理DESCRIBE请求
/*
DESCRIBE rtsp://192.168.1.201:5545/2_1 RTSP/1.0
Accept: application/sdp
CSeq: 2
User-Agent: Lavf58.12.100
Session: YXN_wZ_GR

RTSP/1.0 200 OK
Session: YXN_wZ_GR
Content-Length: 296
CSeq: 2

v=0
o=- 0 0 IN IP4 127.0.0.1
s=No Name
c=IN IP4 192.168.1.201
t=0 0
a=tool:libavformat 58.37.100
m=video 0 RTP/AVP 96
a=rtpmap:96 H264/90000
a=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z00AKp2oHgCJ+WbgICAoAAADAAgAAAMBlCA=,aO48gA==; profile-level-id=4D002A
a=control:streamid=0
*/
func(this *Push)DESCRIBE(conn net.Conn,msgbytes []byte,mystruct *SocketClientSelf)(int){
	strbuf:=string(msgbytes)
	if strbuf[0:8]!="DESCRIBE" {
		return 0
	}
	mystruct.LastRecvBuf=msgbytes
	post:= RegFind(strbuf,"^DESCRIBE rtsp://[\\w\\W]+?:[0-9]+/([\\w\\W]+?) RTSP/[0-9.]+[\\w\\W]+?CSeq:[\\s]*([0-9]+)")

	if len(post)!=3{
		return 0
	}

	mystruct.Flag="Play"
	this.MyServer.PushersLock.RLock()
	v,ok:=this.MyServer.Pushers[post[1]]
	this.MyServer.PushersLock.RUnlock()

	if ok {
		sendstr:=fmt.Sprintf("RTSP/1.0 200 OK\nSession: %s\nContent-Length: %d\nCSeq: %s\n\n%s",mystruct.Session,len(v.Sdp),post[2],v.Sdp)
		this.MyServer.SendData(conn,[]byte(sendstr),mystruct)
	}

	mystruct.LastRecvBuf=make([]byte,0)
	return 1
}



//处理PLAY请求
/*
PLAY rtsp://192.168.1.201:5545/2_1 RTSP/1.0
Range: npt=0.000-
CSeq: 4
User-Agent: Lavf58.12.100
Session: YXN_wZ_GR

RTSP/1.0 200 OK
Session: YXN_wZ_GR
Range: npt=0.000-
CSeq: 4
*/
func(this *Push)PLAY(conn net.Conn,msgbytes []byte,mystruct *SocketClientSelf)(int){
	strbuf:=string(msgbytes)
	if strbuf[0:4]!="PLAY" {
		return 0
	}
	mystruct.LastRecvBuf=msgbytes
	post:= RegFind(strbuf,"^PLAY[\\w\\W]+?CSeq:[\\s]*([0-9]+)[\\w\\W]+")
	if len(post)!=2{
		return 0
	}
	mystruct.Ready=true
	this.MyServer.PushersLock.Lock()
	vv,ok:=this.MyServer.Pushers[mystruct.Channel]
	if ok{
		vv.PlayersLock.Lock()
		vv.Players[mystruct.Connstr]=mystruct
		vv.PlayersLock.Unlock()
	}
	this.MyServer.PushersLock.Unlock()
	sendstr:=fmt.Sprintf("RTSP/1.0 200 OK\nSession: %s\nRange: npt=0.000-\nCSeq: %s\n\n",mystruct.Session,post[1])
	this.MyServer.SendData(conn,[]byte(sendstr),mystruct)
	fmt.Println("客户端",mystruct.Connstr,"开始播放频道",mystruct.Channel)
	mystruct.LastRecvBuf=make([]byte,0)
	return 1
}


//处理正常的RTP数据帧
func(this *Push)Data(conn net.Conn,msgbytes []byte,mystruct *SocketClientSelf)(int){
	if mystruct.Ready==false || mystruct.Flag!="Push"{
		return 0
	}
	if len(msgbytes)<10 {
		mystruct.LastRecvBuf=msgbytes
		return -1
	}
	mystruct.LastRecvBuf=make([]byte,0)
	msgbytes0:=msgbytes
	sendbytes:=make([]byte,0)
	for{
		if len(msgbytes0)<4{
			mystruct.LastRecvBuf=msgbytes0
			break
		}
		ln:=BytesToInt16(msgbytes0[2:4])
		lnn:=int(ln)
		bl:=BytesToInt8(msgbytes0[1:2])

		if bl>=200 && bl<=207{
			//RTCP包，推送端推送会话质量信息
			fmt.Println(bl)
		}

		if lnn<0 || lnn>20000 {
			mystruct.LastRecvBuf=make([]byte,0)
			break
		}
		if lnn+4>len(msgbytes0) {
			mystruct.LastRecvBuf=msgbytes0
			break
		}
		if lnn+4==len(msgbytes0) {

			sendbytes=BytesCombine(sendbytes,msgbytes0)
			mystruct.LastRecvBuf=make([]byte,0)
			break
		}
		if lnn+4<len(msgbytes0) {
			sendbytes=BytesCombine(sendbytes,msgbytes0[0:4+ln])
			msgbytes0=msgbytes0[4+ln:]
		}
	}
	if len(sendbytes)>0 {
		mystruct.MyChan<-sendbytes
	}

	return 1
}
