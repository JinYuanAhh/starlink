package starIM

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"time"
)

var (
	PrefixMode = "default"
)

func SetPrefixMode(mode string) { //设置前缀 （时间等）
	PrefixMode = mode
}
func TacklePrefix() string {
	if PrefixMode == "" {
		return ""
	} else if PrefixMode == "time" {
		return time.Now().Format("[2006-01-02 15:04:05]")
	}
	return ""
}

func StrConnect(s ...string) string {
	var b strings.Builder
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

func Err(msg string, f ...interface{}) {

	fmt.Printf(StrConnect(TacklePrefix(), color.HiRedString(msg), "\n"), f...)
}

func Warn(msg string, f ...interface{}) {
	fmt.Printf(StrConnect(TacklePrefix(), color.HiYellowString(msg), "\n"), f...)
}

func Succ(msg string, f ...interface{}) {
	fmt.Printf(StrConnect(TacklePrefix(), color.HiGreenString(msg), "\n"), f...)
}

func Normal(msg string, f ...interface{}) {
	fmt.Printf(StrConnect(TacklePrefix(), color.WhiteString(msg), "\n"), f...)
}

func ErrF(msg string, f ...interface{}) {
	fmt.Printf(StrConnect(TacklePrefix(), color.HiRedString(msg)), f...)
}

func WarnF(msg string, f ...interface{}) {
	fmt.Printf(StrConnect(TacklePrefix(), color.HiYellowString(msg)), f...)
}

func SuccF(msg string, f ...interface{}) {
	fmt.Printf(StrConnect(TacklePrefix(), color.HiGreenString(msg)), f...)
}

func NormalF(msg string, f ...interface{}) {
	fmt.Printf(StrConnect(TacklePrefix(), color.WhiteString(msg)), f...)
}

func Debug(msg string, f ...interface{}) {
	fmt.Printf(StrConnect(color.BlueString(msg), "\n"), f...)
}
