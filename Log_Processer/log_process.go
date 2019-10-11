package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
)

type Reader struct {
	logpath string
}

type Writer struct {
	influxDBDsn string
}

type LogProcess struct {
	rc chan []byte
	wc chan string
	logreader Reader
	infwriter Writer
}

type Read interface {
	ReadFromFile()
}

type Write interface {
	WirteToInfluxDB()
}

type LogMessage struct {
	LocalTime 						 	time.Time
	BytesLength 					 	int
	Path, Method, Scheme, StatusCode 	string
	UpstreamTime, RequestTime  		 	float64
}

func (r *Reader) ReadFromFile(rc chan []byte)  {
	//读取模块
	//打开文件
	f, err := os.Open(r.logpath)
	if err != nil {
		panic(fmt.Sprintf("open file error:%s",err.Error()))
	}

	//从文件尾开始读取文件内容
	_, _ = f.Seek(0, 2)  //光标放到最后
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500*time.Millisecond)
			continue
		}
		if err != nil {
			panic(fmt.Sprintf("log file readed failture:%s", err.Error()))
		}
		rc <- line[:len(line)-1]
	}
}

func (l *LogProcess) Process()  {
	//解析模块
	/*
	2019-10-01 22:56:35'715" 18920	[DEBUG]	Listen Server Send the request failed, 12029; File: .\ServerWatcher.cpp, Line: 266, Function: CServerWatcher::ListenServer
	172.0.0.12 -- [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 “-” “keepAliveClient” "-" 1.005 1.854
	*/
	regexp.MustCompile("")
	for	data := range l.rc{
		//l.wc <- strings.ToUpper(string(data))
		l.wc <- string(data)
	}
}

func (w *Writer) WirteToInfluxDB(wc chan string) {
	//写入 InfluxDB 时序数据库
	for data := range wc{
		println(data)
	}
}
