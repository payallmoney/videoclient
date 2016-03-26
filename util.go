package main

import "fmt"

func HttpUrl(url string) string {
	cfg := Cfg()
	server := fmt.Sprint(cfg["server"])
	var ret string
	if (server[len(server)-1:] == "/") {
		server = server[0:len(server)-1]
	}
	if (url[0:1] == "/") {
		ret = "http://" + server + url;
	}else {
		ret = "http://" + server + "/" + url;
	}
	return ret;
}

func KodiUrl(url string) string {
	cfg := Cfg()
	server := fmt.Sprint(cfg["kodi"])
	var ret string
	if (server[len(server)-1:] == "/") {
		server = server[0:len(server)-1]
	}
	if (url[0:1] == "/") {
		ret = "http://" + server + url;
	}else {
		ret = "http://" + server + "/" + url;
	}
	return ret;
}