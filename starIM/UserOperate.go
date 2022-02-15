package starIM

func SetUserOnline(account string, online int) error { //更新在线状态
	sqlStr := "UPDATE users SET online=? WHERE account=?"
	_, err := db.Exec(sqlStr, online, account)
	return err
}
