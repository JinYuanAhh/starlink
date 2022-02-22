package main

// 处理用户发送的消息

import (
	"github.com/gorilla/websocket"
	IM "github.com/starIM"
	"github.com/tidwall/gjson"
)

func dealTextMsg(Conn *websocket.Conn, msg []byte, counter *int) { //处理消息
	sJson := string(msg) //[]byte 转 字符串
	//Account := gjson.Get(sJson, "Userinfo.Account").String() //获取用户账号
	l_msgType := gjson.Get(sJson, "Type").String() //获取信息类型
	switch l_msgType {
	case "SendMsg":
		l_ACI, err := IM.ParseToken(gjson.Get(sJson, "T").String())
		if err != nil {
			return
		}
		l_ac := l_ACI.Ac
		if err != nil {
			IM.Warn("[SendMsg] %s", err)
		}
		ValidUserLogged := l_ac != "" && IM.CheckUserSecretKey(l_ac, l_ACI.SecretKey)
		if ValidUserLogged {
			l_msgContent := gjson.Get(sJson, "Info.Content").String()  //获取信息内容
			l_msgContentType := gjson.Get(sJson, "Info.Type").String() //获取内容类型
			switch l_msgContentType {                                  //依据消息的公私分情况
			case "Public": //公共消息（所有人）
				IM.AddPublicMsg(l_ACI.Ac, l_msgContent)
				lPublicInfo := IM.Query_userPublicInfo(l_ac)
				IM.SendToPublic(1, []byte(IM.GenerateJson(map[string]string{
					"Type":               "Message",
					"Info.Type":          "Public",
					"Info.Content":       l_msgContent,
					"Info.ContentType":   "Text",
					"Info.From.Account":  l_ac,
					"Info.From.Avatar":   gjson.Get(lPublicInfo, "Avatar").String(),
					"Info.From.Nickname": gjson.Get(lPublicInfo, "Nickname").String(),
				})))
			case "Private": //私聊

			case "Group": //群发

			default:
				IM.Warn("UnFinished Func")
			}
		} else { //发出警告
			//IM.Warn("[CheckLogged]")
		}
	case "Signup": //注册
		account := gjson.Get(sJson, "Info.Account").String()
		pwd := gjson.Get(sJson, "Info.Pwd").String()
		_, err := IM.Signup(account, pwd) //注册, 只需要捕获错误
		if err != nil {
			//IM.Warn("[Signup] %s", err) //警告
			Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
				"Type":   "Signup",
				"Status": "Error",
				"Err":    err.Error(),
			})))
		} else {
			//IM.Normal("[Signup] New Acoount At: account: %s, pwd: %s, phoneNumber: %s", account, pwd, phoneNumber) //输出日志
			Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
				"Type":   "Signup",
				"Status": "Success",
			})))
		}
	case "Signin": //登录
		l_Account := gjson.Get(sJson, "Info.Account").String()
		l_Pwd := gjson.Get(sJson, "Info.Pwd").String()
		_, err := IM.Signin(l_Account, l_Pwd) //登陆
		if err == nil {                       //成功
			T, err := IM.GenerateToken(l_Account, IM.T_GetUserSecretKey(l_Account))
			if err != nil {
				IM.Warn("[Signin - GenerateToken] %s", err)
			} else {
				Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
					"Type":   "Signin",
					"Status": "Success",
					"T":      T,
				})))
				IM.JoinUserConn(l_Account, Conn, T)
			}
		} else { // 失败
			Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
				"Type":   "Signin",
				"Status": "Error",
				"Err":    err.Error(),
			})))
			if err.Error() == "password is wrong" {
				//*counter = *counter + 1
				//if *counter >= 4 {
				//	Conn.Close()
				//}
			}
		}
	case "Logout":
		l_ACI, err := IM.ParseToken(gjson.Get(sJson, "T").String())
		if err != nil {
			IM.Warn("%s", err)
			return
		}
		l_ac := l_ACI.Ac
		l_sk := l_ACI.SecretKey
		if l_ac == "" || l_sk == "" {
			return
		}
		IM.Logout(l_ac, l_sk)
	case "Reconnect":
		l_ACI, err := IM.ParseToken(gjson.Get(sJson, "T").String())
		if err != nil {
			return
		}
		for _, v := range IM.Users[l_ACI.Ac] {
			if v.T == gjson.Get(sJson, "T").String() {
				v.Conn = Conn
				Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
					"Type":   "Reconnect",
					"Status": "Success",
				})))
				return
			}
		}
		Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
			"Type":   "Reconnect",
			"Status": "Error",
		})))
	}
}
func dealBinMsg(Conn *websocket.Conn, arg []byte, content []byte) {
	Type := gjson.GetBytes(arg, "Type").String()
	switch Type {
	case "File":
		Phrase := gjson.GetBytes(arg, "Phrase").String()
		switch Phrase {
		case "New":
			err := IM.StartFileSave(gjson.GetBytes(arg, "FN").String(), gjson.GetBytes(arg, "CompSegIndex").String(), gjson.GetBytes(arg, "MD5").String())
			if err != nil {
				Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
					"Type":   "Signup",
					"Phrase": "New",
					"Status": "Error",
					"Err":    err.Error(),
				})))
			} else {
				Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
					"Type":   "Signup",
					"Phrase": "New",
					"Status": "Success",
				})))
			}
		case "Continue":
			l_bool, err := IM.ContinueFileSave(gjson.GetBytes(arg, "FN").String(), content)
			if err != nil && !l_bool {
				Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
					"Type":   "Signup",
					"Phrase": "Continue",
					"Status": "Error",
					"Err":    err.Error(),
				})))
				IM.Warn("[File Continue] %s", err)
				return
			} else if err != nil && l_bool {
				Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
					"Type":   "Signup",
					"Phrase": "Continue",
					"Status": "Error",
					"Err":    err.Error(),
				})))
				return
			} else if err == nil && l_bool {
				Conn.WriteMessage(1, []byte(IM.GenerateJson(map[string]string{
					"Type":   "Signup",
					"Phrase": "Continue",
					"Status": "Success",
				})))
				return
			}
		}
	}
}
