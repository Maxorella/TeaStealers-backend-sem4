package repo

const (
	CreateWordSql              = `INSERT INTO word_etalon (word, transcription) VALUES ($1, $2) RETURNING word_id;`
	SelectWordSql              = `SELECT word_id, word, transcription, audio_id from word_etalon WHERE word = $1;`
	UploadLinkSql              = `UPDATE word_etalon SET audio_id = $1 WHERE word = $2 AND is_deleted = FALSE;`
	GetWordCountSql            = `SELECT COUNT(*) from word_etalon;`
	SelectRandomWordSql        = `SELECT word, transcription, tags, audio_id   FROM word_etalon ORDER BY RANDOM() LIMIT 1;`
	SelectRandomWordWithTagSql = `SELECT word, transcription, tags, audio_id FROM word_etalon WHERE is_deleted = FALSE AND (',' || tags || ',') LIKE '%,' || $1 || ',%' ORDER BY RANDOM() LIMIT 1;`
)
