package main

import (
	"encoding/json"
	"os/exec"
	"os"
	"net"
	"time"
	"fmt"
)

const CONFIG_PATH = "servdown.json"

type ServData struct {
	Host			string
	Port			uint
	Protocol		string//tcp/udp
	TimeoutSec		uint
	IntervalSec		uint
	Successes		uint
	Failures		uint
	UptimeRatio		string
}

func pingServer(host string, timeout uint) bool {
	err := exec.Command("ping", host, "-c 1", fmt.Sprintf("-W %v", timeout)).Run()
	return err != nil;
}

func testConn(host string, port uint, proto string, timeout time.Duration) bool {
	conn, err := net.DialTimeout(proto, fmt.Sprintf("%v:%v", host, port), timeout)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

func main(){

	servData := ServData{}
	file, err := os.OpenFile(CONFIG_PATH, os.O_RDWR, os.ModeAppend)
	if os.IsNotExist(err) {
		servData = ServData{
			Host:			"psilocyber.tech",
			Port:			80,
			Protocol:		"tcp",
			TimeoutSec:		3,
			IntervalSec:	60,
		}
		fmt.Printf("ping data file '%s' does not exist, creating default one\n", CONFIG_PATH)
		file, _ = os.Create(CONFIG_PATH)
		bytes, _ := json.MarshalIndent(servData, "", "\t")
		file.Write(bytes)
	}

	defer file.Close()
	file.Seek(0, 0)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&servData)
	if err != nil {
		fmt.Printf("failed to parse ping data file '%s', err: %v\n", CONFIG_PATH, err.Error())
		return
	}

	for {
		if testConn(servData.Host, servData.Port, servData.Protocol, time.Duration(servData.TimeoutSec) * time.Second) {
			servData.Successes++;
		} else {
			servData.Failures++;
		}

		if servData.Failures == 0 {
			servData.UptimeRatio = "100%"
		} else {
			servData.UptimeRatio = fmt.Sprintf("%v%%", (float32(servData.Successes) / (float32(servData.Successes) + float32(servData.Failures))) * float32(100))
		}
		bytes, err := json.MarshalIndent(servData, "", "\t")
		
		file.Truncate(0)
		file.Seek(0, 0)
		_, err = file.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write, err: %v\n", err.Error())
			return
		}
		file.Sync()
		time.Sleep(time.Duration(servData.IntervalSec) * time.Second)
	}

}