package starIM

import (
	"github.com/gorilla/websocket"
	"sync"
)

var (
	MessageQuene   = make(chan []byte)
	MessageWaiting []byte
	QueneLocker    sync.Mutex
	QueneCond      = sync.NewCond(&QueneLocker)
	QueneWG        sync.WaitGroup
	L              sync.Mutex
)

type Connection struct {
	Conn    *websocket.Conn
	Account string
	T       string
	Skey    string
}

func init() {
	go TackleQuene()
}

func TackleQuene() {
	var c []byte
	for {
		c = <-MessageQuene
		MessageWaiting = c
		QueneCond.Broadcast()
		QueneWG.Wait()
	}
}
