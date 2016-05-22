CREATE DATABASE IF NOT EXISTS gizix;
CREATE TABLE IF NOT EXISTS gizix.users (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  password VARCHAR(256) NOT NULL,
  icon_path VARCHAR(256),
  admin BOOLEAN,
  created_at DATETIME,
  logined_at DATETIME,
  PRIMARY KEY (id),
  CONSTRAINT unique_name UNIQUE (name)
) ENGINE = InnoDB;
CREATE TABLE IF NOT EXISTS gizix.rooms (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  created_at DATETIME,
  called_at DATETIME,
  PRIMARY KEY (id),
  CONSTRAINT unique_name UNIQUE (name)
) ENGINE = InnoDB;
CREATE TABLE IF NOT EXISTS gizix.user_room (
  id INT AUTO_INCREMENT,
  user_id INT NOT NULL,
  room_id INT NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fkey_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
  CONSTRAINT fkey_room
    FOREIGN KEY (room_id)
    REFERENCES rooms(id)
    ON DELETE CASCADE
) ENGINE = InnoDB;
CREATE TABLE IF NOT EXISTS gizix.domain (
  id INT AUTO_INCREMENT,
  name VARCHAR(256),
  PRIMARY KEY (id)
) ENGINE = InnoDB;
CREATE TABLE IF NOT EXISTS gizix.skyway (
  id INT AUTO_INCREMENT,
  api_key VARCHAR(64),
  PRIMARY KEY (id)
) ENGINE = InnoDB;
INSERT IGNORE INTO gizix.users (name, password, admin) VALUES ('Gizix', '$2a$10$Zg9nPS07epk/CT8PlyHtZei4FOGhtyKyl49Xvpmlrh.BHZKgdyYPS', 1);
INSERT IGNORE INTO gizix.domain (id, name) VALUES (1, 'example.com');
INSERT IGNORE INTO gizix.skyway (id, api_key) VALUES (1, 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx');
