package repo

const (
	CreateWordSql              = `INSERT INTO word_etalon (word, transcription, tags) VALUES ($1, $2, $3) RETURNING word_id;`
	SelectWordSql              = `SELECT word_id, word, transcription, audio_id from word_etalon WHERE word = $1;`
	UploadLinkSql              = `UPDATE word_etalon SET audio_id = $1 WHERE word = $2 AND is_deleted = FALSE;`
	GetWordCountSql            = `SELECT COUNT(*) from word_etalon;`
	SelectRandomWordSql        = `SELECT word, transcription, tags, audio_id   FROM word_etalon ORDER BY RANDOM() LIMIT 1;`
	SelectRandomWordWithTagSql = `SELECT word, transcription, tags, audio_id FROM word_etalon WHERE is_deleted = FALSE AND (',' || tags || ',') LIKE '%,' || $1 || ',%' ORDER BY RANDOM() LIMIT 1;`
	InsertTag                  = `INSERT INTO word_tag (tag) VALUES ($1) ON CONFLICT (tag) DO NOTHING;`
	SelectTags                 = `SELECT word_tag FROM word_tag;`
	SelectAllWordsWithTag      = `SELECT word_id, word, transcription, audio_id, tags FROM word_etalon WHERE ',' || tags || ',' LIKE '%,' || $1 || ',%' AND is_deleted = FALSE;`

	Insert1Stat        = `INSERT INTO word_user_try (word_id, result) VALUES ($1, $2);`
	InsertPlusBigStat  = `INSERT INTO user_word_summary (word_id, total_plus, total_minus) VALUES ($1, 1, 0) ON CONFLICT (word_id) DO UPDATE SET total_plus = user_word_summary.total_plus + 1;`
	InsertMinusBigStat = `INSERT INTO user_word_summary (word_id, total_plus, total_minus) VALUES ($1, 0, 1) ON CONFLICT (word_id)  DO UPDATE SET total_minus = user_word_summary.total_minus + 1;`
	SelectBigStat      = `SELECT word_id, total_plus, total_minus FROM user_word_summary WHERE word_id=$1;`
)
