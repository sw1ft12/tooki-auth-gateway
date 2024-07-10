
CREATE TYPE GENDER AS ENUM('Male', 'Female', 'Unknown');
CREATE TYPE ROLE AS ENUM('USER', 'ADMIN', 'SUPER_ADMIN');

CREATE TABLE Users
(
    id UUID PRIMARY KEY DEFAULT  gen_random_uuid(),
    email VARCHAR(50) NOT NULL UNIQUE,
    login VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    name TEXT NOT NULL,
    age INT,
    gender GENDER NOT NULL,
    role ROLE DEFAULT 'USER',
    avatar VARCHAR(100),
    verified BOOL DEFAULT FALSE,
    banned BOOL DEFAULT FALSE,
    created_at DATE DEFAULT CURRENT_DATE,
    updated_at DATE DEFAULT CURRENT_DATE
);

CREATE TABLE Tokens (
    token UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expires_in DATE DEFAULT CURRENT_DATE + 30,
    user_id UUID UNIQUE REFERENCES users(id),
    user_agent TEXT
)
