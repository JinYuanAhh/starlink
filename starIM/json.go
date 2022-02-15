package starIM

import "github.com/tidwall/sjson"

func GenerateJson(m map[string]string) string {
	s := "{}"
	for k, v := range m {
		s, _ = sjson.Set(s, k, v)
	}
	return s
}
