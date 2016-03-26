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
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/test", test)
	m.Get("/reg", reg)
	m.Get("/active", active)
	m.Get("/getvideolist", getvideolist)
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

func getvideolist(r render.Render) {
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
	var res interface{}
	//TODO  1,比较当前列表与列表的差异
	if(isSamelist(files)){
		log.Println("########same#########")
		//
	}else{
		log.Println("**********not same*************")
		//2,有差异则替换
		for i := 0; i < lists.Len(); i++ {
			path := lists.Index(i).Elem().String()
			donwfile(path)
		}
		res = playvideos(files)
	}
	r.JSON(200, res)
}
func isSamelist(newlist []string) bool{

	return false
}

func getPlayer(id int) {
	params := bson.M{"jsonrpc": "2.0", "method": "Playlist.Clear", "params": bson.M{"playlistid": id}, "id": 1}
	jsonrpc(params, false)
}

func playvideos(videolist []string) interface{} {
	playlistid := 1
	//清除播放列表
	//clearPlayList(playlistid)
	//TODO 删除非当前播放的列表
	clearOtherList(playlistid)
	//TODO 增加记录
	//先调用第一条播放记录进行播放
	for i := 0; i < len(videolist); i++ {
		additem(videolist[i], playlistid)
	}
	//TODO 判断是否播放

	//TODO
	play(playlistid)
	//activePlayer()
	return getlist(playlistid)
}
func clearOtherList(id int){
	currlist :=getlist(id)
	log.Printf("%v",currlist)
}
func clearPlayList(id int) interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Playlist.Clear", "params": bson.M{"playlistid": id}, "id": 1}
	return jsonrpc(params, false)
}
func activePlayer() interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.GetActivePlayers", "id": 1}
	return jsonrpc(params, false)
}
func playPause(id int) interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.PlayPause", "id": 1, "params":  bson.M{
		"playerid": id }}
	return jsonrpc(params, false)
}
func play(id int) interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.Open", "id": 1, "params":  bson.M{
		"item":bson.M{"playlistid":id} }}
	return jsonrpc(params, false)
}
func playVideo(video string) interface{} {
	video = strings.Replace(video, "\\", "/", -1)
	params := bson.M{"jsonrpc": "2.0", "method": "Player.Open", "id": 1, "params":  bson.M{
		"item":bson.M{"file":video} }}
	return jsonrpc(params, false)
}
func additem(file string, id int) interface{} {
	params := bson.M{
		"jsonrpc": "2.0",
		"method": "Playlist.Add",
		"id": 1,
		"params": bson.M{
			"playlistid": id,
			"item":bson.M{"file":file}}}
	return jsonrpc(params, false)
}
//func jsonrpc (param interface{} ) interface{} {
//	return kodirpc(param,false)
//}
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

func getlist(id int) interface{} {
	params := bson.M{
		"jsonrpc": "2.0",
		"method": "Playlist.GetItems",
		"params": bson.M{"properties": []string{"runtime", "showtitle", "season", "title", "artist" },
			"playlistid": id},
		"id": 1}

	return jsonrpc(params, false)
}
func print(item interface{}) {
	params_str, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(params_str))
}
func test() string {
	//client, err := jsonrpc.Dial("http", "localhost:8080/")
	var params bson.M
	//params = bson.M{
	//	"jsonrpc": "2.0", "method": "JSONRPC.Introspect", "params": bson.M{ "filter":  bson.M{ "id": "Player.GetActivePlayers", "type": "method" } },
	//	"id": 1 }
	//jsonrpc(params)
	//params =bson.M{
	//	"jsonrpc": "2.0",
	//	"method": "Player.GetActivePlayers",
	//	"id": 1}
	//jsonrpc(params)
	//
	//params =bson.M{
	//	"jsonrpc": "2.0",
	//	"method": "Playlist.GetItems",
	//	"params": bson.M{ "properties": []string{ "runtime", "showtitle", "season", "title", "artist" },
	//		"playlistid": 1},
	//	"id": 1}
	//
	//jsonrpc(params)
	//params =bson.M{
	//	"jsonrpc": "2.0",
	//	"method": "Playlist.Add",
	//	"params": bson.M{
	//		"item":bson.M{"file":"D:\\videowork\\videoclient/video/d4db2188-1f73-4a72-9ddc-f5868ba72714.mp4"},
	//		"playlistid": 1},
	//	"id": 1}
	//jsonrpc(params)
	//
	//params =bson.M{
	//	"jsonrpc": "2.0",
	//	"method": "Playlist.GetItems",
	//	"params": bson.M{ "properties": []string{ "runtime", "showtitle", "season", "title", "artist" },
	//		"playlistid": 1},
	//	"id": 1}
	//
	//jsonrpc(params)
	params = bson.M{"jsonrpc": "2.0", "method": "Player.Open", "params":  bson.M{
		"item":bson.M{"file":"D:\\videowork\\videoclient/video/d4db2188-1f73-4a72-9ddc-f5868ba72714.mp4"}}, "id": 1}
	jsonrpc(params, false)
	//
	//params =bson.M{
	//	"jsonrpc": "2.0",
	//	"method": "Playlist.GetItems",
	//	"params": bson.M{ "properties": []string{ "runtime", "showtitle", "season", "title", "artist" },
	//		"playlistid": 1},
	//	"id": 1}
	//
	//jsonrpc(params)
	//params =bson.M{"jsonrpc": "2.0", "method": "Player.PlayPause", "params": bson.M{ "playerid": 1 }, "id": 1}
	//jsonrpc(params)
	//params =bson.M{"jsonrpc": "2.0", "method": "Playlist.Clear","params": bson.M{"playlistid": 1}, "id": 1}
	//jsonrpc(params)
	return "调用成功"
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
	log.Println(realpath)
	downloadurl := HttpUrl(path)
	log.Println(downloadurl)

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

