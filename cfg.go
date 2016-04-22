package main
import "time"

func Cfg() map[string]interface{} {
	ret := make(map[string]interface{})
	ret["server"] = "127.0.0.1:3000"
	//ret["server"] = "121.40.199.41:3000"
	ret["kodi"] = "localhost:8080"
	ret["checktime"] = time.Minute*10
	//ret["program_name"] = "videoclient"
	//ret["program_path"] = "/home/pi/"

	ret["program_name"] = "videoclient"
	ret["program_path"] = "d:\\"
	return ret;
}
