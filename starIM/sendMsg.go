package starIM

func SendStr_Public(msg []byte) {
	for _, v := range Users {
		for _, vv := range v {
			vv.Conn.WriteMessage(1, msg)
		}
	}
}
func SendStr_Private(msg []byte, account string) {
	for _, v := range Users[account] {
		v.Conn.WriteMessage(1, msg)
	}
}
func SendStr_Group(msg []byte, account []string) {
	for _, v := range account {
		for _, vv := range Users[v] {
			vv.Conn.WriteMessage(1, msg)
		}
	}
}

func SendBin_Public(msg []byte) {
	for _, v := range Users {
		for _, vv := range v {
			vv.Conn.WriteMessage(2, msg)
		}
	}
}
func SendBin_Private(msg []byte, account string) {
	for _, v := range Users[account] {
		v.Conn.WriteMessage(2, msg)
	}
}
func SendBin_Group(msg []byte, account []string) {
	for _, v := range account {
		for _, vv := range Users[v] {
			vv.Conn.WriteMessage(2, msg)
		}
	}
}

func Send_Public(mtype int, msg []byte) {
	for _, v := range Users {
		for _, vv := range v {
			vv.Conn.WriteMessage(mtype, msg)
		}
	}
}
func Send_Private(mtype int, msg []byte, account string) {
	for _, v := range Users[account] {
		v.Conn.WriteMessage(mtype, msg)
	}
}
func Send_Group(mtype int, msg []byte, account []string) {
	for _, v := range account {
		for _, vv := range Users[v] {
			vv.Conn.WriteMessage(mtype, msg)
		}
	}
}
