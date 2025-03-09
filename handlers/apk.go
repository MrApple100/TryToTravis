package handlers

import (
	"apk-server/database"
	"apk-server/models"
	"apk-server/storage"

	"fmt"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UploadApk(c *gin.Context) {
	appID := c.PostForm("app_id")
	versionName := c.PostForm("version_name")
	versionCode, _ := strconv.Atoi(c.PostForm("version_code"))

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный файл"})
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла"})
		return
	}
	defer fileReader.Close()

	fileURL, err := storage.UploadToS3(fileReader, appID+"/"+file.Filename)
	if err != nil {
		log.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки файла в S3"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Errorf("ошибка при начале транзакции: %v", err)
	}
	defer tx.Rollback()

	// Добавление записи в таблицу applications
	_, err = tx.Exec(`
        INSERT INTO applications (app_id, name, package_name)
        VALUES ($1, $2, $3)
        ON CONFLICT (app_id) DO NOTHING;
    `, appID, versionName, file.Filename)
	if err != nil {
		fmt.Errorf("ошибка при добавлении записи в applications: %v", err)
	}

	// Добавление записи в таблицу apk_versions
	_, err = tx.Exec(`
        INSERT INTO apk_versions (app_id, version_name, version_code, file_url)
        VALUES ($1, $2, $3, $4);
    `, appID, versionName, versionCode, fileURL)
	if err != nil {
		fmt.Errorf("ошибка при добавлении записи в apk_versions: %v", err)
	}

	// Завершение транзакции
	if err := tx.Commit(); err != nil {
		fmt.Errorf("ошибка при завершении транзакции: %v", err)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения данных в БД"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "APK успешно загружен", "file_url": fileURL})
}

func GetLatestVersion(c *gin.Context) {
	appID := c.Query("app_id")
	var version models.ApkVersion
	err := database.DB.QueryRow(
		"SELECT app_id, version_name, version_code, file_url, uploaded_at FROM apk_versions WHERE app_id = $1 ORDER BY version_code DESC LIMIT 1",
		appID,
	).Scan(&version.AppID, &version.VersionName, &version.VersionCode, &version.FileURL, &version.UploadedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Версия не найдена"})
		return
	}

	c.JSON(http.StatusOK, version)
}

func DownloadApk(c *gin.Context) {
	appID := c.Query("app_id")
	var fileURL string
	err := database.DB.QueryRow(
		"SELECT file_url FROM apk_versions WHERE app_id = $1 ORDER BY version_code DESC LIMIT 1",
		appID,
	).Scan(&fileURL)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Версия не найдена"})
		return
	}

	c.Redirect(http.StatusFound, fileURL)
}
