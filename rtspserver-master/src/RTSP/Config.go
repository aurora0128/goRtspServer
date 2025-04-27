package RTSP

import (
	"os"
)

type Config struct{
	AllContent string
}
func (this * Config) Init(ConfigFileName string)(bool,error){
	file, err := os.Open(ConfigFileName)
	if err != nil {
		return false,err
	}
	defer file.Close()

	fileinfo, err1 := file.Stat()
	if err1 != nil {
		return false,err1
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err2 := file.Read(buffer)
	if err2 != nil {
		return false ,err2
	}
	this.AllContent=string(buffer)
	return true,nil
}
func (this *Config) GetString(path string,key string)(string){
	d:=RegFind(this.AllContent,"[\\s]+\\["+path+"\\][\\w\\W]+?\n[\\s]*"+key+"[\\s]*=[\\s]*([^\\n^\\r]+)")
	if len(d)==2{
		return d[1]
	}
	return ""
}
func (this *Config) GetInt(path string,key string)(int){
	d:=RegFind(this.AllContent,"[\\s]*\\["+path+"\\][\\w\\W]+\n[\\s]*"+key+"[\\s]*=[\\s]*([^\\s]+)")
	if len(d)==2{
		return Str2Int(d[1])
	}
	return -9999
}