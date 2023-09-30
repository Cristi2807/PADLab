CREATE SCHEMA IF NOT EXISTS inventory
    AUTHORIZATION postgres;

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

ALTER TABLE IF EXISTS inventory."Transactions"
    OWNER to postgres;
