CREATE TABLE user_logs(
    request_id      String        NOT NULL,
    addr            String        NOT NULL,
    session_token   String        NULL,
    user_id         Int64         NULL,
    action          String        NOT NULL,
    chance          Int64         NULL,
    duration        Int64         NULL,
    request_time    DateTime64(9, 'UTC')   NOT NULL
)
ENGINE = MergeTree()
ORDER BY request_time;