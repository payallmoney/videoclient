package main

import (
	"strings"
	"path/filepath"
	"os"
	"net/http"
	"io"
	"crypto/md5"
	"fmt"
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
func downfile(path string,hash string) string {
	realpath := filename(path)
	//检查文件是否存在,如果已存在则不再下载
	if fileexists(realpath){
		//检查hash
		localhash ,_ := ComputeMd5(realpath)
		if(localhash == hash){
			return realpath
		}else{
			os.Remove(realpath);
		}
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

func ComputeMd5(filePath string) (string, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x",hash.Sum(result)), nil
}