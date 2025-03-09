-- init.sql

-- Создание таблицы для приложений
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    app_id VARCHAR(255) UNIQUE,
    name VARCHAR(255),
    package_name VARCHAR(255)
);

-- Создание таблицы для версий APK
CREATE TABLE IF NOT EXISTS apk_versions (
    id SERIAL PRIMARY KEY,
    app_id VARCHAR(255) REFERENCES applications(app_id),
    version_name VARCHAR(50),
    version_code INT,
    file_url TEXT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);