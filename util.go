package main

import "fmt"

func HttpUrl(url string) string {
	cfg := Cfg()
	server := fmt.Sprint(cfg["server"])
	var ret string
	if (url[0:1] == "/") {
		ret = "http://" + server + url;
	}else {
		ret = "http://" + server + "/" + url;
	}
	return ret;
}

