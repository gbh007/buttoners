-- +goose Up

-- Уведомления пользователей
CREATE TABLE notifications(
    id          INT AUTO_INCREMENT   PRIMARY KEY,
    user_id     INT         NOT NULL,
    kind        TEXT        NOT NULL,
    level       TEXT        NOT NULL,
    title       TEXT        NOT NULL,
    body        TEXT,
    `read`      BOOLEAN     NOT NULL DEFAULT 0,
    created     TIMESTAMP   NOT NULL
);
