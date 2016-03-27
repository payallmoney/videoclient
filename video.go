package main

import (
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"reflect"
	"gopkg.in/mgo.v2/bson"
	"bytes"
	"strings"
	"path/filepath"
)

func circle() {
	checktime := Cfg()["checktime"].(time.Duration)
	ticker := time.NewTicker(checktime)
	go func() {
		for _ ,_ = range ticker.C {
			videocheck()
		}
	}()
}

func videocheck() {
	if checking {
		return
	}
	checking = true

	resp, err := http.Get(HttpUrl("/video/list/" + cpuid()))
	checkerr(err)
	body, _ := ioutil.ReadAll(resp.Body)
	var result  interface{}
	json.Unmarshal(body, &result)
	lists := reflect.ValueOf(result)
	if(lists.IsNil()){
		checking = false
		//设备尚未注册
		return
	}
	files := make([]string, lists.Len())
	for i := 0; i < lists.Len(); i++ {
		path := lists.Index(i).Elem().String()
		files[i] = filename(path)
	}
	//反过来调用播放接口
	if (isSamelist(files)) {
		//相同的条件下不做任何事情
	}else {
		//2,有差异则替换
		for i := 0; i < lists.Len(); i++ {
			path := lists.Index(i).Elem().String()
			donwfile(path)
		}
		playvideos(files, false)
	}

	checking = false;
}

func isSamelist(newlist []string) bool {
	items:=getlist().(map[string]interface{})["result"].(map[string]interface{})["items"]
	if(items == nil){
		return false
	}
	list := getlist().(map[string]interface{})["result"].(map[string]interface{})["items"].([]interface{})
	if (len(newlist) != len(list)) {
		return false;
	}else {
		for i := 0; i < len(list); i++ {
			newlabel := filepath.Base(newlist[i])
			label := list[i].(map[string]interface{})["label"].(string)
			if newlabel != label {
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

func playvideos(videolist []string, flag bool) interface{} {
	//清除播放列表
	clearPlayList()
	//增加列表
	for i := 0; i < len(videolist); i++ {
		additem(videolist[i])
	}
	//未开始强制播放的情况下进行平滑播放,播放完当前视频才切换新视频列表
	if (flag) {
		//强制开始播放
		play()
	}else if (!activePlayer()) {
		//如果未开始播放,则进行播放
		play()
	}
	return getlist()
}

func clearPlayList() interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Playlist.Clear", "params": bson.M{"playlistid": 1}, "id": 1}
	return jsonrpc(params, false)
}

func getPlayeditem() string {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.GetItem", "params": bson.M{"playerid": 1}, "id": 1}
	ret := jsonrpc(params, false);
	result, _ := ret.(map[string]interface{})["result"].(map[string]interface{})["item"].(map[string]interface{})["label"].(string)
	return result
}
func activePlayer() bool {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.GetActivePlayers", "id": 1}
	ret := jsonrpc(params, false).(map[string]interface{})["result"].([]interface{});
	if (len(ret) > 0) {
		return true
	}else {
		return false
	}
}
func playPause() interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.PlayPause", "id": 1, "params":  bson.M{
		"playerid": 1 }}
	return jsonrpc(params, false)
}
func play() interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.Open", "id": 1, "params":  bson.M{
		"item":bson.M{"playlistid":1},
		"options":bson.M{"repeat":"all"}}}
	return jsonrpc(params, false)
}
func playVideo(video string) interface{} {
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
	checkerr(err)
	if (flag) {
		log_print(string(params_str))
	}
	resp, err := http.Post(KodiUrl("/jsonrpc"), "application/json", bytes.NewBuffer([]byte(params_str)))
	checkerr(err)
	body, _ := ioutil.ReadAll(resp.Body)
	if (flag) {
		log_printf("=====ret====%s", string(body))
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

	return jsonrpc(params, true)
}


