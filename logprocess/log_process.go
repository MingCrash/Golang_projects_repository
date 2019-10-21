package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Reader struct {
	logpath 	string
}

type Writer struct {
	influxDBDsn string
}

type LogProcess struct {
	rc chan 	[]byte
	wc chan 	LogMessage
	logreader 	Reader
	infwriter 	Writer
}

type Read interface {
	ReadFromFile(chan []byte)
}

type Write interface {
	WirteToInfluxDB(chan LogProcess)
}

//日志信息
type LogMessage struct {
	Time 								time.Time
	Level, TraceId, Url, statCode 		string
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
		if err == io.EOF || err != nil{
			statusMonitorChan <- ErrorTypeNum
			time.Sleep(500*time.Millisecond)
			continue
		}
		//if  {
		//	statusMonitorChan <- ErrorTypeNum
		//	panic(fmt.Sprintf("log file readed failture:%s", err.Error()))
		//	continue
		//}
		rc <- line[:len(line)-1]
	}
}

func (l *LogProcess) Process()  {
	//解析模块
	loc,err := time.LoadLocation("PRC")
	rg, err := regexp.Compile("t=([\\d+\\:\\-\\s]+)[^\n]*Level=([a-z]+)[^\n]*TraceId=(\\d+)[^\n]*Url=([^`]+)[^\n]*code=(\\d+)")
	for	data := range l.rc{
		if err != nil {
			statusMonitorChan <- ErrorTypeNum
			continue
		}
		list := rg.FindStringSubmatch(string(data))
		if len(list) < 6 {
			statusMonitorChan <- ErrorTypeNum
			continue
		}
		msg := LogMessage{}
		msg.Time,_ = time.ParseInLocation("2006-01-02 15:04:05", list[1], loc)
		msg.Level = list[2]
		msg.TraceId = list[3]
		msg.Url = list[4]
		msg.statCode = list[5]
		l.wc <- msg
		statusMonitorChan <- HandlelineTypeNum
	}
}

func (w *Writer) WirteToInfluxDB(wc chan LogMessage) {
	//写入 InfluxDB 时序数据库

	// Create a new HTTPClient
	sn := strings.Split(w.influxDBDsn,"@")
	cle, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     	sn[0],
		Username: 	sn[1],
		Password: 	sn[2],
	})
	if err != nil {
		statusMonitorChan <- ErrorTypeNum
		panic("Error creating InfluxDB Client: ")
	}
	defer cle.Close()

	// Close client resources
	//if err := cle.Close(); err != nil {
	//	fmt.Println("Error closing client: ", err.Error())
	//}

	for data := range wc{
		// Create a point
		tags := map[string]string{
			"Url": 		data.Url,
			"Level":   	data.Level,
			"Code":   	data.statCode,
		}
		fields := map[string]interface{}{
			"Time":   	data.Time,
			"TraceId": 	data.TraceId,
		}
		fmt.Println(tags, fields)

		pt, err := client.NewPoint("Mallorder", tags, fields, time.Now())
		if err != nil {
			statusMonitorChan <- ErrorTypeNum
			log.Fatal(err)
		}

		// Create a new point batch  创建批量注入
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  "logprocess",
			Precision: "s",
		})
		if err != nil {
			statusMonitorChan <- ErrorTypeNum
			fmt.Println("Error creating NewBatchPoints: ", err.Error())
		}

		//add point to point batch
		bp.AddPoint(pt)

		// Write the batch
		if err := cle.Write(bp); err != nil {
			statusMonitorChan <- ErrorTypeNum
			log.Fatal(err)
		}
	}
}


//=====================================================================================================================

const (
	ErrorTypeNum = 0
	HandlelineTypeNum = 1
)

var statusMonitorChan = make(chan int ,100)

//日志系统监控信息
type Systeminfo struct {
	Handleline  	int           `json:"handleline"`		//总处理日志行数
	Tps         	float64       `json:"tps"`				//监控系统吞吐量
	Readchan    	int           `json:"readchan"` 		//当前读日志缓存数(read channel长度)
	Writechan    	int           `json:"writechan"`		//当前读日志缓存数(write channel长度)
	RunTotletime	string        `json:"run_totletime"`	//日志系统运行总时间
	ErrTotal 		int 		  `json:"err_total"`		//处理错误总数
}

type SystemMonitor struct {
	StartTime  		time.Time
	info			Systeminfo
}



func (sm *SystemMonitor) Monitor(lp *LogProcess)  {
	//日志系统监控器模块
	var subline int

	//日志系统监控信息收集协程
	go func() {
		for sy := range statusMonitorChan {
			switch sy {
			case ErrorTypeNum:
				sm.info.ErrTotal ++
			case HandlelineTypeNum:
				sm.info.Handleline ++
			}
		}
	}()

	//每五秒
	ticker := time.NewTicker(5*time.Second)
	go func() {
		for {
			lasttmp := sm.info.Handleline
			//此处在等待channel中的信号，因此执行此段代码时会阻塞5秒
			_,ok := <-ticker.C   	//读取空 channel 会阻塞
			if !ok {
				log.Fatal("Error:channel has been closed")
			}
			subline = sm.info.Handleline - lasttmp
		}
	}()

	//创建http server
	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		sm.info.RunTotletime = string(time.Now().Sub(sm.StartTime))
		sm.info.Readchan = len(lp.rc)
		sm.info.Writechan = len(lp.wc)
		sm.info.Tps = float64(subline/5)

		ret,err := json.MarshalIndent(sm.info,"","\t")
		if err != nil {
			log.Fatal("")
		}
		_, _ = io.WriteString(writer, string(ret))
	})

	_ = http.ListenAndServe(":9091", nil)  //ListenAndServe 有阻塞作用
}