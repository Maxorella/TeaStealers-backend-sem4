package repo

const (
	GetIncompletePhraseModuleSql = `
        SELECT m.id
        FROM phrase_modules m
        JOIN phrase_exercises e ON e.module_id = m.id
        LEFT JOIN exercise_progress p 
            ON p.exercise_id = e.id AND p.exercise_type = 'phrase' AND p.user_id = $1
        GROUP BY m.id
        HAVING COUNT(*) FILTER (WHERE p.status = 'completed') < COUNT(*)
        ORDER BY m.id
        LIMIT 1
    `

	GetIncompleteWordModuleSql = `
        SELECT m.id
        FROM word_modules m
        JOIN word_exercises e ON e.module_id = m.id
        LEFT JOIN exercise_progress p 
            ON p.exercise_id = e.id AND p.exercise_type = 'word' AND p.user_id = $1
        GROUP BY m.id
        HAVING COUNT(*) FILTER (WHERE p.status = 'completed') < COUNT(*)
        ORDER BY m.id
        LIMIT 1
    `

	GetWordModuleExercisesWithProgressSql = `
        SELECT e.id, e.exercise_type, e.words, e.transcriptions, e.audio, e.translations, e.module_id,
               COALESCE(p.status, 'none') AS status
        FROM word_exercises e
        LEFT JOIN exercise_progress p 
            ON p.exercise_id = e.id AND p.exercise_type = 'word' AND p.user_id = $1
        WHERE e.module_id = $2
        ORDER BY e.id
    `

	GetWordModuleExercisesSql = `
        SELECT id, exercise_type, words, transcriptions, audio, translations, module_id, 'none' AS status
        FROM word_exercises
        WHERE module_id = $1
        ORDER BY id
    `

	GetPhraseModuleExercisesWithProgressSql = `
        SELECT e.id, e.exercise_type, e.sentence, e.translate, e.transcription, e.audio, e.chain, e.module_id,
               COALESCE(p.status, 'none') AS status
        FROM phrase_exercises e
        LEFT JOIN exercise_progress p 
            ON p.exercise_id = e.id AND p.exercise_type = 'phrase' AND p.user_id = $1
        WHERE e.module_id = $2
        ORDER BY e.id
    `

	GetPhraseModuleExercisesSql = `
        SELECT id, exercise_type, sentence, translate, transcription, audio, chain, module_id, 'none' AS status
        FROM phrase_exercises
        WHERE module_id = $1
        ORDER BY id
    `

	// node sql
	SelectPhraseModulesSql = `
        SELECT id, title 
        FROM phrase_modules 
        ORDER BY id
    `

	SelectWordModulesSql = `
        SELECT id, title 
        FROM word_modules 
        ORDER BY id
    `

	CreateWordExerciseSql = `
INSERT INTO word_exercises (
    exercise_type,
    words,
    transcriptions,
    audio,
    translations,
    module_id
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
`
	CreatePhraseExerciseSql = `
INSERT INTO phrase_exercises (
    exercise_type,
    sentence,
    translate,
    transcription,
    audio,
    chain,
    module_id
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;
`
	UpsertExerciseProgressSql = `
        INSERT INTO exercise_progress (user_id, exercise_id, exercise_type, status)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, exercise_id, exercise_type)
        DO UPDATE SET status = EXCLUDED.status, updated_at = CURRENT_TIMESTAMP
        RETURNING id;
    `
	// new sql
	SelectWordSql                 = `SELECT word_id, word, transcription, audio_link, topic from word_etalon WHERE word = $1 AND is_deleted = FALSE;`
	CreateWordSql                 = `INSERT INTO word_etalon (word, transcription, audio_link, topic) VALUES ($1, $2, $3, $4) RETURNING word_id;`
	InsertWordTip                 = `INSERT INTO word_tip (phonema, tip_text, tip_audio_link, tip_video_link) VALUES ($1, $2, $3, $4);`
	SelectWordTip                 = `SELECT phonema, tip_text, tip_audio_link, tip_video_link from word_tip WHERE phonema = $1;`
	SelectWordWithProgressByTopic = `SELECT 
    we.word,
    we.transcription,
    we.audio_link,
    we.topic,
    COALESCE(
        (SELECT wp.progress 
         FROM word_progress wp 
         WHERE wp.word_id = we.word_id), 
    0) AS progress
FROM 
    word_etalon we
WHERE 
    we.topic = $1`
	SelectRandomWordSql          = `SELECT word_id, word, transcription, topic, audio_link   FROM word_etalon ORDER BY RANDOM() LIMIT 1;`
	SelectRandomWordWithTopicSql = `SELECT word_id, word, transcription, topic, audio_link   FROM word_etalon WHERE is_deleted = FALSE AND topic=$1 ORDER BY RANDOM() LIMIT 1;`

	//old sql
	UploadLinkSql         = `UPDATE word_etalon SET audio_id = $1 WHERE word = $2 AND is_deleted = FALSE;`
	GetWordCountSql       = `SELECT COUNT(*) from word_etalon;`
	InsertTag             = `INSERT INTO word_tag (tag) VALUES ($1) ON CONFLICT (tag) DO NOTHING;`
	SelectTags            = `SELECT word_tag FROM word_tag;`
	SelectAllWordsWithTag = `SELECT word_id, word, transcription, audio_id, tags FROM word_etalon WHERE tags=$1 AND is_deleted = FALSE;`

	Insert1Stat        = `INSERT INTO word_user_try (word_id, result) VALUES ($1, $2);`
	InsertPlusBigStat  = `INSERT INTO user_word_summary (word_id, total_plus, total_minus) VALUES ($1, 1, 0) ON CONFLICT (word_id) DO UPDATE SET total_plus = user_word_summary.total_plus + 1;`
	InsertMinusBigStat = `INSERT INTO user_word_summary (word_id, total_plus, total_minus) VALUES ($1, 0, 1) ON CONFLICT (word_id)  DO UPDATE SET total_minus = user_word_summary.total_minus + 1;`
	SelectBigStat      = `SELECT word_id, total_plus, total_minus FROM user_word_summary WHERE word_id=$1;`

	CreateWordTip = `INSERT INTO word_tip (phonema, tip_text, tip_picture, tip_audio) VALUES ($1, $2,$3,$4);`
)
