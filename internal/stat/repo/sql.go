package repo

const (
	CreateUpdateWordStatSql = `INSERT INTO word_progress (word_id, word, topic, progress) VALUES ($1, $2, $3, $4) ON CONFLICT (word_id) DO UPDATE SET progress = EXCLUDED.progress;`
	InsertTopic             = `INSERT INTO word_topic (topic) VALUES ($1) ON CONFLICT (topic) DO NOTHING;`
	SelectAllTopics         = `SELECT topic_id, topic FROM word_topic;`
	SelectProgressByTopic   = `SELECT 
    we.topic, 
    COUNT(*) AS total_words,
    COALESCE(SUM(CASE WHEN wp.progress = 1 THEN 1 ELSE 0 END), 0) AS words_with_progress_1
FROM 
    word_etalon we
LEFT JOIN 
    word_progress wp ON we.word_id = wp.word_id
WHERE 
    we.topic = $1 
    AND we.is_deleted = FALSE
GROUP BY 
    we.topic;`
)
