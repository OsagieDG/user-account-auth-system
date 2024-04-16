
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.users (
    id UUID NOT NULL PRIMARY KEY,
    username VARCHAR(20) NOT NULL,
    email VARCHAR(50) NOT NULL,
    encryptedpassword VARCHAR(60) NOT NULL,
    isadmin BOOLEAN NOT NULL DEFAULT false,
    UNIQUE (username, email)
);

