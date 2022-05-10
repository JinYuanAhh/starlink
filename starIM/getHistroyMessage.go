package starIM

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strconv"
)

func GetHistroyPublicMessages(ps int, pn int) string {
	var account, content, l string
	var i = 0
	var msgs = "[]"
	sqlStr := "SELECT fromaccount, content FROM sl_msgs_public WHERE canceled = 0 ORDER BY time DESC"
	rows, err := db.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		return ""
	} else {
		for ; i < (pn-1)*ps; i++ {
			rows.Next()
		}
		i = 0
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
			if i == ps {
				break
			}
		}
		return msgs
	}
}
