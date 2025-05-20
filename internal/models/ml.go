package models

import "html"

type MlAnswer struct {
	Transcription string `json:"transcription,omitempty"`
	Text          string `json:"text,omitempty"`
	MlError       string `json:"error,omitempty"`
}

func (ml *MlAnswer) Sanitize() {
	ml.Transcription = html.EscapeString(ml.Transcription)
	ml.MlError = html.EscapeString(ml.MlError)
}
