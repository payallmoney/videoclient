package main

import (
	"github.com/go-martini/martini"
	"runtime"
	"os/exec"
	"strings"
	"log"
	qrcode "github.com/skip2/go-qrcode"
	"path/filepath"
	"os"
	"net/http"
	"io/ioutil"
	"fmt"
)

func main() {
	m := martini.Classic()
	m.Get("/reg", reg)
	m.Get("/active", active)
	m.Get("/qr", cpuqr)
	m.Run()
}

func reg() string {
	fmt.Println(HttpUrl("/video/reg/" + cpuid()));
	resp, err := http.Get(HttpUrl("/video/reg/" + cpuid()))
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
	return string(body);
}

func active() string {
	resp, err := http.Get(HttpUrl("/video/active/" + cpuid()))
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body);
}

func cpuid() string {
	cpuid := ""
	if runtime.GOOS == "windows" {
		// "wmic cpu get ProcessorId /format:csv"
		strs, err := exec.Command("wmic", "cpu", "get", "ProcessorId", "/format:csv").Output()
		if err != nil {
			log.Fatal(err)
		}
		id := strings.Split(strings.Split(string(strs), "\r\n")[2], ",")[1];
		cpuid = strings.Trim(id, " ")
	}else {
		// "cat /proc/cpuinfo | grep Serial | cut -d ':' -f 2'"
		strs, _ := exec.Command("cat", "/proc/cpuinfo", "|", "grep", "Serial", "|", "cut", "-d", "':'", "-f", "2'").Output();
		cpuid = strings.Trim(string(strs), " ")
	}
	return cpuid
}

func cpuqr() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	picpath := dir + "/cpuqr.png"
	err := qrcode.WriteFile("raspi:" + cpuid(), qrcode.Highest, 256, picpath)
	if err != nil {
		log.Fatal(err)
	}
	return picpath
}

