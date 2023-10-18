-- Database: catalog_pad_lab

DROP DATABASE IF EXISTS catalog_pad_lab;

CREATE DATABASE catalog_pad_lab
    WITH
    OWNER = cristi
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

\c catalog_pad_lab;
	
CREATE SCHEMA IF NOT EXISTS catalog;
--------------------------------

CREATE TABLE catalog."Products"
(
    "Id"       uuid                   NOT NULL DEFAULT gen_random_uuid(),
    "Color"    character varying(255) NOT NULL,
    "Size"     character varying(255) NOT NULL,
    "Price"    character varying(255) NOT NULL,
    "Brand"    character varying(255) NOT NULL,
    "Category" character varying(255) NOT NULL,
    "Model"    character varying(255) NOT NULL,
    PRIMARY KEY ("Id")
) TABLESPACE pg_default;

