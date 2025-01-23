
CREATE TABLE IF NOT EXISTS auth.sessions (
    id SERIAL PRIMARY KEY,
    userid UUID REFERENCES auth.users(id) NOT NULL,
    token TEXT NOT NULL,
    expiresat TIMESTAMPTZ NOT NULL
);