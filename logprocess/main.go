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

	time.Sleep(60*time.Second)

}