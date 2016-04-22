package main

import (
	"github.com/go-martini/martini"
	"log"
	"path/filepath"
	"os"
	"net/http"
	"io/ioutil"
	"github.com/martini-contrib/render"
	"encoding/json"
	"strconv"
	"runtime"
	"os/exec"
"strings"
	"bytes"
	"time"
)

var checking = false
var rootpath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

func main() {
	//设置日志

	log.SetFlags(log.LstdFlags | log.Llongfile)

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/versionup", versionup)
	//m.Run()
	m.RunOnAddr(":10002")
}

func versionup(r render.Render) {
	cfg := Cfg()
	program_name := cfg["program_name"].(string)
	program_path := cfg["program_path"].(string)
	resp, err := http.Get(HttpUrl("/program/version"))
	checkerr(err)
	body, _ := ioutil.ReadAll(resp.Body)
	var result  map[string]interface{}
	json.Unmarshal(body, &result)
	currentversion := float64(getCurrentVersion())
	log.Println(currentversion, result["version"])
	//不相等就替换,解决版本回退的问题
	if currentversion != result["version"].(float64) {
		//如果小于 , 执行更新
		newFile := downfile(result["src"].(string))
		//停止当前程序
		err:=stopProgram(program_name)
		checkerr(err)
		//重命名为.old
		Rename()
		//将新程序拷贝为当前程序
		log.Println(program_path+program_name,newFile)
		err = os.Link(newFile,program_path+program_name)
		checkerr(err)
		//设置当前程序为可执行
		os.Chmod(program_path+program_name,777)
		//运行当前程序
		go runExe(program_path+program_name)
		checkerr(err)
		//设置当前程序版本
		setCurrentVersion(int(result["version"].(float64)))
	}
	r.JSON(200, result)
}

func runExe(name string){
	cmd :=exec.Command(name)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	if(err!=nil){
		log.Println(err.Error())
		log.Println("当前执行已被更新杀死")
	}
}
func Rename(){
	cfg := Cfg()
	program_name := cfg["program_name"].(string)
	program_path := cfg["program_path"].(string)
	oldpath := program_path+program_name
	if fileexists(oldpath) {
		//无法remove,结束程序需要时间
		err := os.Remove(oldpath)
		for err != nil{
			time.Sleep(time.Millisecond * 100)
			err = os.Remove(oldpath)
		}
		checkerr(err)
	}
	/*
	newpath := oldpath+".old"
	newpath1 := newpath
	idx :=0
	for fileexists(newpath1){
		err := os.Remove(newpath1)
		//checkerr(err)
		if(err !=nil){
			log.Println(err)
			newpath1 = newpath + strconv.Itoa(idx)
			idx++
		}
	}
	if fileexists(oldpath){
		log.Println(oldpath, newpath1)
		err := os.Rename(oldpath, newpath1)
		checkerr(err)
	}
	*/
}

func stopProgram(name string) error{
	//windows
	//wmic process get parentprocessid,name|find "test.exe"
	pid , hasprogrm := findPid(name)
	log.Println(pid , hasprogrm )
	if(hasprogrm){
		program,err := os.FindProcess(pid)
		checkerr(err)
		err =program.Kill()
		checkerr(err)

	}
	return nil
}

func findPid(name string) (pid int ,flag bool){
	log.Printf("name==%v\r\n","\""+name+"\"")
	if runtime.GOOS == "windows" {
		//windows
		//wmic process where name='test.exe' get processid
		cmd := exec.Command("cmd" ,"/c","wmic", "process", "where", "name='"+name+"'","get","processid"  )
		out,err := cmd.CombinedOutput()
		if err !=nil{
			panic(err)
		}
		id := strings.TrimSpace(strings.Split(string(out), "\r\n")[1]);
		if(id == ""){
			flag = false
		}else{
			flag = true
			pid,_ = strconv.Atoi( id)
		}
	}else{
		//linux
		///netstat -lnp | grep /main | awk '{print $7}' | awk -F '/' '{print $1}'
		strs, _ := exec.Command("bash","-c","netstat", "-lnp", "|", "grep", "/"+name, "|", "awk", "'{print", "$7}'", "|", "awk", "-F", "'/'", "'{print", "$1}'").Output();
		str  := strings.TrimSpace(string(strs))
		if(str == "") {
			flag = false
		}else{
			flag = true
			pid, _ = strconv.Atoi(str)
		}
	}
	log.Printf("pid==%v\r\n",pid)
	return
}

//func decodeGBK(val string) string{
//	reader := strings.NewReader(val)
//	transformer := transform.NewReader(reader, simplifiedchinese.GBK.NewDecoder())
//	bytes, err := ioutil.ReadAll(transformer)
//	checkerr(err)
//	return string(bytes)
//}
func getCurrentVersion() int {
	var rootpath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	file, _ := os.Open(rootpath + "/version")
	data := make([]byte, 1000)
	count, err := file.Read(data)
	var version int
	version, err = strconv.Atoi(string(data[:count]))
	if err != nil  {
		version = 0
	}
	return version
}

func setCurrentVersion(version int) {
	var rootpath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	err :=ioutil.WriteFile(rootpath + "/version", []byte(strconv.Itoa(version)), 0644)
	checkerr(err)
}

