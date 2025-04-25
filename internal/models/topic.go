package models

import "html"

type TopicsList struct {
	Topics []OneTopic `json:"topics"`
}

func (tl *TopicsList) Sanitize() {
	for _, topic := range tl.Topics {
		topic.Sanitize()
	}
}

type OneTopic struct {
	TopicId   *int   `json:"topic_id,omitempty"`
	Topic     string `json:"topic"`
	AllWords  *int   `json:"all_words"`
	TrueWords *int   `json:"true_words"`
}

func (tl *OneTopic) Sanitize() {
	tl.Topic = html.EscapeString(tl.Topic)
}
