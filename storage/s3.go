package storage

import (
	"crypto/tls"
	"net/http"

	"bytes"
	"fmt"
	"log"
	"mime/multipart"

	"apk-server/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Client *s3.S3
var s3Bucket string

// InitS3 инициализирует клиента S3, используя параметры из конфигурации.
func InitS3(cfg *config.Config) error {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(cfg.AWSRegion),
		Endpoint: aws.String("https://storage.yandexcloud.net"), // Укажите endpoint Yandex Object Storage
		Credentials: credentials.NewStaticCredentials(
			cfg.AWSAccessKey,
			cfg.AWSSecretKey,
			"",
		),
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // Отключает проверку сертификата
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ошибка при создании сессии S3: %v", err)
	}
	log.Println(*sess.Config.Region)
	log.Println(cfg.AWSAccessKey)

	s3Client = s3.New(sess)
	s3Bucket = cfg.S3Bucket
	return nil
}

// UploadToS3 загружает файл в S3 и возвращает публичный URL.
func UploadToS3(file multipart.File, fileName string) (string, error) {
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении файла: %v", err)
	}

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(buffer.Bytes()),
		ACL:    aws.String("public-read"), // Убедитесь, что это необходимо
	})
	if err != nil {
		return "", fmt.Errorf("ошибка при загрузке файла в S3: %v", err)
	}

	fileURL := "https://storage.yandexcloud.net/" + s3Bucket + "/" + fileName
	log.Printf("Файл успешно загружен: %s", fileURL)
	return fileURL, nil
}
