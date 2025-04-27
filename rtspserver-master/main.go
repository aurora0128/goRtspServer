package main

import (
	"RTSP"
)

func main(){
	MyConfig:=new(RTSP.Config)
	MyConfig.Init("./Config.ini")
	ob:=new(RTSP.RTSPServer)
	ob.Init("RTSP服务器",MyConfig.GetString("RTSP","Port"),MyConfig.GetInt("RTSP","FrameBuffer"))
	ob.Run()
}
