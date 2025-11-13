CREATE TABLE IF NOT EXISTS users
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    email       TEXT NOT NULL UNIQUE,
    pass_hash   BLOB NOT NULL,
    is_admin    INTEGER
);