package models

type WordUserStat struct {
	WordId        *int   `json:"id,omitempty"`
	Word          string `json:"word"`
	Transcription string `json:"transcription,omitempty"`
	User          string `json:"user,omitempty"`
	Progress      *int   `json:"progress"`
}
