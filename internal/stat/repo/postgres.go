package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
	"strconv"
)

type StatRepo struct {
	db     *sql.DB
	logger logger.Logger
	// cfg    *config.Config
	// metricsC metrics.MetricsHTTP
}

func NewRepository(db *sql.DB, logger logger.Logger) *StatRepo {
	return &StatRepo{db: db, logger: logger}
}

func (r *StatRepo) BeginTx(ctx context.Context) (models.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *StatRepo) UpdateWordStat(ctx context.Context, tx models.Transaction, word *models.WordData, result int) error {
	requestId := utils2.GetRequestIDFromCtx(ctx)
	wordIdStr := strconv.Itoa(*word.WordID)

	res, err := tx.ExecContext(ctx, CreateUpdateWordStatSql, wordIdStr, word.Word, word.Topic, result)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "UpdateWordStat", err)
		return fmt.Errorf("failed to update word stat: %w", err)
	}

	// Проверяем количество затронутых строк
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "UpdateWordStat", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		errMsg := "no rows were affected by the update"
		r.logger.LogError(requestId, logger.RepositoryLayer, "UpdateWordStat", errors.New(errMsg))
		return errors.New(errMsg)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "UpdateWordStat",
		fmt.Sprintf("successfully updated word stat, word_id: %s, rows affected: %d", wordIdStr, rowsAffected))
	return nil
}

func (r *StatRepo) InsertTopic(ctx context.Context, tx models.Transaction, wordTopic string) error {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	_, err := tx.Exec(InsertTopic, wordTopic)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "InsertTopic", err)
		return fmt.Errorf("failed to insert topic: %w", err)
	}
	r.logger.LogInfo(requestId, logger.RepositoryLayer, "InsertTopic", "topic inserted successfully: "+wordTopic)
	return nil
}

func (r *StatRepo) SelectAllTopics(ctx context.Context, tx models.Transaction) (*models.TopicsList, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)
	res, err := r.db.QueryContext(ctx, SelectAllTopics)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "SelectAllTopics", err)
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer res.Close()

	topics := models.TopicsList{} // TODO тут мб по другому надо инициализировать
	for res.Next() {
		topic := models.OneTopic{new(int), "", new(int), new(int)}
		if err := res.Scan(topic.TopicId, &topic.Topic); err != nil {
			r.logger.LogError(requestId, logger.RepositoryLayer, "SelectAllTopics", err)
			return nil, err
		}
		topic.TopicId = nil
		topics.Topics = append(topics.Topics, topic)
	}

	if err := res.Err(); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "SelectAllTopics", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	for _, topic := range topics.Topics {
		newtopic, err := r.GetTopicProgress(ctx, tx, topic.Topic)
		if err != nil {
			return nil, err
		}
		topic = *newtopic // todo тут может неправильно быть??
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "SelectAllTopics",
		fmt.Sprintf("successfully got topics"))
	return &topics, nil
}

/*
func (r *StatRepo) GetTopicProgress(ctx context.Context, tx models.Transaction, topic string) (*models.OneTopic, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	res := r.db.QueryRow(SelectProgressByTopic, topic)

	gotTopic := &models.OneTopic{new(int), topic, new(int), new(int)}
	// gotTopic.Topic = topic
	var all_words, true_words sql.NullInt32
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Scan(&all_words, &true_words); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetTopicProgress", err)
		return &models.OneTopic{}, err
	}

	if all_words.Valid {
		*gotTopic.AllWords = int(all_words.Int32)
	} else {
		return &models.OneTopic{}, errors.New("all_words not valid")
	}
	if true_words.Valid {
		*gotTopic.TrueWords = int(true_words.Int32)
	} else {
		return &models.OneTopic{}, errors.New("true_words not valid")

	}

	return gotTopic, nil
}

*/

func (r *StatRepo) GetTopicProgress(ctx context.Context, tx models.Transaction, topic string) (*models.OneTopic, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	res := r.db.QueryRow(SelectProgressByTopic, topic)

	gotTopic := &models.OneTopic{
		Topic:     topic,
		TrueWords: new(int),
		AllWords:  new(int),
	}

	var all_words, true_words sql.NullInt32
	err := res.Scan(&gotTopic.Topic, &all_words, &true_words)

	// Обрабатываем случай, когда нет строк
	if err == sql.ErrNoRows {
		//	*gotTopic.AllWords = -1
		//	*gotTopic.TrueWords = -1
		return gotTopic, errors.New("no rows get topic")
	}

	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetTopicProgress", err)
		return nil, err
	}

	if all_words.Valid {
		*gotTopic.AllWords = int(all_words.Int32)
	} else {
		*gotTopic.AllWords = -1
	}

	if true_words.Valid {
		*gotTopic.TrueWords = int(true_words.Int32)
	} else {
		*gotTopic.TrueWords = -1
	}

	return gotTopic, nil
}

/*
func (r *StatRepo) GetTopicStat(ctx context.Context, tx models.Transaction, topic string) error {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	res, err := r.db.QueryRow(SelectTopicProgres, topic)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "UpdateWordStat", err)
		return fmt.Errorf("failed to update word stat: %w", err)
	}

	// Проверяем количество затронутых строк
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "UpdateWordStat", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		errMsg := "no rows were affected by the update"
		r.logger.LogError(requestId, logger.RepositoryLayer, "UpdateWordStat", errors.New(errMsg))
		return errors.New(errMsg)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "UpdateWordStat",
		fmt.Sprintf("successfully updated word stat, word_id: %s, rows affected: %d", wordIdStr, rowsAffected))
	return nil
}

*/
