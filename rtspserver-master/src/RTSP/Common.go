package RTSP

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)
func BytesToInt8(bys []byte) int {
	bin_buf := bytes.NewBuffer(bys)
	var x int8
	binary.Read(bin_buf, binary.BigEndian, &x)
	return int(x)
}
func BytesToInt16(bys []byte) int16 {
	bin_buf := bytes.NewBuffer(bys)
	var x int16
	binary.Read(bin_buf, binary.BigEndian, &x)
	return x
}
func Strtime()string{
	return time.Now().Format("2006-01-02 15:04:05")
}
func Log(a ...interface{}){
	dd:=make([]interface{},0)
	dd=append(dd,Strtime())
	for _,k := range a {
		dd=append(dd,k)
	}
	fmt.Println(dd...)
}
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func  GetRandomString(l int) string {
	bytes0 := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes0[r.Intn(len(bytes0))])
	}
	return string(result)
}
func  RegFind(buf string,regstr string)([]string){
	//解析正则表达式，如果成功返回解释器
	//res := make([][][]byte,0)
	ret := make([]string,0)
	re, _  := regexp.Compile(regstr)
	//根据规则提取关键信息
	res := re.FindAllStringSubmatch(buf, -1)
	for _,v :=range res{
		for _,v1 := range v{
			ret = append(ret,v1)
		}
	}
	return ret
}
func Str2Int(str string)int{
	val, _ := strconv.ParseInt(str, 10, 32)
	return int(val)
}