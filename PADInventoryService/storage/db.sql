 -- Database: inventory_pad_lab

DROP DATABASE IF EXISTS inventory_pad_lab;

CREATE DATABASE inventory_pad_lab
    WITH
    OWNER = cristi
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

\c inventory_pad_lab;

CREATE SCHEMA IF NOT EXISTS inventory;

--------------------------------

CREATE TABLE inventory."Transactions"
(
    "Id"            uuid                     NOT NULL DEFAULT gen_random_uuid(),
    "ShoesId"       uuid                     NOT NULL,
    "CreationDate"  timestamp with time zone NOT NULL,
    "Quantity"      bigint                   NOT NULL,
    "OperationType" smallint                 NOT NULL,
    PRIMARY KEY ("Id")
) TABLESPACE pg_default;

