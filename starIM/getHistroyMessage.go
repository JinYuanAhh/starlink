package starIM

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strconv"
)

func GetHistroyPublicMessages() string {
	var account, content, l string
	var i = 0
	var msgs = "[]"
	sqlStr := "SELECT fromaccount, content FROM msgs_public WHERE canceled = 0"
	rows, err := db.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		return ""
	} else {
		for rows.Next() {
			err = rows.Scan(&account, &content)
			if err != nil {
				return ""
			}
			lPublicInfo := Query_userPublicInfo(account)
			l, _ = sjson.Set(msgs, StrConnect(strconv.Itoa(i), ".from.avatar"), gjson.Get(lPublicInfo, "Avatar").String())
			l, _ = sjson.Set(l, StrConnect(strconv.Itoa(i), ".from.nickname"), gjson.Get(lPublicInfo, "Nickname").String())
			l, _ = sjson.Set(l, StrConnect(strconv.Itoa(i), ".from.account"), account)
			l, _ = sjson.Set(l, StrConnect(strconv.Itoa(i), ".content"), content)
			msgs = l
			i++
		}
		return msgs
	}
}
