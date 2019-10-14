package main

import (
	"bufio"
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"image"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
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
	Time 						time.Time
	Level, TraceId, Url 		string
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
	loc,err := time.LoadLocation("PRC")
	rg, err := regexp.Compile("t=([\\d+\\:\\-\\s]+)[^\n]*Level=([a-z]+)[^\n]*TraceId=(\\d+)[^\n]*Url=([^`]+)")
	for	data := range l.rc{
		if err != nil {
			continue
		}
		list := rg.FindStringSubmatch(string(data))
		if len(list) < 5 {
			continue
		}
		msg := LogMessage{}
		//msg.Timelocl = time.Now()
		msg.Time,_ = time.ParseInLocation("2006-01-02 15:04:05", list[1], loc)
		msg.Level = list[2]
		msg.TraceId = list[3]
		msg.Url = list[4]

		l.wc <- msg
	}
}

func (w *Writer) WirteToInfluxDB(wc chan LogMessage) {
	//写入 InfluxDB 时序数据库

	// Create a new HTTPClient
	//http://localhost:8086&zhuzhiming&suzuki
	sn := strings.Split(w.influxDBDsn,"&")
	cle, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     	sn[0],
		Username: 	sn[1],
		Password: 	sn[2],
	})
	if err != nil {
		panic("Error creating InfluxDB Client: ")
	}
	defer cle.Close()

	// Close client resources
	if err := cle.Close(); err != nil {
		fmt.Println("Error closing client: ", err.Error())
	}

	for data := range wc{
		//println(data.Url,data.TraceId,data.Level,data.Time.Format("2006-01-02 15:04:05"),data.Timelocl.Format("2006-01-02 15:04:05"))
		// Create a point
		tags := map[string]string{
			"Url": data.Url,
			"Level":   data.Level,
			"TraceId": data.TraceId,
		}
		fields := map[string]interface{}{
			"Time":   data.Time,

		}
		pt, err := client.NewPoint("mallorder", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		// Create a new point batch  创建批量注入
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  MyDB,
			Precision: "s",
		})
		if err != nil {
			fmt.Println("Error creating NewBatchPoints: ", err.Error())
		}

		//add point to point batch
		bp.AddPoint(pt)

		// Write the batch
		if err := cle.Write(bp); err != nil {
			log.Fatal(err)
		}
	}
}
