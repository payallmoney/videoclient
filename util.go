package main

import (
	"log"
	"os"
"encoding/json"
	"path/filepath"
	"io"
	"runtime/debug"
	"reflect"
)

func HttpUrl(url string) string {
	cfg := Cfg()
	server := cfg["server"].(string)
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
	server := cfg["kodi"].(string)
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
	var rootpath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		var logfile, logfileerr = os.OpenFile(rootpath + "/client.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if logfileerr != nil {
			log.Fatalf("error opening file: %v", logfileerr)
		}
		mWriter := io.MultiWriter(os.Stdout, logfile)
		log.SetOutput(mWriter)

		log.Println(err)
		log.Println(string(debug.Stack()))
		logfile.Close();
		log.SetOutput(os.Stdout)
	}
}


func log_print(msg string) {
	var logfile, logfileerr = os.OpenFile(rootpath + "/client.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if logfileerr != nil {
		log.Fatalf("error opening file: %v", logfileerr)
	}
	mWriter := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(mWriter)
	log.Println(msg)
	logfile.Close();
	log.SetOutput(os.Stdout)
}

func log_printf(format string, msg string) {
	var logfile, logfileerr = os.OpenFile(rootpath + "/client.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if logfileerr != nil {
		log.Fatalf("error opening file: %v", logfileerr)
	}
	mWriter := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(mWriter)
	log.Printf(format + "\r\n", msg)
	logfile.Close();
	log.SetOutput(os.Stdout)
}


func js(item interface{}) string {
	params_str, err := json.Marshal(item)
	checkerr(err)
	return (string(params_str))
}

func IsZero(val interface{}) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && IsZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && IsZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}