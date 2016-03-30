package main

import (
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"reflect"
	"path/filepath"
)

func circle() {
	checktime := Cfg()["checktime"].(time.Duration)
	ticker := time.NewTicker(checktime)
	go func() {
		for _  = range ticker.C {
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

	if(result == nil || lists.IsNil()){
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



