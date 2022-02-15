package starIM

import (
	"github.com/gorilla/websocket"
)

type User struct {
	conn     *websocket.Conn
	platform string
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
		if v.platform == platform {
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
