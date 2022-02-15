package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	IM "github.com/starIM"
	"github.com/tidwall/gjson"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  2560,
		WriteBufferSize: 4608,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		IM.Err("[ERR] Upgrader:", err)
		return
	}
	connectMsg := []byte(`{"Type": "ConnectMsg"}`) //连接成功时发送消息
	conn.WriteMessage(1, connectMsg)               //发送
	//IM.Normal("[Conn] A New Connection")
	defer conn.Close()
	for {
		msgType, msg, err := conn.ReadMessage()

		if err != nil {
			//IM.Err("[ReadMsg] %s", err)
			conn.Close()
			return
		}
		//IM.Normal("[Msg] Received:%s\n[Msg] Type:%d", msg, msgType)
		if msgType == 1 {
			vd := gjson.ValidBytes(msg)
			if vd {
				go dealTextMsg(conn, msg)

			} else {
				//IM.Warn("[Msg] invalid json")
			}
		}
		if msgType == 2 && bytes.Contains(msg, []byte{'|'}) {
			b_args := bytes.SplitN(msg, []byte{'|'}, 2)
			//Msg 分割后成为 args
			arg, err := IM.Base64_Decode(string(b_args[0]))
			if err != nil {
				IM.Warn("[Base64 Decode - main.go]%s", err)
				continue
			}
			content, err := IM.Base64_Decode(string(b_args[1]))
			if err != nil {
				IM.Warn("[Base64 Decode - main.go]%s", err)
				continue
			}
			token, err := IM.ParseToken(gjson.GetBytes(arg, "T").String())
			if err != nil || token.Ac == "" || token.P == "" {
				continue
			}
			dealBinMsg(conn, arg, content)
		}
	}
}

func main() {
	fmt.Println(color.HiYellowString("-- Star Link Server --"))
	var err error
	http.HandleFunc("/", wsHandler) //HTTP服务挂载
	err = IM.Conn()                 //连接数据库
	if err != nil {                 //错误
		IM.Err("[DB]:", err)
	} else {
		IM.Succ("[DB] Connected.")
	}
	http.ListenAndServe(":8889", nil) //监听
}

func init() {
	IM.SetPrefixMode("time")
	os.MkdirAll(IM.StrConnect(filepath.Dir(os.Args[0]), "/Resources/Files"), os.ModePerm)
	os.MkdirAll(IM.StrConnect(filepath.Dir(os.Args[0]), "/Resources/FilesInfo"), os.ModePerm)
}
