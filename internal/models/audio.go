package models

// UserSignUpData represents user information for signup
type AudioFileData struct {
	FileName string `json:"file_name"`
	Word     string `json:"word"`
	Link     string `json:"link"`
}

/*
func (usSignUp *UserSignUpData) Sanitize() {
	usSignUp.Email = html.EscapeString(usSignUp.Email)
	usSignUp.Phone = html.EscapeString(usSignUp.Phone)
	usSignUp.Password = html.EscapeString(usSignUp.Password)
}
*/
