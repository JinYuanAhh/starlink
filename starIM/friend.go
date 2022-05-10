package starIM

import "errors"

func Is_FriendRequestExists(from string, to string) (bool, string) {
	sqlStr := "SELECT addition FROM sl_friendrequests WHERE fromAc=? AND toAc=?"
	var addition string
	err := db.QueryRow(sqlStr, from, to).Scan(&addition)
	if err != nil {
		return false, ""
	} else {
		return true, addition
	}
}
func NewFriendRequest(from string, to string, addition string) (bool, error) {
	if from == to {
		return false, errors.New("can't be self")
	}
	if from == "" || to == "" {
		return false, errors.New("can't be blank")
	}
	e, a := Is_FriendRequestExists(from, to)
	if e { //已存在
		if a == addition { //完全相同
			return false, errors.New("exists")
		} else {
			sqlStr := "UPDATE sl_friendrequests SET addition=? WHERE fromAc=? AND toAc=?"
			_, err := db.Exec(sqlStr, addition, from, to)
			if err != nil {
				return false, err
			} else {
				return true, errors.New("changed")
			}
		}
	} else {
		sqlStr := "INSERT INTO sl_friendrequests (fromAc, toAc, addition) VALUES (?, ?, ?)"
		_, err := db.Exec(sqlStr, from, to, addition)
		if err != nil {
			return false, nil
		}
		return true, nil
	}
}
func TackleFriendRequest(from string, to string, status string) (bool, error) {
	e, _ := Is_FriendRequestExists(from, to)
	if e { //已存在
		sqlStr := "UPDATE sl_friendrequests SET status=? WHERE fromAc=? AND toAc=?"
		_, err := db.Exec(sqlStr, status, from, to)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, errors.New("not exists")
}

func DeleteFriendRequest(from string, to string) (bool, error) {
	e, _ := Is_FriendRequestExists(from, to)
	if e { //已存在
		sqlStr := "DELETE FROM sl_friendrequests WHERE fromAc=? AND toAc=?"
		_, err := db.Exec(sqlStr, from, to)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, errors.New("not exist")
}
