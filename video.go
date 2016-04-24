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
	videocheck()
	checktime := Cfg()["checktime"].(time.Duration)
	ticker := time.NewTicker(checktime)
	go func() {
		for _  = range ticker.C {
			videocheck()
		}
	}()
}

func videocheck() {
	log_print(" do .. videocheck...")
	if checking {
		log_print(" do .. videocheck...return  ... checking  ")

		return
	}
	checking = true
	id :=cpuid()

	resp, err := http.Get(HttpUrl("/video/list/" + id))
	checkerr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkerr(err)
	var result  interface{}
	err =json.Unmarshal(body, &result)
	checkerr(err)
	lists := reflect.ValueOf(result)
	if(result == nil || lists.IsNil()){
		checking = false
		//设备尚未注册
		return
	}
	files := make([]string, lists.Len())
	hashs := make([]string, lists.Len())
	for i := 0; i < lists.Len(); i++ {
		path := lists.Index(i).Elem().MapIndex(reflect.ValueOf("src")).Elem().String()
		hash := lists.Index(i).Elem().MapIndex(reflect.ValueOf("hash")).Elem().String()
		files[i] = filename(path)
		hashs[i] = hash
	}
	//检查一遍视频
	for i := 0; i < lists.Len(); i++ {
		path := lists.Index(i).Elem().MapIndex(reflect.ValueOf("src")).Elem().String()
		downfile(path,hashs[i])
	}
	if(!isSamelist(files)){
		//播放列表不同时
		playvideos(files, false)
	}else if (!activePlayer()) {
		//如果未开始播放,则进行播放
		play()
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

