-- Some example data to populate a database
-- This should only be run from the reload script which drops and creates the database fresh

-- add this to the database with the following:
-- psql <connection string> -f test_tables.sql

INSERT INTO rbac_role (name) VALUES ('DEFAULT');

INSERT INTO member (first_name, last_name, username, password, dob, phone, membership_status_id, rbac_role_id) VALUES 
('Name', 'One', 'email1@address.com', 'aaaaaaaa', '1970-1-1', '316-555-1234', 1, 1),
('Name', 'Two', 'email2@address.com', 'aaaaaaaa', '1970-1-1', '316-555-1234', 1, 1),
('John', 'Doe', 'email3@address.com', 'aaaaaaaa', '1970-1-1', '316-555-1234', 1, 1),
('Jane', 'Doe', 'email4@address.com', 'aaaaaaaa', '1970-1-1', '316-555-1234', 1, 1);

