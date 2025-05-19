package models

import "html"

type CreateWordData struct {
	Exercise      string `json:"exercise"`
	Word          string `json:"word"`
	ModuleId      *int   `json:"id"`
	Transcription string `json:"transcription"`
	AudioLink     string `json:"audio_link"`
	Translation   string `json:"translation"`
}

type CreateWordDataList struct {
	Exercise      string   `json:"exercise"`
	Word          []string `json:"word"`
	ModuleId      *int     `json:"id"`
	Transcription []string `json:"transcription"`
	AudioLink     []string `json:"audio_link"`
	Translation   []string `json:"translation"`
}

type CreatePhraseData struct {
	Exercise      string   `json:"exercise"`
	Sentence      string   `json:"word"`
	Transcription string   `json:"transcription"`
	ModuleId      *int     `json:"id"`
	AudioLink     string   `json:"audio"`
	Translate     string   `json:"translate"`
	Chain         []string `json:"chain"`
}

type ExerciseProgress struct {
	UserID       *int   `json:"user_id"`
	ExerciseID   *int   `json:"exercise_id"`
	ExerciseType string `json:"exercise_type"`
	Status       string `json:"status"`
}

type IdStruct struct {
	Id *int `json:"id"`
}

func (wd *CreateWordData) Sanitize() {
	wd.Word = html.EscapeString(wd.Word)
	wd.Transcription = html.EscapeString(wd.Transcription)
	wd.Translation = html.EscapeString(wd.Translation)
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
