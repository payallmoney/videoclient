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
	"encoding/json"
	"reflect"
	"io"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2/bson"

"bytes"
)

func main() {
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

func getvideolist(r render.Render) {
	fmt.Println(cpuid());
	resp, err := http.Get(HttpUrl("/video/list/" + cpuid()))
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body));
	var result  interface{}
	json.Unmarshal(body, &result)
	lists := reflect.ValueOf(result)
	fmt.Println(reflect.TypeOf(lists.Index(0).Elem().String()));
	fmt.Println(lists.Index(0).Elem());
	ret := make([]string, lists.Len())
	//items := make([]bson.M, lists.Len())
	for i := 0; i < lists.Len(); i++ {
		path := lists.Index(i).Elem().String()
		realpath := donwfile(path)
		ret[i] = realpath
		//items[i] = bson.M{"file":ret[0]}
	}
	print(ret);
	//反过来调用播放接口

	res := playvideos(ret)

	r.JSON(200, res)
}

func getPlayer(id int){
	params :=bson.M{"jsonrpc": "2.0", "method": "Playlist.Clear","params": bson.M{"playlistid": id}, "id": 1}
	jsonrpc(params)
}

func playvideos(videolist []string) interface{}{
	//清除播放列表
	getlist()
	activePlayer()

	playlistid := 1
	clearPlayList(playlistid)
	getlist()
	//先调用第一条播放记录进行播放

	//playVideo(videolist[0])
	//play(playlistid)
	for i := 0; i < len(videolist); i++ {
		additem(videolist[i],playlistid)
	}

	activePlayer()
	return  getlist()
}
func clearPlayList(id int) interface{}{
	params :=bson.M{"jsonrpc": "2.0", "method": "Playlist.Clear","params": bson.M{"playlistid": id}, "id": 1}
	return jsonrpc(params)
}
func activePlayer()interface{}{
	params :=bson.M{"jsonrpc": "2.0", "method": "Player.GetActivePlayers", "id": 1}
	return jsonrpc(params)
}
func playPause(id int)interface{}{
	params :=bson.M{"jsonrpc": "2.0", "method": "Player.PlayPause", "id": 1, "params":  bson.M{
			"playerid": id }}
	return jsonrpc(params)
}
func play(id int)interface{}{
	params :=bson.M{"jsonrpc": "2.0", "method": "Player.Open", "id": 1, "params":  bson.M{
		"item":bson.M{"playlistid":id} }}
	return jsonrpc(params)
}
func playVideo(video string)interface{}{
	video = strings.Replace(video,"\\","/",-1)
	params :=bson.M{"jsonrpc": "2.0", "method": "Player.Open", "id": 1, "params":  bson.M{
		"item":bson.M{"file":video} }}
	return jsonrpc(params)
}
func additem(file string,id int )interface{}{
	file = strings.Replace(file,"\\","/",-1)
	params :=bson.M{
		"jsonrpc": "2.0",
		"method": "Playlist.Add",
		"id": 1,
		"params": bson.M{
			"playlistid": id,
			"item":bson.M{"file":file}}}
	return jsonrpc(params)
}
func jsonrpc (param interface{}) interface{} {
	params_str,err := json.Marshal(param)
	fmt.Println(string(params_str))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(KodiUrl("/jsonrpc" ),"application/json", bytes.NewBuffer( []byte(params_str)))
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var result  interface{}
	json.Unmarshal(body, &result)
	fmt.Printf("=====ret====%s\r\n",body)
	return result
}
func getlist() interface{}{
	params :=bson.M{
		"jsonrpc": "2.0",
		"method": "Playlist.GetItems",
		"params": bson.M{ "properties": []string{ "runtime", "showtitle", "season", "title", "artist" },
			"playlistid": 1},
		"id": 1}

	return jsonrpc(params)
}
func print(item interface{}){
	params_str,err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(params_str))
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
	params =bson.M{"jsonrpc": "2.0", "method": "Player.Open", "params":  bson.M{
		"item":bson.M{"file":"D:\\videowork\\videoclient/video/d4db2188-1f73-4a72-9ddc-f5868ba72714.mp4"}}, "id": 1}
	jsonrpc(params)
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
	return rootpath + "/video" + path[idx:]
}
func donwfile(path string) string {
	realpath := filename(path)
	fmt.Println(realpath)
	downloadurl := HttpUrl(path)
	fmt.Println(downloadurl)

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

