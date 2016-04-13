package main

import (
	"github.com/go-martini/martini"
	"log"
	"path/filepath"
	"os"
	"net/http"
	"io/ioutil"
	"github.com/martini-contrib/render"
	"github.com/skip2/go-qrcode"
	"encoding/json"
)

var checking = false
var rootpath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

func main() {
	//设置日志

	log.SetFlags(log.LstdFlags | log.Llongfile)

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/reg", reg)
	m.Get("/active", active)
	m.Get("/check", check)
	m.Get("/status", status)
	m.Get("/qr", cpuqr)
	//m.Run()
	m.RunOnAddr(":10001")
}

func reg() string {
	resp, err := http.Get(HttpUrl("/client/reg/" + cpuid()))
	checkerr(err)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body);
}

func active() string {
	resp, err := http.Get(HttpUrl("/client/active/" + cpuid()))
	checkerr(err)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body);
}

func check(r render.Render) {
	videocheck()
	circle()
	//反过来调用播放接口
	var res interface{}
	res = "true"

	r.JSON(200, res)
}

func cpuqr() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	picpath := dir + "/cpuqr.png"
	err := qrcode.WriteFile("raspi:" + cpuid(), qrcode.Highest, 256, picpath)
	checkerr(err)
	return picpath
}



func status(r render.Render){
	resp, err := http.Get(HttpUrl("/client/status/" + cpuid()))
	checkerr(err)
	body, _ := ioutil.ReadAll(resp.Body)
	var result  interface{}
	json.Unmarshal(body, &result)
	r.JSON(200, result)
}




