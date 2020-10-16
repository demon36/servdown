package main

import (
	"encoding/json"
	"os/exec"
	"os"
	"time"
	"fmt"
)

const CONFIG_PATH = "data.json"

type PingData struct {
	Server			string
	TimeoutSec		uint
	IntervalSec		uint
	Successes		uint
	Failures		uint
	UptimeRatio		string
}

func main(){

	pingData := PingData{}
	file, err := os.OpenFile(CONFIG_PATH, os.O_RDWR, os.ModeAppend)
	if os.IsNotExist(err) {
		pingData = PingData{
			Server:			"psilocyber.tech",
			TimeoutSec:		3,
			IntervalSec:	60,
		}
		fmt.Printf("ping data file '%s' does not exist, creating default one\n", CONFIG_PATH)
		file, _ = os.Create(CONFIG_PATH)
		bytes, _ := json.MarshalIndent(pingData, "", "\t")
		file.Write(bytes)
	}

	defer file.Close()
	file.Seek(0, 0)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&pingData)
	if err != nil {
		fmt.Printf("failed to parse ping data file '%s', err: %v\n", CONFIG_PATH, err.Error())
		return
	}

	for {
		err := exec.Command("ping", pingData.Server, "-c 1", fmt.Sprintf("-W %v", pingData.TimeoutSec)).Run()
		if err == nil {
			pingData.Successes++;
		} else {
			pingData.Failures++;
		}

		if pingData.Failures == 0 {
			pingData.UptimeRatio = "100%"
		} else {
			pingData.UptimeRatio = fmt.Sprintf("%v%%", (float32(pingData.Successes) / (float32(pingData.Successes) + float32(pingData.Failures))) * float32(100))
		}
		bytes, err := json.MarshalIndent(pingData, "", "\t")
		
		file.Truncate(0)
		file.Seek(0, 0)
		_, err = file.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write, err: %v\n", err.Error())
			return
		}
		file.Sync()
		time.Sleep(time.Duration(pingData.IntervalSec) * time.Second)
	}

}