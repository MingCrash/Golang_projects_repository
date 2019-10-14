package main

import (
	"fmt"
	"strings"
)

func main()  {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	//rd := Reader{
	//	logpath:"MallOrder-log.txt",
	//}
	//wd := Writer{
	//	influxDBDsn:"name&pass",
	//}
	//
	//lp := &LogProcess{
	//	rc: make(chan []byte),
	//	wc: make(chan LogMessage),
	//	logreader: rd,
	//	infwriter: wd,
	//}
	//go lp.logreader.ReadFromFile(lp.rc)
	//go lp.Process()
	//go lp.infwriter.WirteToInfluxDB(lp.wc)
	//
	//time.Sleep(30*time.Second)

	//loc,_ := time.LoadLocation("PRC")
	//fmt.Println(time.Now().Format("2006/01&02 15:04:05"))
	//fmt.Println(time.Now().Format(time.RFC822))
	////fmt.Println(time.Parse("2016-01-02 15:04:05", "2018-04-23 12:24:51"))
	//fmt.Println(time.ParseInLocation("2006-01-02 15:04:05", "2012-05-11 14:06:06", loc))
	//fmt.Println(time.Now().Date())
	//fmt.Println(time.Now().Clock())
	//fmt.Println(time.Now().Second())
	//fmt.Println()

	css := "http://localhost:8086&zhuzhiming&suzuki"
	fmt.Println(strings.Split(css,"&"))
}