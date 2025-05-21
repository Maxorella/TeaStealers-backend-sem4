package models

type MlAnswer struct {
	Transcription string `json:"transcription,omitempty"`
	Text          string `json:"text,omitempty"`
	MlError       string `json:"error,omitempty"`
}
