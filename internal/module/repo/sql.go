package repo

const (
	CreateModuleWord   = `INSERT INTO word_modules (title) VALUES ($1) RETURNING id`
	CreateModulePhrase = `INSERT INTO phrase_modules (title) VALUES ($1) RETURNING id`
)
