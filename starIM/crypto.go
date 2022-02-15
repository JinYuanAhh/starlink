package starIM

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func Sha256(str string) string {
	w := sha256.New()
	io.WriteString(w, str)
	bw := w.Sum(nil)
	return hex.EncodeToString(bw)
}

func Base64_Encode(s []byte) string {
	return base64.StdEncoding.EncodeToString(s)
}
func Base64_Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
