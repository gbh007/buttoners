-- +goose Up

-- Данные пользователя
CREATE TABLE users(
    id          INT AUTO_INCREMENT   PRIMARY KEY,
    login       TEXT        NOT NULL UNIQUE,
    password    TEXT        NOT NULL,
    salt        TEXT        NOT NULL,
    created     TIMESTAMP   NOT NULL,
    updated     TIMESTAMP  
);

-- Сессии пользователей
CREATE TABLE sessions(
    token       VARCHAR(64)             PRIMARY KEY,
    user_id     INT         NOT NULL    REFERENCES users (id),
    created     TIMESTAMP   NOT NULL,
    used        TIMESTAMP,
    updated     TIMESTAMP  
);
