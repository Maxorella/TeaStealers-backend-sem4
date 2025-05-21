package models

type TopicsList struct {
	Topics []OneTopic `json:"topics"`
}

type OneTopic struct {
	TopicId   *int   `json:"topic_id,omitempty"`
	Topic     string `json:"topic"`
	AllWords  *int   `json:"all_words"`
	TrueWords *int   `json:"true_words"`
}
