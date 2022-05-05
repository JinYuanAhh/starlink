package starIM

import (
	"github.com/tidwall/sjson"
)

func GenerateJson(m map[string]string) []byte {
	var s string
	for k, v := range m {
		s, _ = sjson.Set(s, k, []byte(v))
	}
	return []byte(s)
}
