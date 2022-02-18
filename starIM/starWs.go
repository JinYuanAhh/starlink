package starIM

import (
	"github.com/gorilla/websocket"
	"time"
)

type User struct {
	Conn     *websocket.Conn
	Platform string
}

var Users = make(map[string][]*User)

func JoinUserPlatform(account string, conn *websocket.Conn, platform string) int {
	if !CheckUserExist(account) {
		Users[account] = []*User{}
	}
	e, i := CheckUserPlatformExist(account, platform)
	if e {
		return i
	} else {
		Users[account] = append(Users[account], &User{conn, platform})
		return len(Users[account]) - 1
	}
}

func DelUserAllPlatform(account string) {
	delete(Users, account)
}
func DelUserPlatform(account string, platform string) {
	e, i := CheckUserPlatformExist(account, platform)
	if e {
		Users[account] = append(Users[account][:i], Users[account][i+1:]...)
	}
}
func DelUsers(account []string) {
	for _, v := range account {
		DelUserAllPlatform(v)
	}
}
func CheckUserPlatformExist(account string, platform string) (bool, int) {
	if !CheckUserExist(account) {
		return false, -1
	}
	exist := false
	index := -1
	for i, v := range Users[account] {
		if v.Platform == platform {
			exist = true
			index = i
		}
	}
	return exist, index
}
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
		Send_Public(websocket.PingMessage, []byte{})
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
