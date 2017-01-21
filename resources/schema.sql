-- Database init
CREATE USER scratche WITH UNENCRYPTED PASSWORD 'scratche';
CREATE DATABASE "scratche";
GRANT ALL ON DATABASE "scratche" TO "scratche";

-- Switch to the scratche db
\connect "scratche";
CREATE EXTENSION fuzzystrmatch; -- create the fuzzy match ext
-- Now connect as the scratche user.
set role "scratche";

-- User

CREATE TABLE "user" (
    "uid" text NOT NULL,
    "email" text NOT NULL DEFAULT '',
    "firstname" text NOT NULL DEFAULT '',
    "lastname" text NOT NULL DEFAULT '',
    "hash" text NOT NULL DEFAULT '',
    "gender" text NOT NULL DEFAULT 'Undefined',
    "phone" text NOT NULL DEFAULT '',
    "address" text NOT NULL DEFAULT '',
    "timezone" text NOT NULL DEFAULT '',

    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "last_update" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ON "user" ("uid");
CREATE UNIQUE INDEX ON "user" ("email");

-- Card

CREATE TABLE "card" (
    "uid" text NOT NULL,
    "owner_uid" text NOT NULL,

    "text" text NOT NULL DEFAULT '',
    "position" int NOT NULL DEFAULT 0,
    "state" text NOT NULL DEFAULT 'CardActive',

    -- rich information
    -- could not be set
    -- r stands for rich
    "r_category" int NOT NULL DEFAULT 0,
    "r_image" text DEFAULT '',
    "r_url" text DEFAULT  '',
    "r_title" text DEFAULT  '',

    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "last_update" timestamp with time zone NOT NULL DEFAULT now(),
    "archive_time" timestamp with time zone,
    "deletion_time" timestamp with time zone
);

CREATE UNIQUE INDEX ON "card" ("uid");
ALTER TABLE "card" ADD CONSTRAINT "card_owner_uid" FOREIGN KEY ("owner_uid") REFERENCES "user" ("uid") MATCH FULL;

----------------------
-- Domains
-- Use by Bing to analyze the content of a card
-- by assigning Category to domains.
----------------------

CREATE TABLE "domain" (
   "domain" text NOT NULL,
   "category" int NOT NULL DEFAULT 0,
   "weight" int NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "domain" ("domain", "category");

-- DomainResult
-- Analyzes results computed using domains.

CREATE TABLE "domain_result" (
    "uid" text NOT NULL,
    "card_uid" text NOT NULL,
    "card_text" text NOT NULL DEFAULT '',
    "category" int[] NOT NULL,
    "domains" text[] NOT NULL,
    "weight" int NOT NULL DEFAULT 0,
    "creation_time" timestamp with time zone DEFAULT now()
);

ALTER TABLE "domain_result" ADD CONSTRAINT "domain_result_card_uid" FOREIGN KEY ("card_uid") REFERENCES "card" ("uid") MATCH FULL;
CREATE UNIQUE INDEX ON "domain_result" ("uid");
CREATE INDEX ON "domain_result" ("card_text");
CREATE INDEX ON "domain_result" ("category");

----------------------
-- KGType
-- Matching
----------------------

CREATE TABLE "kg_type" (
    "type" text NOT NULL,
    "category" int NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "kg_type" ("type");

CREATE TABLE "kg_result" (
    "uid" text NOT NULL,
    "card_uid" text NOT NULL,
    "card_text" text NOT NULL DEFAULT '',
    "types" text DEFAULT '',
    "description" text DEFAULT '',
    "category" int[] NOT NULL,
    "creation_time" timestamp with time zone DEFAULT now()
);

ALTER TABLE "kg_result" ADD CONSTRAINT "kg_result_card_uid" FOREIGN KEY ("card_uid") REFERENCES "card" ("uid") MATCH FULL;
CREATE UNIQUE INDEX ON "kg_result" ("uid");
CREATE INDEX ON "kg_result" ("card_text");
CREATE INDEX ON "kg_result" ("category");

----------------------
-- Emailing
----------------------

CREATE TABLE "emailing_sent" (
    "uid" text NOT NULL,
    "owner_uid" text NOT NULL,
    "type" text NOT NULL,
    "creation_time" timestamp with time zone DEFAULT now()
);

CREATE UNIQUE INDEX ON "emailing_sent" ("uid");
CREATE INDEX ON "emailing_sent" ("owner_uid");
CREATE INDEX ON "emailing_sent" ("owner_uid","creation_time");

CREATE TABLE "emailing_unsubscribe" (
    "owner_uid" text NOT NULL,
    "creation_time" timestamp with time zone DEFAULT NULL
);

CREATE UNIQUE INDEX ON "emailing_unsubscribe" ("owner_uid");

----------------------
-- DB Schema
----------------------

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
insert into "card" (uid,owner_uid,text,position) values ('abcdabcd-1234-1234-1234-abcdabcdabcd','12341234-1234-1234-1234-123412341234', 'Text of a card', 0);
insert into "card" (uid,owner_uid,text,position) values ('1bcdabcd-1234-1234-1234-abcdabcdabcd','12341234-1234-1234-1234-123412341234', 'Text of another card', 1);
insert into "card" (uid,owner_uid,text,position) values ('2bcdabcd-1234-1234-1234-abcdabcdabcd','12341234-1234-1234-1234-123412341234', 'Your bones don''t break, mine do. That''s clear. Your cells react to bacteria and viruses differently than mine. You don''t get sick, I do. That''s also clear. But for some reason, you and I react the exact same way to water. We swallow it too fast, we choke. We get some in our lungs, we drown.', 2);
