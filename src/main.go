package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
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
)

func wsHandler(w http.ResponseWriter, r *http.Request) { //websocket 处理器函数
	conn, err := upgrader.Upgrade(w, r, nil) //升级HTTP协议为websocket
	if err != nil {
		IM.Err("[ERR] Upgrader:", err)
		return
	}

	var SigninWrongCount int                       //SigninWrong Counter
	connectMsg := []byte(`{"Type": "ConnectMsg"}`) //连接成功时发送的消息
	conn.WriteMessage(1, connectMsg)               //发送
	//IM.Normal("[Conn] A New Connection")
	defer conn.Close()
	for { //死循环读消息
		msgType, msg, err := conn.ReadMessage() //没有新消息时会阻塞
		if err != nil {                         //发生错误则关闭连接
			//IM.Err("[ReadMsg] %s", err)
			conn.Close()
			return
		}
		IM.Debug("[Msg] Received:%s\n[Msg] Type:%d", msg, msgType) //DEBUG:: 显示收到的消息
		if msgType == websocket.TextMessage {                      //如果是文本消息 大部分
			if gjson.ValidBytes(msg) { //验证是否为有效的json
				go dealTextMsg(conn, msg, &SigninWrongCount) //处理
			} else {
				//IM.Warn("[Msg] invalid json")
			}
		} else if msgType == websocket.BinaryMessage && bytes.Contains(msg, []byte{'|'}) { //如果是二进制消息 关于文件的
			b_args := bytes.SplitN(msg, []byte{'|'}, 2) //文件消息标准格式： base64|base64 所以先分割成为2份
			//Msg 分割后成为 b_args
			arg, err := IM.Base64_Decode(string(b_args[0])) //解码参数部分
			if err != nil {
				IM.Warn("[Base64 Decode - main.go]%s", err)
				continue
			}
			content, err := IM.Base64_Decode(string(b_args[1])) //解码内容部分
			if err != nil {
				IM.Warn("[Base64 Decode - main.go]%s", err)
				continue
			}
			token, err := IM.ParseToken(gjson.GetBytes(arg, "T").String()) //验证身份
			if err != nil || token.Ac == "" || token.P == "" {
				continue
			}
			go dealBinMsg(conn, arg, content) //处理
		}
	}
}

func main() {
	fmt.Println(color.HiYellowString("-- Star Link Server --"))
	var err error
	//http.HandleFunc("/", wsHandler) //HTTP服务挂载
	//FUNC:: IM.Ping(pingTickTime) //开启Ping/Pong功能
	err = IM.Conn() //连接数据库
	if err != nil { //错误
		IM.Err("[DB]:", err)
	} else {
		IM.Succ("[DB] Connected.")
	}
	InitMux()
	http.ListenAndServe(":8889", R) //监听
}

func init() {
	IM.SetPrefixMode("time")                                                                  //设置logsystem的消息前自动加上时间
	os.MkdirAll(IM.StrConnect(filepath.Dir(os.Args[0]), "/Resources/Files"), os.ModePerm)     //创建必要目录 - 文件
	os.MkdirAll(IM.StrConnect(filepath.Dir(os.Args[0]), "/Resources/FilesInfo"), os.ModePerm) //创建必要目录 - 未完成上传文件信息
}
