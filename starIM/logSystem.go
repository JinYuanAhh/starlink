package starIM

import (
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
)

var (
	dbg = log.New(os.Stdout, "<DEBUG> ", log.Lshortfile|log.Ltime|log.Ldate)
)

func SetFlags(flag int) {
	log.SetFlags(flag)
}

func StrConnect(s ...string) string {
	var b strings.Builder
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

func Err(msg string, f ...interface{}) {

	log.Printf(StrConnect(color.HiRedString(msg), "\n"), f...)
}

func Warn(msg string, f ...interface{}) {
	log.Printf(StrConnect(color.HiYellowString(msg), "\n"), f...)
}

func Succ(msg string, f ...interface{}) {
	log.Printf(StrConnect(color.HiGreenString(msg), "\n"), f...)
}

func Normal(msg string, f ...interface{}) {
	log.Printf(StrConnect(color.WhiteString(msg), "\n"), f...)
}

func ErrF(msg string, f ...interface{}) {
	log.Printf(StrConnect(color.HiRedString(msg)), f...)
}

func WarnF(msg string, f ...interface{}) {
	log.Printf(StrConnect(color.HiYellowString(msg)), f...)
}

func SuccF(msg string, f ...interface{}) {
	log.Printf(StrConnect(color.HiGreenString(msg)), f...)
}

func NormalF(msg string, f ...interface{}) {
	log.Printf(StrConnect(color.WhiteString(msg)), f...)
}

func Debug(msg string, f ...interface{}) {
	dbg.Printf(StrConnect(color.BlueString(msg), "\n"), f...)
}
