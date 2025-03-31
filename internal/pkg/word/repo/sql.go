package repo

const (
	CreateWordSql = `INSERT INTO word_etalon (word, transcription) VALUES ($1, $2) RETURNING word_id`
	SelectWordSql = `SELECT word_id, word, transcription, audio_link from word_etalon WHERE word = $1`
)
