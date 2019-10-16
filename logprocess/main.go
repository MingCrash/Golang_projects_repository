package main

import (
	"flag"
	"runtime"
)

func main()  {
	var logpath, influxsn string
	flag.StringVar(&logpath, "logpath", "MallOrder-log.log","set log path of target for reading")
	flag.StringVar(&influxsn, "influxsn", "http://localhost:8086@zhuzhiming@suzuki","set influDB's host and client in format:http://localhost:8086@zhuzhiming@suzuki")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	rd := Reader{
		logpath:"MallOrder-log.log",
	}
	wd := Writer{
		influxDBDsn:"http://192.168.0.48:8086@zhuzhiming@suzuki",
	}

	lp := &LogProcess{
		rc: make(chan []byte,100),   		//添加缓存，缓冲不同协程的处理速度差问题
		wc: make(chan LogMessage,100),  	//添加缓存，缓冲不同协程的处理速度差问题
		logreader: rd,
		infwriter: wd,
	}
	sy := &SystemMonitor{}
	go lp.logreader.ReadFromFile(lp.rc)
	go lp.Process()
	go lp.infwriter.WirteToInfluxDB(lp.wc)

	sy.Monitor(lp)

	//ticker := time.NewTicker(5*time.Second)
	//go func() {
	//	for {
	//		//此处在等待channel中的信号，因此执行此段代码时会阻塞5秒
	//		fmt.Println("未接受信号,等待")
	//		<-ticker.C
	//		fmt.Println("已接受信号，处理定时函数")
	//	}
	//}()
	//
	//time.Sleep(11*time.Second)
}