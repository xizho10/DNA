package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	//"time"
)

var mux sync.Mutex
var i int = 0

// 定义http请求的处理方法
func handlerHello(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

	}
	request := make(map[string]interface{})
	err = json.Unmarshal(body, &request)
	if err != nil {

	}
	mux.Lock()
	i++
	fmt.Println(i)
	mux.Unlock()
	//fmt.Println(r.URL.Path,request)
	resp := map[string]interface{}{
		"Action":  "zx",
		"Result":  "restful-dsggshshhhshitirurteetwt",
		"Error":   0,
		"Desc":    "",
		"Version": "1.0.0",
	}
	ret, _ := json.Marshal(resp)
	//time.Sleep(time.Second / 500)
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(ret))
}

func main() {

	// 注册http请求的处理方法
	http.HandleFunc("/", handlerHello)
	fmt.Println("Server on 50334 start.")
	// 在8086端口启动http服务，会一直阻塞执行
	err := http.ListenAndServe(":50334", nil)
	if err != nil {
		log.Println(err)
	}

	// http服务因故停止后 才会输出如下内容
	fmt.Println("Server on 8086 stopped")
	os.Exit(0)
}
