package repo

const (
	CreateWordSql = `INSERT INTO word_etalon (word, transcription) VALUES ($1, $2) RETURNING word_id`
	SelectWordSql = `SELECT word_id, word, transcription, audio_link from word_etalon WHERE word = $1`
	UploadLinkSql = `UPDATE word_etalon SET audio_link = $1 WHERE word = $2 AND is_deleted = FALSE`
)
