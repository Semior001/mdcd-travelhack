CREATE DATABASE test;
\c test;

DROP TABLE IF EXISTS "images";
DROP SEQUENCE IF EXISTS images_id_seq;
CREATE SEQUENCE images_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."images" (
    "id" integer DEFAULT nextval('images_id_seq') NOT NULL,
    "bar_code" text,
    "img_format" text,
    "img_type_id" integer,
    "local_filename" text,
    "user_id" bigint,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    CONSTRAINT "images_bar_code_key" UNIQUE ("bar_code"),
    CONSTRAINT "images_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "images_img_type_id_fkey" FOREIGN KEY (img_type_id) REFERENCES img_types(id) NOT DEFERRABLE,
    CONSTRAINT "images_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) NOT DEFERRABLE
) WITH (oids = false);


DROP TABLE IF EXISTS "img_types";
CREATE TABLE "public"."img_types" (
    "id" integer NOT NULL,
    "code" text NOT NULL,
    CONSTRAINT "img_types_code_key" UNIQUE ("code"),
    CONSTRAINT "img_types_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

DROP TABLE IF EXISTS "users";
DROP SEQUENCE IF EXISTS users_id_seq;

CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."users" (
    "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "email" text NOT NULL,
    "password" text,
    "privileges" jsonb,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    CONSTRAINT "users_email_key" UNIQUE ("email"),
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

-- inserting initial values

INSERT INTO "img_types" ("id", "code") VALUES
(0,	'source'),
(1,	'background');

INSERT INTO "users" ("id", "email", "password", "privileges", "created_at", "updated_at") VALUES
(8,	'blah@blah.blah',	'blah',	NULL,	'2020-03-09 09:40:56.001599+00',	'2020-03-09 09:40:56.001599+00');