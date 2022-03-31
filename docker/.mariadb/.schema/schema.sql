CREATE TABLE IF NOT EXISTS videos (
    id          VARCHAR(36) NOT NULL,
    title       VARCHAR(64) NOT NULL,
    v_status    INT NOT NULL, -- v_state, because state and status are SQL keywords
    uploaded_at TIMESTAMP DEFAULT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT pk PRIMARY KEY (id),
    CONSTRAINT unique_title UNIQUE (title)
);

CREATE TABLE IF NOT EXISTS uploads (
    id          VARCHAR(36) NOT NULL,
    v_id        VARCHAR(64) NOT NULL,
    v_status    INT NOT NULL, -- v_state, because state and status are SQL keywords
    uploaded_at TIMESTAMP DEFAULT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT pk PRIMARY KEY (id),
    CONSTRAINT fk_v_id FOREIGN KEY (v_id) REFERENCES videos (id),
);
