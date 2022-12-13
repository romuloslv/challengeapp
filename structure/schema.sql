CREATE TABLE IF NOT EXISTS accounts (
    person_id   VARCHAR(11) PRIMARY KEY UNIQUE NOT NULL,
    first_name  VARCHAR(30) NOT NULL,
    last_name   VARCHAR(20) NOT NULL,
    web_address VARCHAR(50),
    date_birth  VARCHAR(10)
);