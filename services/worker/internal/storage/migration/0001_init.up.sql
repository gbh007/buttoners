CREATE TABLE task_results(
    user_id         Int64         NOT NULL,
    chance          Int64         NOT NULL,
    duration        Int64         NOT NULL,
    result          Int64         NOT NULL,
    result_text     String        NOT NULL,
    error_text      String        NOT NULL,
    start_time      DateTime64(9, 'UTC')   NOT NULL,
    end_time        DateTime64(9, 'UTC')   NOT NULL
)
ENGINE = MergeTree()
ORDER BY start_time;