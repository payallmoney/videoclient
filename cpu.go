package main

import (
"runtime"
"os/exec"
"strings"
)

func cpuid() string {
	cpuid := ""
	if runtime.GOOS == "windows" {
		// "wmic cpu get ProcessorId /format:csv"
		strs, err := exec.Command("wmic", "cpu", "get", "ProcessorId", "/format:csv").Output()
		checkerr(err)
		id := strings.Split(strings.Split(string(strs), "\r\n")[2], ",")[1];
		cpuid = strings.TrimSpace(id)
	}else {
		// "cat /proc/cpuinfo | grep Serial | cut -d ':' -f 2"
		strs,_ :=exec.Command("bash","-c","cat /proc/cpuinfo | grep Serial | cut -d ':' -f 2").Output()
		cpuid = strings.TrimSpace(string(strs))
	}
	return cpuid
}



