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
	wc chan LogMessage
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
	Time, Level, TraceId, Url string
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
		println(line[:len(line)-1])
		rc <- line[:len(line)-1]
	}
}

func (l *LogProcess) Process()  {
	//解析模块
	rg, err := regexp.Compile("t=([\\d+\\:\\-\\s]+)[^\n]*Level=([a-z]+)[^\n]*TraceId=(\\d+)[^\n]*Url=([^`]+)")
	for	data := range l.rc{
		if err != nil || len(list) < 5{
			continue
		}
		list := rg.FindStringSubmatch(string(data))
		loms := LogMessage {
			Time: 		list[1],
			Level:		list[2],
			TraceId: 	list[3],
			Url: 		list[4],
		}
		l.wc <- loms
	}
}

func (w *Writer) WirteToInfluxDB(wc chan LogMessage) {
	//写入 InfluxDB 时序数据库
	for data := range wc{
		println(data)
	}
}
