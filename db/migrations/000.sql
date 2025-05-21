DROP TABLE IF EXISTS exercise_progress;
DROP TABLE IF EXISTS phrase_exercises;
DROP TABLE IF EXISTS word_exercises;
DROP TABLE IF EXISTS phrase_modules;
DROP TABLE IF EXISTS word_modules;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS phrase_exercise_type;
DROP TYPE IF EXISTS word_exercise_type;
DROP TABLE IF EXISTS word_tip;

CREATE TYPE word_exercise_type AS ENUM (
    'pronounce',
    'guessWord',
    'pronounceFiew'
);

CREATE TYPE phrase_exercise_type AS ENUM (
    'pronounce',
    'completeChain'
);

CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY,
    passwordHash TEXT CONSTRAINT passwordHash_length CHECK ( char_length(passwordHash) <= 40) NOT NULL,
    levelUpdate INTEGER NOT NULL DEFAULT 1,
    name VARCHAR(50) NOT NULL,
    email TEXT NOT NULL UNIQUE,
    dateCreation TIMESTAMP NOT NULL DEFAULT NOW(),
    isDeleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE word_modules (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
);

CREATE TABLE phrase_modules (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
);

CREATE TABLE word_exercises (
    id SERIAL PRIMARY KEY,
    exercise_type word_exercise_type NOT NULL,
    words TEXT[] NOT NULL,
    transcriptions TEXT[] NOT NULL,
    audio TEXT[] NOT NULL,
    translations TEXT[] NOT NULL,
    module_id INTEGER REFERENCES word_modules(id) ON DELETE CASCADE
);

CREATE TABLE phrase_exercises (
    id SERIAL PRIMARY KEY,
    exercise_type phrase_exercise_type NOT NULL,
    sentence TEXT,
    translate TEXT,
    transcription TEXT,
    audio TEXT NOT NULL,
    chain TEXT[],
    module_id INTEGER REFERENCES phrase_modules(id) ON DELETE CASCADE
);

CREATE TABLE exercise_progress (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id INTEGER NOT NULL,
    exercise_type VARCHAR(10) NOT NULL,  -- "word" или "phrase"
    status VARCHAR(20) NOT NULL DEFAULT 'none',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_progress_per_exercise UNIQUE (user_id, exercise_id, exercise_type)
);

CREATE TABLE IF NOT EXISTS word_tip(
    phonema TEXT PRIMARY KEY,
    tip_text TEXT,
    tip_audio_link TEXT,
    tip_video_link TEXT
);

