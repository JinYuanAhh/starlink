package main

// 处理用户发送的消息

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	IM "github.com/starIM"
	"github.com/tidwall/gjson"
)

func dealTextMsg(conn *websocket.Conn, msg []byte) { //处理消息
	var err error
	sJson := string(msg) //[]byte 转 字符串
	//Account := gjson.Get(sJson, "Userinfo.Account").String() //获取用户账号
	l_msgType := gjson.Get(sJson, "Type").String() //获取信息类型
	switch l_msgType {
	case "SendMsg":
		l_ACI, _ := IM.ParseToken(gjson.Get(sJson, "T").String())
		l_ac := l_ACI.Ac
		l_p := l_ACI.P
		if err != nil {
			IM.Warn("[SendMsg] %s", err)
		}
		ValidUserLogged := l_ac != "" && l_p != ""
		if ValidUserLogged {
			l_msgContent := gjson.Get(sJson, "Msginfo.Content").String()  //获取信息内容
			l_msgContentType := gjson.Get(sJson, "Msginfo.Type").String() //获取内容类型
			switch l_msgContentType {                                     //依据消息的公私分情况
			case "Public": //公共消息（所有人）
				IM.SendStr_Public([]byte(l_msgContent))
			case "Private": //私聊
				account := gjson.Get(sJson, "Msginfo.To").String()
				IM.SendStr_Private([]byte(l_msgContent), account)
			case "Group": //群发
				var accounts []string
				err = json.Unmarshal([]byte(gjson.Get(sJson, "Msginfo.To").String()), &accounts) //获取发送给的账号
				if err != nil {
					//IM.Normal("[Get] get ids when GroupMsg failed") //发送失败
				} else {
					//IM.SendStr_Group([]byte(l_msgContent), accounts) //发送
				}
			default:
				IM.Normal("UnFinished Func")
			}
		} else { //发出警告
			//IM.Warn("[CheckLogged]")
		}
	case "Signup": //注册
		account := gjson.Get(sJson, "Info.Account").String()
		pwd := gjson.Get(sJson, "Info.Pwd").String()
		phoneNumber := gjson.Get(sJson, "Info.PhoneNumber").String()
		_, err := IM.Signup(account, pwd, phoneNumber) //注册, 只需要捕获错误
		if err != nil {
			//IM.Warn("[Signup] %s", err) //警告
		} else {
			//IM.Normal("[Signup] New Acoount At: account: %s, pwd: %s, phoneNumber: %s", account, pwd, phoneNumber) //输出日志
		}
	case "Signin": //登录
		l_Account := gjson.Get(sJson, "Info.Account").String()
		l_Pwd := gjson.Get(sJson, "Info.Pwd").String()
		l_Platform := gjson.Get(sJson, "Info.P").String()
		if e, _ := IM.CheckUserPlatformExist(l_Account, l_Platform); e {
			conn.WriteMessage(1, IM.Msg_Signin_Logged())
		}
		_, err := IM.Signin(l_Account, l_Pwd) //登陆
		if err == nil {                       //成功
			T, err := IM.GenerateToken(l_Account, l_Platform)
			if err != nil {
				IM.Warn("[Signin] %s", err)
			} else {
				conn.WriteMessage(1, IM.Msg_Signin_Success(T))
				IM.JoinUserPlatform(l_Account, conn, l_Platform)
			}
		} else { // 失败
			IM.Normal("[Signin] %s", err)
		}
	case "Logout":
		l_ACI, _ := IM.ParseToken(gjson.Get(sJson, "T").String())
		l_ac := l_ACI.Ac
		l_p := l_ACI.P
		if l_ac == "" || l_p == "" {
			return
		}
		IM.Logout(l_ac, l_p)
	}
}

func dealBinMsg(conn *websocket.Conn, arg []byte, content []byte) {
	Type := gjson.GetBytes(arg, "Type").String()
	switch Type {
	case "File":
		Phrase := gjson.GetBytes(arg, "Phrase").String()
		switch Phrase {
		case "New":
			err := IM.StartFileSave(gjson.GetBytes(arg, "FN").String(), gjson.GetBytes(arg, "CompSegIndex").String(), gjson.GetBytes(arg, "MD5").String())
			if err != nil {
				conn.WriteMessage(1, IM.Msg_File_New_Err(err))
			}else {
				conn.WriteMessage(1, IM.Msg_File_New_Success())
			}
		case "Continue":
			l_bool, err := IM.ContinueFileSave(gjson.GetBytes(arg, "FN").String(), content)
			if err != nil && !l_bool {
				conn.WriteMessage(1, IM.Msg_File_Continue_Err(err))
				IM.Warn("[File Continue] %s", err)
				return
			}else if err != nil && l_bool {
				conn.WriteMessage(1, IM.Msg_File_Continue_Err(err))
				return
			}else if err == nil && l_bool {
				conn.WriteMessage(1, IM.Msg_File_Continue_Success())
				return
			}
		}
	}
}