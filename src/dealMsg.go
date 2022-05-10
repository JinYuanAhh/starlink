package main

// 处理用户发送的消息

import (
	"github.com/gorilla/websocket"
	IM "github.com/starIM"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func dealTextMsg(Connection *IM.Connection, msg []byte, counter *int) { //处理消息
	//Account := gjson.GetBytes(msg, "Userinfo.Account").String() //获取用户账号
	l_msgType := gjson.GetBytes(msg, "Type").String() //获取信息类型
	switch l_msgType {
	case "Message":
		ValidUserLogged := Connection.Account != ""
		if ValidUserLogged {
			l_msgContent := gjson.GetBytes(msg, "Info.Content").String()  //获取信息内容
			l_msgContentType := gjson.GetBytes(msg, "Info.Type").String() //获取内容类型
			switch l_msgContentType {                                     //依据消息的公私分情况
			case "Public": //公共消息（所有人）
				go IM.AddPublicMsg(Connection.Account, l_msgContent)
				IM.MessageQuene <- IM.GenerateJson(map[string]string{
					"Type":              "Message",
					"Info.Type":         "Public",
					"Info.Content":      l_msgContent,
					"Info.ContentType":  "Text",
					"Info.From.Account": Connection.Account,
				})
			case "Private": //私聊

			case "Group": //群发

			default:
				IM.Err("UnFinished Func")
			}
		} else { //发出警告
			//IM.Warn("[CheckLogged]")
		}
	case "Signup": //注册
		account := gjson.GetBytes(msg, "Info.Account").String()
		pwd := gjson.GetBytes(msg, "Info.Pwd").String()
		_, err := IM.Signup(account, pwd) //注册, 只需要捕获错误
		if err != nil {
			//IM.Warn("[Signup] %s", err) //警告
			go ConnWriteMessage(Connection.Conn, 1, []byte(IM.GenerateJson(map[string]string{
				"Type":        "Signup",
				"Info.Status": "Error",
				"Error":       err.Error(),
			})))
		} else {
			//IM.Normal("[Signup] New Acoount At: account: %s, pwd: %s, phoneNumber: %s", account, pwd, phoneNumber) //输出日志
			go ConnWriteMessage(Connection.Conn, 1, []byte(IM.GenerateJson(map[string]string{
				"Type":        "Signup",
				"Info.Status": "Success",
			})))
		}
	case "Signin": //登录
		l_Account := gjson.GetBytes(msg, "Info.Account").String()
		l_Pwd := gjson.GetBytes(msg, "Info.Pwd").String()
		l_sk := IM.T_GetUserSecretKey(l_Account)
		_, err := IM.Signin(l_Account, l_Pwd) //登陆
		if err == nil {                       //成功
			T, err := IM.GenerateToken(l_Account, l_sk)
			if err != nil {
				IM.Warn("[Signin - GenerateToken] %s", err)
			} else {
				go ConnWriteMessage(Connection.Conn, 1, []byte(IM.GenerateJson(map[string]string{
					"Type":        "Signin",
					"Info.Status": "Success",
					"T":           T,
				})))
				Connection.T = T
				Connection.Account = l_Account
			}
		} else { // 失败
			go ConnWriteMessage(Connection.Conn, 1, []byte(IM.GenerateJson(map[string]string{
				"Type":        "Signin",
				"Info.Status": "Error",
				"Error":       err.Error(),
			})))
			if err.Error() == "password is wrong" {
				//*counter = *counter + 1
				//if *counter >= 4 {
				//	Connection.Conn.Close()
				//}
			}
		}
	case "Logout":
		Connection.Account = ""
		Connection.T = ""
		Connection.Skey = ""
	case "Reconnect":
		l_T := gjson.GetBytes(msg, "T").String()
		l, err := IM.ParseToken(l_T)
		if err != nil {
			go ConnWriteMessage(Connection.Conn, 1, []byte(IM.GenerateJson(map[string]string{
				"Type":        "Reconnect",
				"Info.Status": "Error",
			})))
		}
		Connection.T = l_T
		Connection.Account = l.Ac
		go ConnWriteMessage(Connection.Conn, 1, []byte(IM.GenerateJson(map[string]string{
			"Type":        "Reconnect",
			"Info.Status": "Success",
		})))
		return
	case "FriendRequest":
		ValidUserLogged := Connection.Account != ""
		if ValidUserLogged {
			rs := gjson.GetManyBytes(msg, "Info.Type", "Info.To", "Info.Addition", "Info.From")
			switch rs[0].String() { //switch 类型
			case "New":
				succ, err := IM.NewFriendRequest(Connection.Account, rs[1].String(), rs[2].String())
				if succ {
					l, _ := sjson.SetBytes(msg, "Info.From", Connection.Account)
					IM.MessageQuene <- l
				} else {
					go ConnWriteMessage(Connection.Conn, 1, []byte(IM.GenerateJson(map[string]string{
						"Type":        "FriendRequest",
						"Info.Status": "Error",
						"Error":       err.Error(),
					})))
				}
			case "Tackle":
				if ok, err := IM.TackleFriendRequest(rs[3].String(), Connection.Account, rs[2].String()); ok {
					if rs[2].String() == "Agree" {
						l, _ := sjson.Set(IM.UserPrivateInfo(rs[3].String()), IM.StrConnect("friends.", rs[1].String(), ".Relation"), "N") //Normal Friend 普通朋友
						err1 := IM.SetUserPrivateInfo(rs[3].String(), l)
						l, _ = sjson.Set(IM.UserPrivateInfo(rs[1].String()), IM.StrConnect("friends.", rs[3].String(), ".Relation"), "N")
						err2 := IM.SetUserPrivateInfo(rs[1].String(), l)
						if err != nil || err2 != nil {
							IM.Err("%s", err1.Error())
							go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
								"Type":        "FriendRequest",
								"Info.Type":   "Tackle",
								"Info.Status": "Error",
								"Info.Error":  err1.Error() + err2.Error(),
							}))
						} else {
							go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
								"Type":        "FriendRequest",
								"Info.Type":   "Tackle",
								"Info.Status": "Success",
							}))
							l, _ := sjson.SetBytes(msg, "Info.To", Connection.Account)
							IM.MessageQuene <- l
						}
					} else if rs[2].String() == "Refuse" {
						b, e := IM.DeleteFriendRequest(rs[3].String(), Connection.Account)
						if b {
							go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
								"Type":        "FriendRequest",
								"Info.Type":   "Tackle",
								"Info.Status": "Success",
							}))
							l, _ := sjson.SetBytes(msg, "Info.To", Connection.Account)
							IM.MessageQuene <- l
						} else {
							go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
								"Type":        "FriendRequest",
								"Info.Type":   "Tackle",
								"Info.Status": "Error",
								"Info.Error":  e.Error(),
							}))
						}
					}

				} else {
					go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
						"Type":        "FriendRequest",
						"Info.Type":   "Tackle",
						"Info.Status": "Error",
						"Info.Error":  err.Error(),
					}))
				}

			}
		}
	}
}
func dealBinMsg(Connection *IM.Connection, arg []byte, content []byte) {
	Type := gjson.GetBytes(arg, "Type").String()
	switch Type {
	case "File":
		itype := gjson.GetBytes(arg, "Info.Type").String()
		switch itype { //Type
		case "New":
			rs := gjson.GetManyBytes(arg, "Info.Filename", "Info.Sha", "Info.Args", "Info.Size")
			err := IM.CreateFile(rs[1].String(), rs[0].String(), Connection.Account, rs[2].String(), rs[3].String())
			if err != nil {
				go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
					"Type":        "File",
					"Info.Type":   "New",
					"Info.Status": "Error",
					"Info.Error":  err.Error(),
				}))
			}
		case "Append":
			if err := IM.AppendFile(gjson.GetBytes(arg, "Info.Sha").String(), Connection.Account, content); err != nil {
				go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
					"Type":        "File",
					"Info.Type":   "Append",
					"Info.Status": "Error",
					"Info.Error":  err.Error(),
				}))
			} else {
				go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
					"Type":        "File",
					"Info.Type":   "New",
					"Info.Status": "Success",
				}))
			}
		case "Complete":
			if err := IM.CompleteFile(gjson.GetBytes(arg, "Info.Sha").String(), Connection.Account); err != nil {
				go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
					"Type":        "File",
					"Info.Type":   "Complete",
					"Info.Status": "Error",
					"Info.Error":  err.Error(),
				}))
			} else {
				go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
					"Type":        "File",
					"Info.Type":   "Complete",
					"Info.Status": "Success",
				}))
			}
		}
	}
}

