package main

import (
	"os"
"runtime"
"os/exec"
"strings"
"path/filepath"
	"github.com/skip2/go-qrcode"
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
		// "cat /proc/cpuinfo | grep Serial | cut -d ':' -f 2'"
		strs, _ := exec.Command("cat", "/proc/cpuinfo", "|", "grep", "Serial", "|", "cut", "-d", "':'", "-f", "2'").Output();
		cpuid = strings.TrimSpace(string(strs))
	}
	return cpuid
}

func cpuqr() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	picpath := dir + "/cpuqr.png"
	err := qrcode.WriteFile("raspi:" + cpuid(), qrcode.Highest, 256, picpath)
	checkerr(err)
	return picpath
}

