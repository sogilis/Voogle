CREATE TABLE IF NOT EXISTS videos (
    id              VARCHAR(36) NOT NULL,
    title           VARCHAR(64) NOT NULL,
    video_status    INT NOT NULL,
    uploaded_at     DATETIME,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT pk PRIMARY KEY (id),
    CONSTRAINT unique_title UNIQUE (title)
);

CREATE TABLE IF NOT EXISTS uploads (
    id              VARCHAR(36) NOT NULL,
    video_id        VARCHAR(36) NOT NULL,
    upload_status   INT NOT NULL,
    uploaded_at     DATETIME,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT pk PRIMARY KEY (id),
    CONSTRAINT fk_v_id FOREIGN KEY (video_id) REFERENCES videos (id)
);
