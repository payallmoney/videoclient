package main

import (
	"github.com/go-martini/martini"
	"runtime"
	"os/exec"
	"strings"
	"log"
)

func main() {
	m := martini.Classic()
	m.Get("/",index)
	m.Get("/cpuid", cpuid)
	m.Run()
}

func index() string {
	return "Hello world!"
}
func cpuid() string {
	cpuid:=""
	if runtime.GOOS == "windows" {
		// "wmic cpu get ProcessorId /format:csv"
		strs, err := exec.Command("wmic","cpu","get","ProcessorId","/format:csv").Output()
		if err != nil {
			log.Fatal(err)
		}
		id :=strings.Split(strings.Split(string(strs),"\r\n")[2],",")[1];
		cpuid = strings.Trim(id," ")
	}else {
		// "cat /proc/cpuinfo | grep Serial | cut -d ':' -f 2'"
		strs , _ := exec.Command("cat","/proc/cpuinfo","|","grep","Serial","|","cut","-d","':'","-f","2'").Output();
		cpuid = strings.Trim(string(strs)," ")
	}
	return cpuid
}