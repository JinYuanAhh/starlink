package main

import (
	"github.com/gorilla/mux"
	IM "github.com/starIM"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
)

func InitMux() {
	R.HandleFunc("/ws", wsHandler)
	R.HandleFunc("/api/query/{statement}", apiQueryHandler).Methods("POST")
}

func apiQueryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statement := vars["statement"]
	body, _ := ioutil.ReadAll(r.Body)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	//w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Add("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	//w.Header().Add("Access-Control-Allow-Credentials", "true")
	switch statement {
	case "userPublicInfo":
		w.Write([]byte(IM.Query_userPublicInfo(gjson.GetBytes(body, "Account").String())))
	case "userPrivateInfo":
		w.Write([]byte(IM.Query_userPrivateInfo(gjson.GetBytes(body, "T").String())))
	case "historyPublicMessages":
		rs := gjson.GetManyBytes(body, "PS", "PN") //PageSize PageNumber
		w.Write([]byte(IM.GetHistroyPublicMessages(int(rs[0].Int()), int(rs[1].Int()))))
	case "file":
		sha := gjson.GetBytes(body, "Sha").String()
		f, err := IM.QueryFile(sha)
		if err != nil {
			w.Write([]byte(err.Error()))

		} else {
			w.Write(IM.GenerateJson(map[string]string{
				"Filename": f.Filename,
				"Owner":    f.Owner,
				"Args":     f.Args,
				"Date":     f.Date.String(),
				"Complete": strconv.FormatBool(f.Complete),
				"Size":     f.Size,
			}))
		}
	}
}
