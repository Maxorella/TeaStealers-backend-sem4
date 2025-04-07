DROP TABLE IF EXISTS user_acc;
CREATE TABLE IF NOT EXISTS user_acc(
  user_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  login VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  pass_salt VARCHAR(255) NOT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

DROP TABLE IF EXISTS word_etalon;

CREATE TABLE IF NOT EXISTS word_etalon (
    word_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    word VARCHAR(255) NOT NULL,
    transcription TEXT,
    audio_id TEXT,
    tags TEXT,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

DROP TABLE IF EXISTS word_tip;

CREATE TABLE IF NOT EXISTS word_tip(
    tip_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    -- word_id INT, не нужна привязка к слову, подсказки общие
    phonema TEXT NOT NULL,
    tip TEXT NOT NULL
);

DROP TABLE IF EXISTS word_user_try;

CREATE TABLE IF NOT EXISTS word_user_try(
    try_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    word_id INT REFERENCES word_etalon(word_id) ON DELETE CASCADE,
    got_transcription TEXT NOT NULL,
    result BOOLEAN NOT NULL
    );

DROP TABLE IF EXISTS user_word_summary;

CREATE TABLE IF NOT EXISTS user_word_summary(
    user_id INT NOT NULL REFERENCES user_acc(user_id),
    word_id INT REFERENCES word_etalon(word_id),
    total_plus INT NOT NULL, -- кол-во правильных произношений
    total_minus INT NOT NULL, -- кол-во неправильных произношений
    PRIMARY KEY (user_id, word_id)
);

