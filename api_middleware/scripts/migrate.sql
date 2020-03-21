DROP TABLE IF EXISTS "images" CASCADE;
DROP SEQUENCE IF EXISTS images_id_seq;
DROP TABLE IF EXISTS "img_types" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP SEQUENCE IF EXISTS users_id_seq;

CREATE SEQUENCE images_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;
CREATE TABLE "img_types" (
                                      "id" integer NOT NULL,
                                      "code" text NOT NULL,
                                      CONSTRAINT "img_types_code_key" UNIQUE ("code"),
                                      CONSTRAINT "img_types_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;
CREATE TABLE "users" (
                         "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
                         "email" text NOT NULL,
                         "password" text NOT NULL,
                         "privileges" jsonb DEFAULT '{}'::jsonb,
                         "created_at" timestamptz DEFAULT NOW(),
                         "updated_at" timestamptz DEFAULT NOW(),
                         CONSTRAINT "users_email_key" UNIQUE ("email"),
                         CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE TABLE "images" (
    "id" integer DEFAULT nextval('images_id_seq') NOT NULL,
    "bar_code" text,
    "img_format" text,
    "img_type_id" integer,
    "local_filename" text,
    "user_id" bigint,
    "created_at" timestamptz DEFAULT NOW(),
    "updated_at" timestamptz DEFAULT NOW(),
    CONSTRAINT "images_bar_code_key" UNIQUE ("bar_code"),
    CONSTRAINT "images_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "images_img_type_id_fkey" FOREIGN KEY (img_type_id) REFERENCES img_types(id) NOT DEFERRABLE,
    CONSTRAINT "images_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) NOT DEFERRABLE
) WITH (oids = false);

-- inserting initial values

INSERT INTO "img_types" ("id", "code") VALUES
(0,	'source'),
(1,	'background');