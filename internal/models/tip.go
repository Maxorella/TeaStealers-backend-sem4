package models

type TipData struct {
	TipID        *int   `json:"id,omitempty"`
	Phonema      string `json:"phonema"`
	TipText      string `json:"text"`
	TipMediaLink string `json:"media_link"`
	TipAudioLink string `json:"audio_link"`
}
