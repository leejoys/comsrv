DROP TABLE IF EXISTS comments;

CREATE TABLE IF NOT EXISTS comments(
id SERIAL PRIMARY KEY,
author TEXT NOT NULL,
content TEXT NOT NULL,
pubtime BIGINT NOT NULL,
parentpost BIGINT,
parentcomment BIGINT
);

INSERT INTO comments (id, author, content, pubtime, parentpost, parentcomment) 
VALUES (0, 'Админ', 'Содержание комментария', 0, 0, 0);
