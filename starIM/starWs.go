package starIM

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type User struct {
	Conn     *websocket.Conn
	Platform string
	Pinged   bool
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
		Users[account] = append(Users[account], &User{conn, platform, false})
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

func TickForPing(p time.Duration, c chan map[string]User) {
	ticker := time.NewTicker(p) //Ping/Pong定时器
	defer ticker.Stop()
	fmt.Println("d")
	for range ticker.C {
		tc := make(chan int) //time chan
		Send_Public(websocket.PingMessage, []byte{})
		fmt.Println("d")
		go func() { time.Sleep(time.Second * 15); tc <- 0 }()
		select {
		case <-tc: //15秒的最大等待时间到了
			//fmt.Println("tc")
			close(c)
			for k, v := range Users {
				for _, u := range v {
					if !u.Pinged {
						DelUserPlatform(k, u.Platform)
					} else {
						u.Pinged = false
					}
				}
			}
			c = make(chan map[string]User)
		case l_value := <-c: //有新用户回复了这个消息
			//fmt.Println("c")
			for k, v := range l_value {
				for _, u := range Users[k] {
					if u.Platform == v.Platform && u.Conn == v.Conn {
						u.Pinged = true
					}
				}
			}
		}
	}

}
