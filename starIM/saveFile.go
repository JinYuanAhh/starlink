package starIM

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"time"
)

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
	_, err := os.Stat(fn)
	return os.IsNotExist(err)
}

func CreateFile(sha string, fn string, owner string) error {
	abFn := StrConnect("./Resources/Files/", sha)
	if IsFileExists(abFn) {
		return errors.New("exist")
	} else {
		_, err := os.Create(abFn)
		if err != nil {
			sqlStr := "INSERT INTO sl_files (filename, sha, complete, segIndex, owner, date) VALUES (?,?,?,'',?,?)"
			db.Exec(sqlStr, fn, sha, 0, 0, owner, time.Now().Format("2006-01-02 15:04:05"))
		}
		return err
	}
}

func QueryFile(sha string) error {
	sqlStr := "SELECT (filename, complete, segIndex, owner, date) FROM sl_files WHERE sha=?"
	var (
		filename string
		complete bool
		segIndex int
		owner    string
		date     time.Time
	)
	err := db.QueryRow(sqlStr, sha).Scan(&filename, &complete, &segIndex, &owner, &date)

	return err
}
func CheckFileOwner(sha string, owner string) error {
	sqlStr := "SELECT owner FROM sl_files WHERE sha=?"
	var (
		Towner string
	)
	err := db.QueryRow(sqlStr, sha).Scan(&Towner)
	if err != nil && Towner != owner {
		err = errors.New("not your file")
	}
	return err
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

func CompleteFile(sha string) {

}
