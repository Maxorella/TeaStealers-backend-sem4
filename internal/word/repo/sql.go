package repo

const (
	// new sql
	SelectWordSql = `SELECT word_id, word, transcription, audio_link, topic from word_etalon WHERE word = $1 AND is_deleted = FALSE;`
	CreateWordSql = `INSERT INTO word_etalon (word, transcription, audio_link, topic) VALUES ($1, $2, $3, $4) RETURNING word_id;`
	InsertTopic   = `INSERT INTO word_topic (topic) VALUES ($1) ON CONFLICT (topic) DO NOTHING;`

	//old sql
	UploadLinkSql              = `UPDATE word_etalon SET audio_id = $1 WHERE word = $2 AND is_deleted = FALSE;`
	GetWordCountSql            = `SELECT COUNT(*) from word_etalon;`
	SelectRandomWordSql        = `SELECT word_id, word, transcription, tags, audio_id   FROM word_etalon ORDER BY RANDOM() LIMIT 1;`
	SelectRandomWordWithTagSql = `SELECT word_id, word, transcription, tags, audio_id FROM word_etalon WHERE is_deleted = FALSE AND tags=$1 ORDER BY RANDOM() LIMIT 1;`
	InsertTag                  = `INSERT INTO word_tag (tag) VALUES ($1) ON CONFLICT (tag) DO NOTHING;`
	SelectTags                 = `SELECT word_tag FROM word_tag;`
	SelectAllWordsWithTag      = `SELECT word_id, word, transcription, audio_id, tags FROM word_etalon WHERE tags=$1 AND is_deleted = FALSE;`

	Insert1Stat        = `INSERT INTO word_user_try (word_id, result) VALUES ($1, $2);`
	InsertPlusBigStat  = `INSERT INTO user_word_summary (word_id, total_plus, total_minus) VALUES ($1, 1, 0) ON CONFLICT (word_id) DO UPDATE SET total_plus = user_word_summary.total_plus + 1;`
	InsertMinusBigStat = `INSERT INTO user_word_summary (word_id, total_plus, total_minus) VALUES ($1, 0, 1) ON CONFLICT (word_id)  DO UPDATE SET total_minus = user_word_summary.total_minus + 1;`
	SelectBigStat      = `SELECT word_id, total_plus, total_minus FROM user_word_summary WHERE word_id=$1;`

	CreateWordTip = `INSERT INTO word_tip (phonema, tip_text, tip_picture, tip_audio) VALUES ($1, $2,$3,$4);`
	SelectWordTip = `SELECT phonema, tip_text, tip_picture, tip_audio from word_tip WHERE phonema = $1;`
)
