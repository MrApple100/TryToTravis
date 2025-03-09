package models

import "time"

type Application struct {
    ID          int    `db:"id"`
    AppID       string `db:"app_id"`
    Name        string `db:"name"`
    PackageName string `db:"package_name"`
}

type ApkVersion struct {
    ID          int       `db:"id"`
    AppID       string    `db:"app_id"`
    VersionName string    `db:"version_name"`
    VersionCode int       `db:"version_code"`
    FileURL     string    `db:"file_url"`
    UploadedAt  time.Time `db:"uploaded_at"`
}
