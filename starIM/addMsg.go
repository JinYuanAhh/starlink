package starIM

import "time"

func AddPublicMsg(fromAccount string, content string) (bool, error) {
	sqlStr := "INSERT INTO msgs_public (time, fromAccount, content, canceled) VALUES (?, ?, ?, 0)"
	_, err := db.Exec(sqlStr, time.Now().Format("2006-01-02 15:04:05"), fromAccount, content)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
