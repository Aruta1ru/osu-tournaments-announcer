CREATE DATABASE posts;
\c posts;

CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY,
    username VARCHAR NOT NULL,
    avatar VARCHAR NOT NULL,
    country_code VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS forumposts (
    id INT PRIMARY KEY,
    title VARCHAR NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    edited_at TIMESTAMP WITH TIME ZONE,
    picture_preview VARCHAR NOT NULL,
    is_valid BOOLEAN,
    is_notified BOOLEAN DEFAULT false,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS forumpost_links (
    id SERIAL PRIMARY KEY,
    forumpost_id INT NOT NULL,
    name VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
    FOREIGN KEY (forumpost_id) REFERENCES forumposts (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notify_servers (
    id SERIAL PRIMARY KEY,
    server_id BIGINT NOT NULL UNIQUE,
    channel_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);