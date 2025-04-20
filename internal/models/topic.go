package models

import "html"

type TopicsList struct {
	Topics []string `json:"topics"`
}

func (tl *TopicsList) Sanitize() {
	for i, topic := range tl.Topics {
		tl.Topics[i] = html.EscapeString(topic)
	}
}

type OneTopic struct {
	Topic string `json:"topic"`
}
