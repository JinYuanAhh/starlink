package starIM

func SendToPublic(Mtype int, content []byte) {
	for _, v := range Users {
		for _, vv := range v {
			vv.Conn.WriteMessage(Mtype, content)
		}
	}
}
func SendToAccount(Mtype int, content []byte, account string) {
	for _, v := range Users[account] {
		v.Conn.WriteMessage(Mtype, content)
	}
}
