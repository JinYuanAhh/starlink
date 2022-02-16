package starIM

func Msg_Signin_Success(T string) []byte {
	return []byte(GenerateJson(map[string]string{
		"Type":   "Signin",
		"Status": "Success",
		"T":      T,
	}))
}
func Msg_Signin_Err(err error) []byte {
	return []byte(GenerateJson(map[string]string{
		"Type":   "Signin",
		"Status": "Error",
		"Err":    err.Error(),
	}))
}
func Msg_File_New_Success() []byte {
	return []byte(GenerateJson(map[string]string{
		"Type":   "File",
		"Phrase": "New",
		"Status": "Success",
	}))
}
func Msg_File_New_Err(err error) []byte {
	return []byte(GenerateJson(map[string]string{
		"Type":   "File",
		"Phrase": "New",
		"Status": "Error",
		"Err":    err.Error(),
	}))
}
func Msg_File_Continue_Success() []byte {
	return []byte(GenerateJson(map[string]string{
		"Type":   "File",
		"Phrase": "Continue",
		"Status": "Success",
	}))
}
func Msg_File_Continue_Err(err error) []byte {
	return []byte(GenerateJson(map[string]string{
		"Type":   "File",
		"Phrase": "Continue",
		"Status": "Error",
		"Err":    err.Error(),
	}))
}
