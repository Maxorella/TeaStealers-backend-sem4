package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/pkg/config"
	"github.com/TeaStealers-backend-sem4/internal/pkg/minio/helpers"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

// MinClient интерфейс для взаимодействия с Minio
type MinClient interface {
	InitMinio() error                                    // Метод для инициализации подключения к Minio
	CreateOne(file helpers.FileDataType) (string, error) // Метод для создания одного объекта в бакете Minio
	GetOne(objectID string) (string, error)              // Метод для получения одного объекта из бакета Minio
	DeleteOne(objectID string) error                     // Метод для удаления одного объекта из бакета Minio
}

// minClient реализация интерфейса MinioClient
type minClient struct {
	mc   *minio.Client // Клиент Minio
	conf *config.Config
}

func NewMinioClient(conf *config.Config) MinClient {
	return &minClient{conf: conf} // Возвращает новый экземпляр minioClient с указанным именем бакета
}

// InitMinio подключается к Minio и создает бакет, если не существует
func (m *minClient) InitMinio() error {
	ctx := context.Background()
	// Подключение к Minio
	client, err := minio.New(m.conf.MinioService.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.conf.MinioService.MinioRootUser, m.conf.MinioService.MinioRootPassword, ""),
		Secure: m.conf.MinioService.MinioUseSSL,
	})
	if err != nil {
		return err
	}

	m.mc = client

	// Проверка бакета и его создание, если не существует
	exists, err := m.mc.BucketExists(ctx, m.conf.MinioService.BucketName)
	if err != nil {
		return err
	}
	if !exists {
		err := m.mc.MakeBucket(ctx, m.conf.MinioService.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateOne создает один объект в бакете Minio.
// Метод принимает структуру fileData, которая содержит имя файла и его данные.
// В случае успешной загрузки данных в бакет, метод возвращает nil, иначе возвращает ошибку.
func (m *minClient) CreateOne(file helpers.FileDataType) (string, error) {
	objectID := uuid.New().String()

	reader := bytes.NewReader(file.Data)

	_, err := m.mc.PutObject(context.Background(), m.conf.MinioService.BucketName, objectID, reader, int64(len(file.Data)), minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("ошибка при создании объекта %s: %v", file.FileName, err)
	}

	url, err := m.mc.PresignedGetObject(context.Background(), m.conf.MinioService.BucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании URL для объекта %s: %v", file.FileName, err)
	}

	return url.String(), nil
}

// GetOne получает предварительный url на один объект из бакета Minio по его идентификатору.
func (m *minClient) GetOne(objectID string) (string, error) {
	url, err := m.mc.PresignedGetObject(context.Background(), m.conf.MinioService.BucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении URL для объекта %s: %v", objectID, err)
	}

	return url.String(), nil
}

// DeleteOne удаляет один объект из бакета Minio по его идентификатору.
func (m *minClient) DeleteOne(objectID string) error {
	err := m.mc.RemoveObject(context.Background(), m.conf.MinioService.BucketName, objectID, minio.RemoveObjectOptions{})
	if err != nil {
		return err // Возвращаем ошибку, если не удалось удалить объект.
	}
	return nil
}
