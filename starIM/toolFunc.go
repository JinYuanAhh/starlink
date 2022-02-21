package starIM

import (
	"database/sql"
	"math/rand"
	"time"
	"unsafe"
)

// 工具函数
func T_getUserAmount() (int, error) { //获取用户数量
	var amount int
	sqlStr := "SELECT COUNT(id) AS amount FROM users"
	err := db.QueryRow(sqlStr).Scan(&amount)
	if err == nil {
		return amount, err
	} else {
		return -1, err
	}
}
func T_IsAccountExist(account string) bool { //判断账户是否存在
	sqlStr := "SELECT id FROM users WHERE account=?"
	var id int
	return !(db.QueryRow(sqlStr, account).Scan(&id) == sql.ErrNoRows)
}
func T_RandString(n int) string { //随机生成字符串 skey用
	const letterBytes = "-_1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func CheckTokenValid(token string) bool {
	_, err := ParseToken(token)
	if err != nil {
		return false
	} else {
		return true
	}
}
func Query_userPublicInfo(account string) string {
	var str string
	sqlStr := "SELECT publicInfo FROM users WHERE account=?"
	err := db.QueryRow(sqlStr, account).Scan(&str)
	if err != nil {
		return "{}"
	}
	return str
}
func Query_userPrivateInfo(token string) string {
	var str string
	l_ACI, err := ParseToken(token)
	if err != nil {
		return ""
	} else {
		sqlStr := "SELECT privateInfo FROM users WHERE account=?"
		err = db.QueryRow(sqlStr, l_ACI.Ac).Scan(&str)
		if err != nil {
			return "{}"
		}
		return str
	}
}
