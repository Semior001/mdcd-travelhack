DROP TABLE IF EXISTS "images" CASCADE;
DROP SEQUENCE IF EXISTS images_id_seq;
DROP TABLE IF EXISTS "users" CASCADE;
DROP SEQUENCE IF EXISTS users_id_seq;

CREATE SEQUENCE images_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

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
    "bar_code" text NOT NULL,
    "mime" text NOT NULL,
    "img_type" text NOT NULL,
    "local_filename" text NOT NULL,
    "user_id" bigint DEFAULT 1,
    "created_at" timestamptz DEFAULT NOW(),
    "updated_at" timestamptz DEFAULT NOW(),
    CONSTRAINT "images_bar_code_key" UNIQUE ("bar_code"),
    CONSTRAINT "images_pkey" PRIMARY KEY ("id")

    -- Make sure there is no foreign key in test DB
    -- CONSTRAINT "images_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) NOT DEFERRABLE
) WITH (oids = false);