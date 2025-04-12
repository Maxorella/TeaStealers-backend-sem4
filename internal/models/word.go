package models

import "html"

type CreateWordData struct {
	Word          string `json:"word"`
	Transcription string `json:"transcription"`
	Tags          string `json:"tags"`
}

func (wd *CreateWordData) Sanitize() {
	wd.Word = html.EscapeString(wd.Word)
	wd.Transcription = html.EscapeString(wd.Transcription)
}

type TagsList struct {
	Tags []string `json:"tags"`
}
type OneTag struct {
	Tag string `json:"tag"`
}

type WordStat struct {
	Id         int `json:"id"`
	TotalPlus  int `json:"plus"`
	TotalMinus int `json:"minus"`
}

type WordData struct {
	WordID        int    `json:"word_id"`
	Word          string `json:"word"`
	Tags          string `json:"tags"`
	Transcription string `json:"transcription"`
	Link          string `json:"link"`
}

func (wd *WordData) Sanitize() {
	wd.Word = html.EscapeString(wd.Word)
	wd.Tags = html.EscapeString(wd.Tags)
	wd.Transcription = html.EscapeString(wd.Transcription)
	wd.Link = html.EscapeString(wd.Link)
}

type TipData struct {
	TipID      int    `json:"tipID,omitempty"`
	Phonema    string `json:"phonema"`
	TipText    string `json:"tipText"`
	TipPicture string `json:"tipPicture"`
	TipAudio   string `json:"tipAudio"`
}

func (td *TipData) Sanitize() {
	td.Phonema = html.EscapeString(td.Phonema)
	td.TipText = html.EscapeString(td.TipText)
	td.TipPicture = html.EscapeString(td.TipPicture)
	td.TipAudio = html.EscapeString(td.TipAudio)
}
