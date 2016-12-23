-- Database init
CREATE USER scratche WITH UNENCRYPTED PASSWORD 'scratche';
CREATE DATABASE "scratche";
GRANT ALL ON DATABASE "scratche" TO "scratche";

-- Switch to the scratche db as the scratche user.
\connect "scratche";
set role "scratche";

-- User

CREATE TABLE "user" (
    "uid" text NOT NULL,
    "email" text NOT NULL DEFAULT '',
    "firstname" text NOT NULL DEFAULT '',
    "lastname" text NOT NULL DEFAULT '',
    "password" text NOT NULL DEFAULT '',
    "gender" text NOT NULL DEFAULT 'Undefined',
    "phone" text NOT NULL DEFAULT '',
    "address" text NOT NULL DEFAULT '',
    "timezone" text NOT NULL DEFAULT '',

    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "last_update" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ON "user" ("uid");
CREATE INDEX ON "user" ("email");

-- Card

CREATE TABLE "card" (
    "uid" text NOT NULL,
    "user_uid" text NOT NULL,

    "text" text NOT NULL DEFAULT '',
    "position" int NOT NULL DEFAULT 0,
    "state" text NOT NULL DEFAULT 'CardActive',

    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "last_update" timestamp with time zone NOT NULL DEFAULT now(),
    "deletion_time" timestamp with time zone
);

CREATE UNIQUE INDEX ON "card" ("uid");
ALTER TABLE "card" ADD CONSTRAINT "card_owner_uid" FOREIGN KEY ("user_uid") REFERENCES "user" ("uid") MATCH FULL;

-- DB Schema

CREATE TABLE "db_schema" (
    "version" int NOT NULL DEFAULT 0,
    "update_time" timestamp with time zone NOT NULL DEFAULT now()
);

--

INSERT INTO "db_schema" VALUES (
    1,
    now()
);

----------------------

insert into "user" (uid) values ('12341234-1234-1234-1234-123412341234');
insert into "card" (uid,user_uid,text) values ('abcdabcd-abcd-abcd-abcd-abcdabcdabcd','12341234-1234-1234-1234-123412341234', 'Text of a card');
insert into "card" (uid,user_uid,text) values ('1bcdabcd-abcd-abcd-abcd-abcdabcdabcd','12341234-1234-1234-1234-123412341234', 'Text of another card');
insert into "card" (uid,user_uid,text) values ('2bcdabcd-abcd-abcd-abcd-abcdabcdabcd','12341234-1234-1234-1234-123412341234', 'Your bones don''t break, mine do. That''s clear. Your cells react to bacteria and viruses differently than mine. You don''t get sick, I do. That''s also clear. But for some reason, you and I react the exact same way to water. We swallow it too fast, we choke. We get some in our lungs, we drown.');
