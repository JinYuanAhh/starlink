package starIM

import (
	"github.com/gorilla/websocket"
	"time"
)

type User struct {
	Conn *websocket.Conn
	//Platform string
	T string //token
}

var Users = make(map[string][]*User)

func JoinUserConn(account string, conn *websocket.Conn, token string) int {
	if !CheckUserExist(account) {
		Users[account] = []*User{}
	}
	//e, i := CheckUserPlatformExist(account, platform)
	//if e {
	//	return i
	//} else {
	Users[account] = append(Users[account], &User{conn, token})
	return len(Users[account]) - 1
	//}
}

func DelUserAllConns(account string) {
	delete(Users, account)
}

//func DelUserPlatform(account string, platform string) {
//	e, i := CheckUserPlatformExist(account, platform)
//	if e {
//		Users[account] = append(Users[account][:i], Users[account][i+1:]...)
//	}
//}
func DelUserConn(account string, T string) {
	for i, v := range Users[account] {
		if v.T == T {
			Users[account] = append(Users[account][i:], Users[account][:i+1]...)
		}
	}
}
func DelUsers(account []string) {
	for _, v := range account {
		DelUserAllConns(v)
	}
}

//func CheckUserPlatformExist(account string, platform string) (bool, int) {
//	if !CheckUserExist(account) {
//		return false, -1
//	}
//	exist := false
//	index := -1
//	for i, v := range Users[account] {
//		if v.Platform == platform {
//			exist = true
//			index = i
//		}
//	}
//	return exist, index
//}
func CheckUserExist(account string) bool {
	exist := false
	_, ok := Users[account]
	exist = ok
	return exist
}

func Ping(p time.Duration) {
	ticker := time.NewTicker(p) //Ping/Pong定时器
	defer ticker.Stop()
	for range ticker.C {
		SendToPublic(websocket.PingMessage, []byte{})
		for _, v := range Users {
			for _, vv := range v {
				vv.Conn.SetReadDeadline(time.Now().Add((time.Second * 10)))
				vv.Conn.SetPongHandler(func(str string) error {
					vv.Conn.SetReadDeadline(time.Now().Add((time.Second * 10)))
					return nil
				})
			}
		}
	}

} //Ping(How long)
