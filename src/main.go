package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	IM "github.com/starIM"
	"github.com/tidwall/gjson"
)

var (
	upgrader = websocket.Upgrader{ //升级HTTP协议为websocket的结构
		ReadBufferSize:  2560,
		WriteBufferSize: 4608,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	pingTickTime = time.Minute //每隔 pingTickTime 时间发送一次Ping消息
	R            = mux.NewRouter()
	L            sync.Mutex
)

func wsHandler(w http.ResponseWriter, r *http.Request) { //websocket 处理器函数

	var Connection = new(IM.Connection)
	conn, err := upgrader.Upgrade(w, r, nil) //升级HTTP协议为websocket
	if err != nil {
		IM.Err("[ERR] Upgrader:", err)
		return
	}
	Connection.Conn = conn
	//IM.Normal("[Conn] A New Connection")
	defer conn.Close()
	go dealChanMsg(Connection)
	for { //死循环读消息
		msgType, msg, err := conn.ReadMessage() //没有新消息时会阻塞
		IM.DEBUGCOUNTER++
		IM.Debug("%d", IM.DEBUGCOUNTER)
		if err != nil { //发生错误则关闭连接
			conn.Close()
			return
		}
		//IM.Debug("%s", msg)                                            //DEBUG:: 显示收到的消息
		if msgType == websocket.TextMessage && gjson.ValidBytes(msg) { //如果是文本消息 大部分 //验证是否为有效的json
			go dealTextMsg(Connection, msg) //处理
		} else if msgType == websocket.BinaryMessage && bytes.Contains(msg, []byte{'|'}) { //如果是二进制消息 关于文件的
			b_args := bytes.SplitN(msg, []byte{'|'}, 2) //文件消息标准格式： base64|bin 所以先分割成为2份
			//Msg 分割后成为 b_args
			arg, err := IM.Base64_Decode(string(b_args[0])) //解码参数部分
			if err != nil {
				IM.Warn("[Base64 Decode - main.go]%s", err)
				continue
			}
			go dealFileMsg(Connection, arg, b_args[1]) //处理
		} else if msgType == websocket.BinaryMessage {
			go dealBinMsg(Connection, msg)
		}
	}
}

func main() {
	fmt.Println(color.HiYellowString("-- Star Link Server --"))
	var err error
	//http.HandleFunc("/", wsHandler) //HTTP服务挂载
	err = IM.Conn() //连接数据库
	if err != nil { //错误
		IM.Err("[DB]:", err)
	} else {
		IM.Succ("[DB] Connected.")
	}
	InitMux()
	http.ListenAndServe(":8889", R) //监听
}

func init() { //设置logsystem的消息前自动加上时间
	os.MkdirAll(IM.StrConnect(filepath.Dir(os.Args[0]), "/Resources/Files"), os.ModePerm) //创建必要目录 - 文件
	//os.MkdirAll(IM.StrConnect(filepath.Dir(os.Args[0]), "/Resources/FilesInfo"), os.ModePerm) //创建必要目录 - 未完成上传文件信息
}
