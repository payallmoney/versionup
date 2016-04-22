package main

import (
	"fmt"
	"log"
	"os"
"encoding/json"
	"path/filepath"
	"io"
	"runtime/debug"
	"reflect"
	"net/http"
"strings"
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

	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
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

func filename(path string) string {
	idx := strings.LastIndex(path, "/")
	rootpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkerr(err)
	ret := rootpath + "/program" + path[idx:];
	ret, err = filepath.Abs(ret);
	checkerr(err)
	return ret
}

func fileexists(file string) bool{
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false;
	}else{
		return true;
	}
}

func checkerr(err error) {
	var rootpath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		var logfile, logfileerr = os.OpenFile(rootpath + "/up.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if logfileerr != nil {
			log.Fatalf("error opening file: %v", logfileerr)
		}
		mWriter := io.MultiWriter(os.Stdout, logfile)
		log.SetOutput(mWriter)

		log.Println(err)
		log.Println(string(debug.Stack()))
		logfile.Close();
		log.SetOutput(os.Stdout)
		os.Exit(0)
	}
}
