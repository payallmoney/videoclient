package main

import (
	"strings"
	"path/filepath"
	"os"
	"net/http"
	"io"
)


func filename(path string) string {
	idx := strings.LastIndex(path, "/")
	rootpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkerr(err)
	ret := rootpath + "/video" + path[idx:];
	ret, err = filepath.Abs(ret);
	checkerr(err)
	return ret
}
func downfile(path string) string {
	realpath := filename(path)
	//检查文件是否存在,如果已存在则不再下载
	if fileexists(realpath){
		return realpath
	}
	downloadurl := HttpUrl(path)

	out, err := os.Create(realpath)
	checkerr(err)
	defer out.Close()
	resp, err := http.Get(downloadurl)
	checkerr(err)
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
