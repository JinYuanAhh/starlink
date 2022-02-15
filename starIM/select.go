package starIM

func SelectUserIdByAccount(account string) (int, error) { //以账号查询id
	sqlStr := "SELECT id FROM users WHERE account=?"
	var result int
	err := db.QueryRow(sqlStr, account).Scan(&result)
	return result, err
}
func SelectUserSecertKeyByAccount(Account string) (string, error) { //以账号查询secretKey
	sqlStr := "SELECT secretKey FROM users WHERE account=?"
	var result string
	err := db.QueryRow(sqlStr, Account).Scan(&result)
	return result, err
}
