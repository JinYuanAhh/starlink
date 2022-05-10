package starIM

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type SFile struct {
	Filename string
	Complete bool
	SegIndex int
	Owner    string
	Date     time.Time
	Args     string
	Size     string
}

func GetFileSha256(fn string) (string, error) {
	hash := sha256.New()
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fc, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	hash.Write(fc)
	Shastr := hex.EncodeToString(hash.Sum(nil))
	return Shastr, nil
}

func IsFileExists(fn string) bool {
	x, err := os.Stat(fn)
	fmt.Println(x)
	return !os.IsNotExist(err)
}

func CreateFile(sha string, fn string, owner string, args string, fullsize string) error {
	abFn := StrConnect("./Resources/Files/", sha)
	if !CheckFileInfoExist(sha) {
		return errors.New("exist")
	} else {
		f, err := os.Create(abFn)
		go f.Close()
		if err == nil {
			sqlStr := "INSERT INTO sl_files (filename, sha, complete, segIndex, owner, date, fullsize, args) VALUES (?,?,?,?,?,?,?,?)"
			_, err = db.Exec(sqlStr, fn, sha, 0, 0, owner, time.Now().Format("2006-01-02 15:04:05"), fullsize, args)
		}
		return err
	}
}

func QueryFile(sha string) (SFile, error) {
	sqlStr := "SELECT (filename, complete, segIndex, owner, date, args, size) FROM sl_files WHERE sha=?"
	s := SFile{}
	err := db.QueryRow(sqlStr, sha).Scan(&s.Filename, &s.Complete, &s.SegIndex, &s.Owner, &s.Date, &s.Args, &s.Size)
	return s, err
}
func CheckFileOwner(sha string, owner string) error {
	sqlStr := "SELECT owner FROM sl_files WHERE sha=? AND owner=?"
	row := db.QueryRow(sqlStr, sha, owner)
	if row.Err() == sql.ErrNoRows {
		return errors.New("not your file or no such file")
	}
	return row.Err()
}

func CheckFileCompleted(sha string) bool {
	sqlStr := "SELECT complete FROM sl_files WHERE sha=?"
	var comp interface{}
	err := db.QueryRow(sqlStr, sha).Scan(&comp)
	return err == nil && comp.([]uint8)[0] == 1
}

func CheckFileInfoExist(sha string) bool {
	sqlStr := "SELECT filename FROM sl_files WHERE sha=?"
	err := db.QueryRow(sqlStr, sha).Scan()
	return err == sql.ErrNoRows
}

func AppendFile(sha string, owner string, content []byte) error {
	if err := CheckFileOwner(sha, owner); err != nil {
		return err
	}
	f, err := os.OpenFile(StrConnect("./Resources/Files/", sha), os.O_APPEND, 0777)
	defer f.Close()
	if err != nil {
		return err
	} else {
		_, err = f.Write(content)
		if err != nil {
			return err
		} else {
			sqlStr := "UPDATE sl_files SET segIndex=segIndex+1 WHERE sha=?"
			_, err = db.Exec(sqlStr, sha)
			return err
		}
	}
}

func CompleteFile(sha string, owner string) error {
	if CheckFileCompleted(sha) {
		return errors.New("completed")
	}
	if err := CheckFileOwner(sha, owner); err != nil {
		return err
	}
	s, err := GetFileSha256(StrConnect("./Resources/Files/", sha))
	if err == nil {
		if strings.EqualFold(s, sha) {
			sqlStr := "UPDATE sl_files SET complete=1 WHERE sha=?"
			_, err = db.Exec(sqlStr, sha)
			return err
		} else {
			return errors.New("sha256 not match")
		}
	} else {
		return err
	}
}
