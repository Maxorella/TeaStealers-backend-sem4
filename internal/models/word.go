package models

import "html"

type CreateWordData struct {
	Word          string `json:"word"`
	Transcription string `json:"transcription"`
}

func (wd *CreateWordData) Sanitize() {
	wd.Word = html.EscapeString(wd.Word)
	wd.Transcription = html.EscapeString(wd.Transcription)
}

type WordData struct {
	WordID        int    `json:"word_id"`
	Word          string `json:"word"`
	Transcription string `json:"transcription"`
	Link          string `json:"link"`
}

func (wd *WordData) Sanitize() {
	wd.Word = html.EscapeString(wd.Word)
	wd.Transcription = html.EscapeString(wd.Transcription)
	wd.Link = html.EscapeString(wd.Link)
}
