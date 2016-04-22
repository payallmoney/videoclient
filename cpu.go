package main

import (
//"runtime"
//"os/exec"
//"strings"
	"github.com/shirou/gopsutil/cpu"
	"log"
)

/*
	https://github.com/MetaScale/nt2/blob/a1ede0dc54a332dcc265e1cfed267d4b4a57e019/tools/is_multicore/is_multicore.cpp
*/
//
//func cpuid() string {
//	cpuid := ""
//	if runtime.GOOS == "windows" {
//		// "wmic cpu get ProcessorId /format:csv"
//		cmd :=exec.Command("cmd", "/c" ,"wmic cpu get ProcessorId /format:csv")
//
//		strs, err := cmd.Output()
//		checkerr(err)
//
//		id := strings.Split(strings.Split(string(strs), "\r\n")[2], ",")[1];
//		cpuid = strings.TrimSpace(id)
//
//		checkerr(err)
//	}else {
//		// "cat /proc/cpuinfo | grep Serial | cut -d ':' -f 2"
//		strs,_ :=exec.Command("bash","-c","cat /proc/cpuinfo | grep Serial | cut -d ':' -f 2").Output()
//		cpuid = strings.TrimSpace(string(strs))
//	}
//	return cpuid
//}

func cpuid()string{
	ids,_ := cpu.Info()
	return ids[0].PhysicalID
}

func main(){
	log.Println(cpuid())
}
