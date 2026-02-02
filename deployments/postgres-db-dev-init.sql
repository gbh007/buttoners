CREATE ROLE legacy_user WITH LOGIN PASSWORD 'legacy_pwd';
CREATE DATABASE legacy_db OWNER legacy_user;