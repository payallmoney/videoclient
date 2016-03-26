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
	"encoding/json"
	"reflect"
	"io"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2/bson"
	"bytes"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/reg", reg)
	m.Get("/active", active)
	m.Get("/check", check)
	m.Get("/qr", cpuqr)
	m.Run()
}

func reg() string {
	log.Println(HttpUrl("/video/reg/" + cpuid()));
	resp, err := http.Get(HttpUrl("/video/reg/" + cpuid()))
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(body)
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
func circle(){
	checktime := Cfg()["checktime"].(time.Duration)
	ticker := time.NewTicker(checktime)
	go func() {
		for  range ticker.C {
			videocheck()
		}
	}()
}

func videocheck() {
	log.Println("====videocheck===")
	resp, err := http.Get(HttpUrl("/video/list/" + cpuid()))
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var result  interface{}
	json.Unmarshal(body, &result)
	lists := reflect.ValueOf(result)
	files := make([]string, lists.Len())
	for i := 0; i < lists.Len(); i++ {
		path := lists.Index(i).Elem().String()
		files[i] = filename(path)
	}
	//反过来调用播放接口
	if(isSamelist(files)){
		//相同的条件下不做任何事情
	}else{
		log.Println("**********not same*************")
		//2,有差异则替换
		for i := 0; i < lists.Len(); i++ {
			path := lists.Index(i).Elem().String()
			donwfile(path)
		}
		playvideos(files,false)
	}
}

func check(r render.Render) {
	videocheck()
	circle()
	//反过来调用播放接口
	var res interface{}
	res ="true"

	r.JSON(200, res)
}
func isSamelist(newlist []string) bool{
	list := getlist().(map[string]interface{})["result"].(map[string]interface{})["items"].([]interface{})
	if(len(newlist) !=len(list)){
		return false;
	}else{
		for i:=0;i<len(list);i++{
			newlabel := filepath.Base(newlist[i])
			label := list[i].(map[string]interface{})["label"].(string)
			if newlabel != label{
				return false;
			}
		}
		return true;
	}
}

func getPlayer(id int) {
	params := bson.M{"jsonrpc": "2.0", "method": "Playlist.Clear", "params": bson.M{"playlistid": id}, "id": 1}
	jsonrpc(params, false)
}

func playvideos(videolist []string,flag bool) interface{} {
	//清除播放列表
	clearPlayList()
	//增加列表
	for i := 0; i < len(videolist); i++ {
		additem(videolist[i])
	}
	//未开始强制播放的情况下进行平滑播放,播放完当前视频才切换新视频列表
	if(flag){
		//强制开始播放
		play()
	}else if(!activePlayer()){
		//如果未开始播放,则进行播放
		play()
	}
	return getlist()
}

func clearPlayList()  interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Playlist.Clear", "params": bson.M{"playlistid": 1}, "id": 1}
	return jsonrpc(params, false)
}


func getPlayeditem() string{
	params := bson.M{"jsonrpc": "2.0", "method": "Player.GetItem", "params": bson.M{"playerid": 1}, "id": 1}
	ret := jsonrpc(params, false);
	result ,_ := ret.( map[string]interface{})["result"].( map[string]interface{})["item"].( map[string]interface{})["label"].(string)
	return result
}
func activePlayer() bool{
	params := bson.M{"jsonrpc": "2.0", "method": "Player.GetActivePlayers", "id": 1}
	ret :=jsonrpc(params, false).(map[string]interface{})["result"].([]interface{});
	if(len(ret)>0){
		return true
	}else{
		return false
	}
}
func playPause() interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.PlayPause", "id": 1, "params":  bson.M{
		"playerid": 1 }}
	return jsonrpc(params, false)
}
func play()  interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.Open", "id": 1, "params":  bson.M{
		"item":bson.M{"playlistid":1},
		"options":bson.M{"repeat":"all"}}}
	return jsonrpc(params, false)
}
func playVideo(video string)  interface{} {
	video = strings.Replace(video, "\\", "/", -1)
	params := bson.M{"jsonrpc": "2.0", "method": "Player.Open", "id": 1, "params":  bson.M{
		"item":bson.M{"file":video} }}
	return jsonrpc(params, false)
}
func additem(file string) interface{} {
	params := bson.M{
		"jsonrpc": "2.0",
		"method": "Playlist.Add",
		"id": 1,
		"params": bson.M{
			"playlistid": 1,
			"item":bson.M{"file":file}}}
	return jsonrpc(params, false)
}

func jsonrpc(param interface{}, flag bool) interface{} {
	params_str, err := json.Marshal(param)
	if (flag) {
		log.Println(string(params_str))
	}
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(KodiUrl("/jsonrpc"), "application/json", bytes.NewBuffer([]byte(params_str)))
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if (flag) {
		log.Printf("=====ret====%s\r\n", body)
	}
	var result  interface{}
	json.Unmarshal(body, &result)
	return result
}

func getlist() interface{} {
	params := bson.M{
		"jsonrpc": "2.0",
		"method": "Playlist.GetItems",
		"params": bson.M{"properties": []string{"runtime", "showtitle", "season", "title", "artist" },
			"playlistid": 1},
		"id": 1}

	return jsonrpc(params, false)
}
func js(item interface{}) string {
	params_str, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}
	return (string(params_str))
}

func filename(path string) string {
	idx := strings.LastIndex(path, "/")
	rootpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	ret := rootpath + "/video" + path[idx:];
	ret, err = filepath.Abs(ret);
	if err != nil {
		log.Fatal(err)
	}
	return ret
}
func donwfile(path string) string {
	realpath := filename(path)
	//检查文件是否存在,如果已存在则不再下载
	if fileexists(realpath){
		return realpath
	}
	downloadurl := HttpUrl(path)

	out, err := os.Create(realpath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	resp, err := http.Get(downloadurl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	return realpath
}

func fileexists(file string) bool{
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false;
	}else{
		return true;
	}
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
	if err != nil {
		log.Fatal(err)
	}
	return picpath
}

