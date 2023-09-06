CREATE TABLE accounts (
    user_id serial PRIMARY KEY,
    username VARCHAR ( 50 ) UNIQUE NOT NULL,
    password VARCHAR ( 50 ) NOT NULL,
    email VARCHAR ( 255 ) UNIQUE NOT NULL,
    created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
);
INSERT INTO accounts VALUES (1, 'john', 'password', 'john@gmail.com', '2023-07-12 00:00:00', '2023-07-12 00:01:00');
