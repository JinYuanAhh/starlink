package starIM

func SelectUserSecertKeyByAccount(Account string) (string, error) { //以账号查询secretKey
	sqlStr := "SELECT secretKey FROM sl_users WHERE account=?"
	var result string
	err := db.QueryRow(sqlStr, Account).Scan(&result)
	return result, err
}
