package repo

const (
	CreateUpdateWordStatSql = `INSERT INTO word_progress (word_id, progress) VALUES ($1, $2) ON CONFLICT (word_id) DO UPDATE SET progress = EXCLUDED.progress;`
)
