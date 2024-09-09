CREATE TABLE IF NOT EXISTS Users(
	email TEXT PRIMARY KEY,
	password TEXT NOT NULL,
	role TEXT NOT NULL,
	verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS Sessions(
	refresh_token TEXT PRIMARY KEY,
	user_email TEXT NOT NULL,
	deviceID TEXT NOT NULL,
	FOREIGN KEY (user_email) REFERENCES Users (email)
);