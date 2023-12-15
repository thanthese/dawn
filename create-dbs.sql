CREATE DATABASE recordings;
create role dawn with login password 'test';
GRANT USAGE, SELECT ON SEQUENCE album_id_seq TO dawn;
