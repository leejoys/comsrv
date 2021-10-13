DROP TABLE IF EXISTS comments;

CREATE TABLE IF NOT EXISTS comments(
id SERIAL PRIMARY KEY,
author TEXT NOT NULL,
content TEXT NOT NULL,
pubtime BIGINT NOT NULL,
parentpost BIGINT,
parentcomment BIGINT
);
