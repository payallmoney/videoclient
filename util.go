package main

import (
	"fmt"
	"log"
	"os"
"encoding/json"
)

func HttpUrl(url string) string {
	cfg := Cfg()
	server := fmt.Sprint(cfg["server"])
	var ret string
	if (server[len(server) - 1:] == "/") {
		server = server[0:len(server) - 1]
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
	if (server[len(server) - 1:] == "/") {
		server = server[0:len(server) - 1]
	}
	if (url[0:1] == "/") {
		ret = "http://" + server + url;
	}else {
		ret = "http://" + server + "/" + url;
	}
	return ret;
}

func checkerr(err error) {
	if err != nil {
		var logfile, logfileerr = os.OpenFile(rootpath + "/client.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if logfileerr != nil {
			log.Fatalf("error opening file: %v", logfileerr)
		}
		log.SetOutput(logfile)
		log.Fatal(err)
		logfile.Close();
		log.SetOutput(nil)
	}
}

func log_print(msg string) {
	var logfile, logfileerr = os.OpenFile(rootpath + "/client.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if logfileerr != nil {
		log.Fatalf("error opening file: %v", logfileerr)
	}
	log.SetOutput(logfile)
	log.Println(msg)
	logfile.Close();
	log.SetOutput(nil)
}

func log_printf(format string, msg string) {
	var logfile, logfileerr = os.OpenFile(rootpath + "/client.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if logfileerr != nil {
		log.Fatalf("error opening file: %v", logfileerr)
	}
	log.SetOutput(logfile)
	log.Printf(format + "\r\n", msg)
	logfile.Close();
	log.SetOutput(nil)
}


func js(item interface{}) string {
	params_str, err := json.Marshal(item)
	checkerr(err)
	return (string(params_str))
}
