package models

import "html"

type MlAnswer struct {
	Transcription string `json:"transcription"`
	MlError       string `json:"error"`
}

func (ml *MlAnswer) Sanitize() {
	ml.Transcription = html.EscapeString(ml.Transcription)
	ml.MlError = html.EscapeString(ml.MlError)
}
