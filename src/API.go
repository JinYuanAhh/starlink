package main

func InitMux() {
	R.HandleFunc("/ws", wsHandler)
}

//TODO ↓ httpQueryHandler
//func apiQueryHandler(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	statement := vars["statement"]
//	body, _ := ioutil.ReadAll(r.Body)
//	switch statement {
//	case "userPublicInfo":
//		w.Write([]byte(IM.Query_userPublicInfo(gjson.GetBytes(body, "Account").String())))
//	case "userPrivateInfo":
//		w.Write([]byte(IM.Query_userPrivateInfo(gjson.GetBytes(body, "T").String())))
//	case "historyPublicMessages":
//		rs := gjson.GetManyBytes(body, "PS", "PN") //PageSize PageNumber
//		w.Write([]byte(IM.GetHistroyPublicMessages(int(rs[0].Int()), int(rs[1].Int()))))
//	case "fileInfo":
//		sha := gjson.GetBytes(body, "Sha").String()
//		f, err := IM.QueryFile(sha)
//		if err != nil {
//			w.Write([]byte(err.Error()))
//		} else {
//			C := "false"
//			if f.Complete[0] == '\u0001' {
//				C = "true"
//			}
//			w.Write(IM.GenerateJson(map[string]string{
//				"Filename": f.Filename,
//				"Owner":    f.Owner,
//				"Args":     f.Args,
//				"Date":     f.Date.String(),
//				"Complete": C,
//				"Size":     fmt.Sprintf("%f", f.Size),
//			}))
//		}
//	case "file":
//		if r.Body == http.NoBody {
//			return
//		}
//		rs := gjson.GetManyBytes(body, "Sha", "Buffsize", "SegIndex")
//		f, err := os.OpenFile(IM.StrConnect("./Resources/Files/", rs[0].String()), os.O_RDONLY, 0777)
//		defer f.Close()
//		if err != nil {
//			w.Write([]byte(err.Error()))
//		} else {
//			buf := make([]byte, rs[1].Int())
//			w.Header().Set("Content-Disposition", "attachment; filename=")
//			w.Header().Set("Content-Type", "application/octet-stream")
//			w.Header().Set("Content-Length", strconv.FormatInt(rs[1].Int(), 10))
//			readSize, err := f.ReadAt(buf, rs[1].Int()*(rs[2].Int()-1))
//			if err != nil && err != io.EOF {
//				w.Write([]byte(err.Error())) //INFO:: Buffsize 单位是B
//			} else {
//				if int64(readSize) < rs[1].Int() {
//
//				}
//				w.Write(buf) //INFO:: Buffsize 单位是B
//				//for {
//				//	readSize, err := f.Read(buf)
//				//	if err != nil && err != io.EOF {
//				//		w.WriteHeader(http.StatusInternalServerError)
//				//		return
//				//	} else if err == io.EOF {
//				//		return
//				//	}
//				//}
//			}
//			return
//		}
//	}
//} // //dd
//↑ httpQueryHandler
