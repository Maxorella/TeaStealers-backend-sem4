-- TODO add user
--DROP TABLE IF EXISTS user_acc;
--CREATE TABLE IF NOT EXISTS user_acc(
--  user_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--  login VARCHAR(255) NOT NULL,
--  password VARCHAR(255) NOT NULL,
--  pass_salt VARCHAR(255) NOT NULL,
--  is_deleted BOOLEAN NOT NULL DEFAULT FALSE
--);


DROP TABLE  IF EXISTS word_etalon CASCADE;
CREATE TABLE IF NOT EXISTS word_etalon (
    word_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    word VARCHAR(255) NOT NULL,
    transcription TEXT NOT NULL,
    audio_link TEXT,
    topic TEXT NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

DROP TABLE IF EXISTS word_topic;
CREATE TABLE IF NOT EXISTS word_topic(
    topic_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    topic TEXT NOT NULL UNIQUE
);

DROP TABLE IF EXISTS word_progress;
CREATE TABLE IF NOT EXISTS word_progress(
    word_id INT PRIMARY KEY REFERENCES word_etalon(word_id),
    -- user_id INT TODO add user
    word TEXT,
    topic TEXT,
    progress INT
);

--DROP TABLE IF EXISTS topic_progress;
--CREATE TABLE IF NOT EXISTS topic_progress(
    -- user_id INT TODO add user double primary key
--    topic TEXT  PRIMARY KEY,
--    true_words INT,
--   all_words INT
--);


DROP TABLE IF EXISTS word_tip;
CREATE TABLE IF NOT EXISTS word_tip(
    phonema TEXT PRIMARY KEY,
    tip_text TEXT ,
    tip_audio_link TEXT,
    tip_video_link TEXT
);


-- DROP TABLE IF EXISTS word_user_try;
-- CREATE TABLE IF NOT EXISTS word_user_try(
--    try_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--    word_id INT REFERENCES word_etalon(word_id) ON DELETE CASCADE,
--    -- got_transcription TEXT NOT NULL,
--    result BOOLEAN NOT NULL
--    );

-- DROP TABLE IF EXISTS user_word_summary;
-- CREATE TABLE IF NOT EXISTS user_word_summary(
--    -- user_id INT NOT NULL REFERENCES user_acc(user_id),
--    word_id INT REFERENCES word_etalon(word_id),
--    total_plus INT NOT NULL, -- кол-во правильных произношений
--    total_minus INT NOT NULL, -- кол-во неправильных произношений
--    PRIMARY KEY (word_id)
-- );

