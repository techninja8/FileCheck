-- schema.sql
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS files (
	id TEXT PRIMARY KEY,
	owner TEXT NOT NULL, -- We'll get the username from the token and append it here
	filename TEXT NOT NULL,
	hash TEXT NOT NULL,
	uploaded_at TIMESTAMP NOT NULL,
	location TEXT NOT NULL
);