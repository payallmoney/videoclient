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
	"time"
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
	id :=cpuid()
	resp, err := http.Get(HttpUrl("/client/reg/" + id))
	checkerr(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body);
}

func active() string {
	id :=cpuid()
	resp, err := http.Get(HttpUrl("/client/active/" + id))
	checkerr(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body);
}

func check(r render.Render) {
	checknetwork(circle)
	//反过来调用播放接口
	var res interface{}
	res = "true"

	r.JSON(200, res)
}



type next func()

func checknetwork( function next){
	resp, err := http.Get(HttpUrl("/test"))
	defer resp.Body.Close()
	checktime := time.Second*10
	ticker := time.NewTimer(checktime)
	var flag bool
	if(err!=nil) {

		flag = false
	}else{
		body, _ := ioutil.ReadAll(resp.Body)
		if(string(body)=="连接测试"){
			flag = true
			log_print("连接测试成功,开始执行!")
			function()
		}else{
			flag = false
		}
	}
	//如果未成功 ,继续测试
	if(!flag){
		log_print("连接测试失败,10秒后重试!")
		go func() {
			for _  = range ticker.C {
				checknetwork( function)
			}
		}()
	}
}

func cpuqr() string {
	id :=cpuid()
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	picpath := dir + "/cpuqr.png"
	err := qrcode.WriteFile("raspi:" + id, qrcode.Highest, 256, picpath)
	checkerr(err)
	return picpath
}



func status(r render.Render){
	id :=cpuid()
	resp, err := http.Get(HttpUrl("/client/status/" + id))
	checkerr(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result  interface{}
	json.Unmarshal(body, &result)
	r.JSON(200, result)
}

func getStatus(){
	id :=cpuid()
	resp, err := http.Get(HttpUrl("/client/status/" + id))
	checkerr(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result  interface{}
	json.Unmarshal(body, &result)
}




