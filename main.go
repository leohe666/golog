package main
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"gopkg.in/yaml.v2"
	"net/http"
	"path"
	"strings"
	"encoding/json"
)

type conf struct {
	Port string `yaml:"port"`
	Path string `yaml:"path"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err  #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}
type JsonResult struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
}
func wirteResult(w http.ResponseWriter, jsonResult JsonResult) {
	w.Header().Set("content-type","text/json")
	ret_json,_ := json.Marshal(jsonResult)
	w.Write(ret_json)
}
func logHandler(w http.ResponseWriter, r *http.Request) {
	paraPath := r.PostFormValue("p")
	msg := r.PostFormValue("m")
	if (r.Method == "GET") {
		query := r.URL.Query()
		msg = query.Get("m")
		paraPath = query.Get("p")
	}
	if (paraPath == "") {
		wirteResult(w, JsonResult{
			Code:   5001,
			Msg:  "文件目录必须填写",
		})
		return
	}
	logFilePath := c.Path+paraPath
	logFileDirPath := path.Dir(logFilePath)
	_, err := os.Stat(logFileDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(logFileDirPath, os.ModePerm)
            if err != nil{
                fmt.Println(logFileDirPath + "文件夹创建失败：", err.Error())
            }
		}
	}
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	gLogger := log.New(file, "", log.Ldate|log.Ltime)
	gLogger.Printf("[%s]%s", machineId, msg)
	wirteResult(w, JsonResult{
		Code:   200,
		Msg:  "处理成功",
	})
}
func getMachineId() {
   	content, err := ioutil.ReadFile("/etc/machine-id")
	if err != nil {
		log.Printf("yamlFile.Get err  #%v", err)
	}
	machineId = strings.Replace(string(content[:]), "\n", "", -1)
}
var c conf
var machineId string
func main() {
	getMachineId()
	c.getConf()
	mux := http.NewServeMux()
	mux.HandleFunc("/", logHandler)
	fmt.Println("logServer started, port:"+c.Port+",logDirRootPath:"+c.Path)
	http.ListenAndServe(":"+c.Port, mux)
   	fmt.Printf(machineId)
}