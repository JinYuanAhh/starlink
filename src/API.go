package main

import (
	"github.com/gorilla/mux"
	IM "github.com/starIM"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

func InitMux() {
	R.HandleFunc("/ws/", wsHandler)
	R.HandleFunc("/api/query/{statement}", apiHandler)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte(IM.Query_userPublicInfo(gjson.Get(string(body), "account").String())))
	case "userPrivateInfo":
		w.Write([]byte(IM.Query_userPrivateInfo(gjson.Get(string(body), "T").String())))
	}
}
