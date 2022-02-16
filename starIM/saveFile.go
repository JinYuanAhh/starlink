package starIM

import (
	"crypto"
	"encoding/hex"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

//func init() {
//
//	fmt.Println(StartFileSave("test.txt", "3", "202cb962ac59075b964b07152d234b70"))
//	fmt.Println(ContinueFileSave("test.txt", []byte("1")))
//	fmt.Println(ContinueFileSave("test.txt", []byte("2")))
//	fmt.Println(ContinueFileSave("test.txt", []byte("3")))
//
//}

func StartFileSave(fn string, CompSegIndex string, MD5 string) error {
	_, err := os.Create(StrConnect(filepath.Dir(os.Args[0]), "/Resources/Files/", fn))
	if err != nil {
		return err
	}
	f, err := os.Create(StrConnect(filepath.Dir(os.Args[0]), "/Resources/FilesInfo/", fn, ".i"))
	if err != nil {
		return err
	} else {
		defer f.Close()
		_, err := f.Write([]byte(GenerateJson(map[string]string{
			"MD5":          MD5,
			"CompSegIndex": CompSegIndex,
			"SegIndex":     "0",
		})))
		if err != nil {
			return err
		}
		return nil
	}

}
func ContinueFileSave(fn string, segment []byte) (bool, error) {
	f, err := os.OpenFile(StrConnect(filepath.Dir(os.Args[0]), "/Resources/Files/", fn), os.O_APPEND|os.O_RDONLY, 0666)
	if err != nil {
		return false, err
	} else {
		defer f.Close()
		fi, err := GetFileInfo(fn)
		if err != nil {
			return false, err
		}
		segmentIndex := gjson.Get(fi, "SegIndex").Int() + 1
		CompleteSegmentIndex := gjson.Get(fi, "CompSegIndex").Int()
		if CompleteSegmentIndex < int64(segmentIndex) {
			return false, errors.New("seg err(too big)")
		} else {
			_, err := f.Write(segment)
			if err != nil {
				return false, err
			}
			err = UpdateFileSegIndexNow(fn, int(segmentIndex))
			if err != nil {
				return false, err
			}
		}
		if CompleteSegmentIndex == int64(segmentIndex) {
			_, err := EndFileSave(fn)
			if err != nil {
				return true, err
			}
		}
		return false, nil
	}
} //返回 是否为最后一段, Err
func EndFileSave(fn string) (bool, error) {
	fi, err := GetFileInfo(fn)
	if err != nil {
		return false, err
	}
	MD5_Origin := gjson.Get(fi, "MD5").String()
	MD5, err := GetFileMD5(StrConnect(filepath.Dir(os.Args[0]), "/Resources/Files/", fn))
	if err != nil {
		return false, err
	}
	if MD5 == MD5_Origin {
		err := os.Remove(StrConnect(filepath.Dir(os.Args[0]), "/Resources/FilesInfo/", fn, ".i"))
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, errors.New("md5 difference")
	}

}
func GetFileInfo(fn string) (string, error) {
	fi, err := ioutil.ReadFile(StrConnect(filepath.Dir(os.Args[0]), "/Resources/FilesInfo/", fn, ".i"))
	if err != nil {
		return "", err
	} else {
		return string(fi), nil
	}
}
func UpdateFileSegIndexNow(fn string, segIndex int) error { //在文件中设置好现在传输了几段了
	fi, err := os.OpenFile(StrConnect(filepath.Dir(os.Args[0]), "/Resources/FilesInfo/", fn, ".i"), os.O_RDWR|os.O_RDONLY, 0666)
	if err != nil {
		return err
	} else {
		defer fi.Close()
		fi_o, _ := GetFileInfo(fn)
		fi_o, _ = sjson.Set(fi_o, "SegIndex", segIndex)
		_, err := fi.Write([]byte(fi_o))
		if err != nil {
			return err
		}
		return nil
	}
}
func GetFileMD5(fn string) (string, error) {
	MD5 := crypto.MD5.New()
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(MD5, f)
	if err != nil {
		return "", err
	}
	MD5str := hex.EncodeToString(MD5.Sum(nil))
	return MD5str, nil
}
