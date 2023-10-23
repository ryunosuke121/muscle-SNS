-- +migrate Up
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS trainings;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS menus;

CREATE TABLE user_groups (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(32) NOT NULL,
    image_url VARCHAR(255),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE menus (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE users (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    image_url VARCHAR(255),
    user_group_id INTEGER,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_group_id) REFERENCES user_groups(id)
);

CREATE TABLE trainings (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(64) NOT NULL,
    menu_id INTEGER NOT NULL,
    times INTEGER NOT NULL,
    weight INTEGER NOT NULL,
    sets INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (menu_id) REFERENCES menus(id)
);

CREATE TABLE posts (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(64) NOT NULL,
    training_id INTEGER NOT NULL,
    comment VARCHAR(255) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (training_id) REFERENCES trainings(id)
);

-- +migrate Down
DROP TABLE posts;
DROP TABLE trainings;
DROP TABLE users;
DROP TABLE user_groups;
DROP TABLE menus;

