package main

import (
	"flag"
	"os"
	"zmqclient/zmqclient"

	"github.com/golang/glog"
)

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set("stderrthreshold", "Info")
	flag.Set("logtostderr", "true")
	flag.Parse()
}

func main() {
	glog.Infoln("Start dfxp Client")
	option := zmqclient.ZmqClientOptions{}
	client := zmqclient.NewZmqClient(&option)
	if client == nil {
		glog.Fatalf("Create zmqClient failed.")
		return
	}

	err := zmqclient.ZmqClientExe()
	if err != nil {
		glog.Errorf("zmq client failed.Err:%s", err)
	}

	glog.Info("zmq client exit")

}
