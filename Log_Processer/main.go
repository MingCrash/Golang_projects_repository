package main

import (
	"runtime"
	"time"
)

func main()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rd := Reader{
		logpath:"MallOrder-log.log",
	}
	wd := Writer{
		influxDBDsn:"name&pass",
	}

	lp := &LogProcess{
		rc: make(chan []byte),
		wc: make(chan LogMessage),
		logreader: rd,
		infwriter: wd,
	}
	go lp.logreader.ReadFromFile(lp.rc)
	go lp.Process()
	go lp.infwriter.WirteToInfluxDB(lp.wc)

	time.Sleep(10*time.Second)

	//str := "t=2019-10-07 00:05:09`Level=error`TraceId=157037790590870`Url=/api/Asyn/run`UserAgent=curl/7.19.7 (x86_64-redhat-linux-gnu) libcurl/7.19.7 NSS/3.27.1 zlib/1.2.3 libidn/1.18 libssh2/1.4.2`Message=MallOrder"
	//rg,_ := regexp.Compile("t=([\\d+\\:\\-\\s]+)[^\n]*Level=([a-z]+)[^\n]*TraceId=(\\d+)[^\n]*Url=([^`]+)")
	//for i,v:= range rg.FindStringSubmatch(str) {
	//	println(i,v)
	//}
}