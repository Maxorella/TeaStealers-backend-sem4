DROP TABLE IF EXISTS user_acc;
CREATE TABLE IF NOT EXISTS user_acc(
  user_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  login VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  pass_salt VARCHAR(255) NOT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

DROP TABLE IF EXISTS words_etalon;

CREATE TABLE IF NOT EXISTS words_etalon (
    word_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    word VARCHAR(255) NOT NULL,
    transcription TEXT, -- мб null, а потом добавлять??
    audio_link TEXT, -- мб null, а потом добавлять??
    -- pronouncation VARCHAR(255)
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

DROP TABLE IF EXISTS word_tip;

CREATE TABLE IF NOT EXISTS word_tip(
    tip_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    -- word_id INT, не нужна привязка к слову, подсказки общие
    phonema TEXT NOT NULL, -- мб получится varchar сделать
    tip TEXT NOT NULL
);

-- REFERENCES advert(id)
-- PRIMARY KEY (user_id, advert_id)
DROP TABLE IF EXISTS word_user_try;

CREATE TABLE IF NOT EXISTS word_user_try(
    try_id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    word_id INT REFERENCES words_etalon(word_id) ON DELETE CASCADE, -- каскад опасно, но пока так
    got_transcription TEXT NOT NULL, -- сейчас тут будет слово распознанное, в перспективе транскрипция
    got_result BOOLEAN NOT NULL -- true - без ошибок, false - допущена ошибка
    );

DROP TABLE IF EXISTS user_word_summary;

CREATE TABLE IF NOT EXISTS user_word_summary(
    user_id INT NOT NULL REFERENCES user_acc(user_id),
    word_id INT REFERENCES words_etalon(word_id),
    total_plus INT NOT NULL, -- кол-во правильных произношений
    total_minus INT NOT NULL, -- кол-во неправильных произношений
    PRIMARY KEY (user_id, word_id)
);

