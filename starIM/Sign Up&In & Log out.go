package starIM

import (
	"database/sql"
	"errors"
)

func Signup(account string, pwd string) (int, error) { //注册
	sqlStr := "INSERT INTO users (id, account, pwd, phoneNum, secretKey, infos) VALUES (?,?,?,'',?,?)"
	if !T_IsAccountExist(account) { //检测账号重复
		{
			v, err := ValidStrFormat(account, StringFormat_Username)
			if err != nil {
				return -1, err
			} else if v {
				return -1, errors.New("account contains illegal character")
			}
			v, err = ValidStrFormat(pwd, StringFormat_Pwd)
			if err != nil {
				return -1, err
			} else if v {
				return -1, errors.New("pwd contains illegal character")
			}
		} //验证账号或密码是否含有特殊字符
		id, err := T_getUserAmount()
		id++
		if id <= 0 { // 获取现有用户数量失败
			return -1, err
		}
		_, err = db.Exec(sqlStr, id, account, pwd, T_RandString(188), GenerateJson(map[string]string{
			"Avatar":         "https://pic2.zhimg.com/v2-5a2635a7bc78b92de9e3daacb38fe442_r.jpg",
			"Friends.0.type": "@PUBLICROOM",
		}))
		if err != nil { //插入记录失败
			return -1, err
		}
		return id, err
	} else { //账号重复
		return -1, errors.New("account exists")
	}
}

func Signin(account string, pwd string) (string, error) { //登录
	sqlStr := "SELECT secretKey, pwd FROM users WHERE account = ?"
	var ( //查询secretKey和pwd 定义变量
		secretKey string
		password  string
	)
	err := db.QueryRow(sqlStr, account).Scan(&secretKey, &password) //查询
	if err != nil {                                                 //查询失败
		if err == sql.ErrNoRows { //无此账号
			return "", errors.New("user doesn't exists")
		} else { //其他
			return "", err
		}
	} else { //查询成功
		if password == pwd { //比对密码
			err = SetUserOnline(account, 1)
			if err != nil {
				return "", errors.New("setUserOnline failed")
			} else {
				return secretKey, nil //成功
			}
		} else { //失败
			return "", errors.New("password is wrong")
		}
	}
}

func Logout(account string, platform string) {
	DelUserPlatform(account, platform)
	if len(Users[account]) == 0 {
		SetUserOnline(account, 0)
	}
}
