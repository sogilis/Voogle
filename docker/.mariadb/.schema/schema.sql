CREATE TABLE IF NOT EXISTS video_state (
    id          INT NOT NULL,
    state_name  ENUM('UPLOADING', 'UPLOADED', 'ENCODING', 'READY', 'FAILURE') NOT NULL,
    CONSTRAINT pk PRIMARY KEY (id)
);

INSERT INTO video_state (id, state_name) VALUES (1, 'UPLOADING');
INSERT INTO video_state (id, state_name) VALUES (2, 'UPLOADED');
INSERT INTO video_state (id, state_name) VALUES (3, 'ENCODING');
INSERT INTO video_state (id, state_name) VALUES (4, 'READY');
INSERT INTO video_state (id, state_name) VALUES (5, 'FAILURE');

CREATE TABLE IF NOT EXISTS videos (
    id          BIGINT UNSIGNED DEFAULT UUID() NOT NULL,
    client_id   BIGINT UNSIGNED DEFAULT UUID() NOT NULL,
    title       VARCHAR(64) NOT NULL,
    v_state     INT NOT NULL,   -- v_state, because state and status are SQL keywords
    last_update DATE NOT NULL DEFAULT NOW(),
    CONSTRAINT pk PRIMARY KEY (id),
    CONSTRAINT fk_v_state FOREIGN KEY (v_state) REFERENCES video_state (id)
);

INSERT INTO videos (title, v_state) VALUES ('Video1', 1);
INSERT INTO videos (title, v_state) VALUES ('Video2', 5);