func dealChanMsg(Connection *IM.Connection) { //处理 MessageQuene 的 Msg
	var c []byte
	IM.QueneCond.L.Lock()
	IM.QueneCond.Wait()
	IM.QueneWG.Add(1)
	IM.QueneCond.L.Unlock()
	for {
		IM.QueneCond.L.Lock()
		//select {
		//case c := <-IM.MessageQuene:
		//	switch gjson.GetBytes(c, "Info.Type").String() {
		//	case "Public":
		//		Connection.Conn.WriteMessage(1, c)
		//	}
		//case c := <-IM.FileQuene:
		//	_ = c //DEBUG - DELETE
		//case c := <-IM.FunctionQuene:
		//	rs := gjson.GetManyBytes(c, "Type", "Info.Type", "Info.From", "Info.To", "Info.Addition")
		//	switch rs[0].String() { //Type
		//	case "FriendRequest":
		//		if (rs[1].String() == "New") {
		//			if (rs[2].String() == Connection.Account) { //发出者
		//				Connection.Conn.WriteMessage(1, IM.GenerateJson(map[string]string{
		//					"Type": "FriendRequest",
		//					"Info.Type": "New",
		//					"Info.Status": "Success",
		//				}))
		//			}else if (rs[3].String() == Connection.Account) { //接收者
		//				Connection.Conn.WriteMessage(1, IM.GenerateJson(map[string]string{
		//					"Type": "FriendRequest",
		//					"Info.Type": "New",
		//					"From": rs[2].String(),
		//					"Addition": rs[4].String(),
		//				}))
		//			}
		//		}
		//	}
		//}
		IM.L.Lock()
		c = IM.MessageWaiting
		IM.QueneWG.Done()
		IM.L.Unlock()
		IM.QueneCond.L.Unlock()
		rs := gjson.GetManyBytes(c, "Type", "Info.Type", "Info.From", "Info.To", "Info.Addition")
		switch rs[0].String() {
		case "FriendRequest":
			if rs[1].String() == "New" {
				if rs[2].String() == Connection.Account { //发出者
					go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
						"Type":        "FriendRequest",
						"Info.Type":   "New",
						"Info.Status": "Success",
					}))
				} else if rs[3].String() == Connection.Account { //接收者
					go ConnWriteMessage(Connection.Conn, 1, IM.GenerateJson(map[string]string{
						"Type":      "FriendRequest",
						"Info.Type": "New",
						"From":      rs[2].String(),
						"Addition":  rs[4].String(),
					}))
				}
			} else if rs[1].String() == "Tackle" {
				if rs[2].String() == Connection.Account { //发出者
					go ConnWriteMessage(Connection.Conn, 1, c)
				}
			}
		case "Message":
			if rs[1].String() == "Public" {
				go ConnWriteMessage(Connection.Conn, 1, c)
			}
		}
		IM.QueneCond.L.Lock()
		IM.QueneCond.Wait()
		IM.QueneCond.L.Unlock()
		IM.QueneWG.Add(1)

	}
}

func ConnWriteMessage(conn *websocket.Conn, mtype int, data []byte) {
	L.Lock()
	conn.WriteMessage(mtype, data)
	L.Unlock()
}
