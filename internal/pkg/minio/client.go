package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/pkg/config"
	"github.com/TeaStealers-backend-sem4/internal/pkg/logger"
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
	mc     *minio.Client // Клиент Minio
	conf   *config.Config
	logger logger.Logger
}

func NewMinioClient(conf *config.Config, logr logger.Logger) MinClient {
	return &minClient{conf: conf, logger: logr} // Возвращает новый экземпляр minioClient с указанным именем бакета
}

func (m *minClient) SetBucketPolicy() error {
	policy := `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "*"
                ]
            },
            "Action": [
                "s3:GetObject",
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::defaultbucket",
                "arn:aws:s3:::defaultbucket/*"
            ]
        }
    ]
}`

	err := m.mc.SetBucketPolicy(context.Background(), m.conf.MinioService.BucketName, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %v", err)
	}
	return nil
}

// InitMinio подключается к Minio и создает бакет, если не существует
func (m *minClient) InitMinio() error {
	ctx := context.Background()
	// Подключение к Minio
	client, err := minio.New(m.conf.MinioService.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.conf.MinioService.MinioRootUser, m.conf.MinioService.MinioRootPassword, ""),
		Secure: m.conf.MinioService.MinioUseSSL,
	})
	m.logger.LogDebug("connected to minio")
	if err != nil {
		m.logger.LogDebug(err.Error())
		return err
	}

	m.mc = client

	// Проверка бакета и его создание, если не существует
	exists, err := m.mc.BucketExists(ctx, m.conf.MinioService.BucketName)
	if err != nil {
		m.logger.LogDebug(err.Error())
		return err
	}
	if !exists {

		err := m.mc.MakeBucket(ctx, m.conf.MinioService.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			m.logger.LogDebug(err.Error())
			return err
		}
		m.logger.LogDebug("created bucket")
	}
	m.logger.LogDebug("minio client init success")
	err = m.SetBucketPolicy()
	return err
}

// CreateOne создает один объект в бакете Minio.
// Метод принимает структуру fileData, которая содержит имя файла и его данные.
// В случае успешной загрузки данных в бакет, метод возвращает nil, иначе возвращает ошибку.
func (m *minClient) CreateOne(file helpers.FileDataType) (string, error) {
	objectID := uuid.New().String()

	reader := bytes.NewReader(file.Data)

	_, err := m.mc.PutObject(context.Background(), m.conf.MinioService.BucketName, objectID, reader, int64(len(file.Data)), minio.PutObjectOptions{})
	if err != nil {
		m.logger.LogDebug(err.Error())
		return "", fmt.Errorf("failed to create object %s: %v", file.FileName, err)
	}

	url, err := m.mc.PresignedGetObject(context.Background(), m.conf.MinioService.BucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		m.logger.LogDebug(err.Error())
		return "", fmt.Errorf("fail create url for object %s: %v", file.FileName, err)
	}
	m.logger.LogDebug("created file with url: " + url.String())
	return objectID, nil
}

// GetOne получает предварительный url на один объект из бакета Minio по его идентификатору.
func (m *minClient) GetOne(objectID string) (string, error) {
	url, err := m.mc.PresignedGetObject(context.Background(), m.conf.MinioService.BucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		m.logger.LogDebug(err.Error())
		return "", fmt.Errorf("fail to get url for object %s: %v", objectID, err)
	}
	m.logger.LogDebug("got url for object")

	return url.String(), nil
}

// DeleteOne удаляет один объект из бакета Minio по его идентификатору.
func (m *minClient) DeleteOne(objectID string) error {
	err := m.mc.RemoveObject(context.Background(), m.conf.MinioService.BucketName, objectID, minio.RemoveObjectOptions{})
	if err != nil {
		m.logger.LogDebug(err.Error())
		return err
	}
	m.logger.LogDebug("deleted object")
	return nil
}
