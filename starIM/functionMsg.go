package starIM

import "fmt"

func Msg_Signin_Success(T string) []byte {
	m := `{"Type":"Signin","Status":"Success","T":"%s"}`
	return []byte(fmt.Sprintf(m, T))
}
func Msg_Signin_Logged() []byte {
	m := `{"Type":"Signin","Status":"Logged"}`
	return []byte(m)
}
func Msg_File_New_Success() []byte {
	m := `{"Type":"File","Phrase":"New","Status":"Success"}`
	return []byte(m)
}
func Msg_File_New_Err(err error) []byte {
	m := `{"Type":"File","Phrase":"New","Status":"Error","Err":"%s"}`
	return []byte(fmt.Sprintf(m, err))
}
func Msg_File_Continue_Success() []byte {
	m := `{"Type":"File","Phrase":"Continue","Status":"Success"}`
	return []byte(m)
}
func Msg_File_Continue_Err(err error) []byte {
	m := `{"Type":"File","Phrase":"Continue","Status":"Error","Err":"%s"}`
	return []byte(fmt.Sprintf(m, err))
}