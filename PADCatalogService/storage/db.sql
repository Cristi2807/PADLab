CREATE SCHEMA IF NOT EXISTS catalog
    AUTHORIZATION postgres;

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

ALTER TABLE IF EXISTS catalog."Products"
    OWNER to postgres;
