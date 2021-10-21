-- schema owner
CREATE USER practicum WITH password 'practicum_project';

-- schema user
CREATE USER practicum_ms WITH password 'practicum_project_ms';

CREATE DATABASE mdb WITH OWNER postgres ENCODING 'UTF8';

--CONNECT TO mdb USING postgres;

-- create schema
CREATE SCHEMA practicum_project AUTHORIZATION practicum_project;

GRANT USAGE ON SCHEMA practicum_project TO practicum_project_ms;

ALTER DEFAULT PRIVILEGES FOR USER practicum_project IN SCHEMA practicum_project GRANT SELECT,INSERT,UPDATE,DELETE,TRUNCATE ON TABLES TO practicum_project;
ALTER DEFAULT PRIVILEGES FOR USER practicum_project IN SCHEMA practicum_project GRANT USAGE ON SEQUENCES TO practicum_project;
ALTER DEFAULT PRIVILEGES FOR USER practicum_project IN SCHEMA practicum_project GRANT EXECUTE ON FUNCTIONS TO practicum_project;
