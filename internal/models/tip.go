package models

import "html"

type TipData struct {
	TipID        *int   `json:"id,omitempty"`
	Phonema      string `json:"phonema"`
	TipText      string `json:"text"`
	TipMediaLink string `json:"media_link"`
	TipAudioLink string `json:"audio_link"`
}

func (td *TipData) Sanitize() {
	td.Phonema = html.EscapeString(td.Phonema)
	td.TipText = html.EscapeString(td.TipText)
	td.TipMediaLink = html.EscapeString(td.TipMediaLink)
	td.TipAudioLink = html.EscapeString(td.TipAudioLink)
}
