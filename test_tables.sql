-- Some example data to populate a database

-- add this to the database with the following:
-- psql <connection string> -f test_tables.sql

DROP TABLE IF EXISTS users;

CREATE TABLE users (
id SERIAL,
name varchar(255) NOT NULL,
email varchar(255) NOT NULL,
CONSTRAINT users_pkey PRIMARY KEY (id)
);

INSERT INTO users (name, email) VALUES 
('Name One', 'email@address.com'),
('Name Two', 'something@example.com'),
('John Doe', 'nomail@junk.com'),
('Jane Doe', 'junk@example.com');