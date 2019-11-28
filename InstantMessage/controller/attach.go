package controller

import (
	"../unit"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	_ = os.MkdirAll("./mnt", os.ModePerm)
}

func AttachUpload (writer http.ResponseWriter, request *http.Request) {
	//使用FormFile获取上传文件
	srcfile,head,err := request.FormFile("file")
	if err!= nil{
		unit.RespFail(writer,err)
	}
	//持久化文件到本地
	suffix := ".png"
	//如果前端文件名称包含 .png
	filename := head.Filename
	tmp := strings.Split(filename,".")
	if len(tmp)>1{
		suffix = "."+tmp[len(tmp)-1]
	}
	//如果前端指定filetype
	filetype := request.FormValue("filetype")
	if len(filetype)>0{
		suffix = filetype
	}
	filename = fmt.Sprintf("%d%04d%s",time.Now().Unix(),rand.Int31(),suffix)
	dstfile,err := os.Create("./mnt/"+filename)
	if err!=nil{
		unit.RespFail(writer,err)
		return
	}
	//将源文件内容copy到新文件
	_, err = io.Copy(dstfile, srcfile)
	if err!=nil{
		unit.RespFail(writer,err)
		return
	}
	//将新文件路径转换成url地址
	unit.RespSuccess(writer,"/mnt/"+filename)
}