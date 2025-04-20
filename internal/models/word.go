package models

import "html"

type CreateWordData struct {
	WordID        *int   `json:"id,omitempty"`
	Word          string `json:"word"`
	Transcription string `json:"transcription"`
	AudioLink     string `json:"audio_link"`
	Topic         string `json:"topic"`
}

func (wd *CreateWordData) Sanitize() {
	wd.Word = html.EscapeString(wd.Word)
	wd.Transcription = html.EscapeString(wd.Transcription)
	wd.Topic = html.EscapeString(wd.Topic)
}

type WordData struct {
	WordID        *int   `json:"id,omitempty"`
	Word          string `json:"word"`
	Topic         string `json:"topic"`
	Transcription string `json:"transcription"`
	AudioLink     string `json:"audio_link"`
	Progress      *int   `json:"progress,omitempty"`
}

func (wd *WordData) Sanitize() {
	wd.Word = html.EscapeString(wd.Word)
	wd.Topic = html.EscapeString(wd.Topic)
	wd.Transcription = html.EscapeString(wd.Transcription)
	wd.AudioLink = html.EscapeString(wd.AudioLink)
}
