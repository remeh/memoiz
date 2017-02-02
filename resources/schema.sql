-- Database init
CREATE USER memoiz WITH UNENCRYPTED PASSWORD 'memoiz';
CREATE DATABASE "memoiz";
GRANT ALL ON DATABASE "memoiz" TO "memoiz";

-- Switch to the memoiz db
\connect "memoiz";
CREATE EXTENSION fuzzystrmatch; -- create the fuzzy match ext
-- Now connect as the memoiz user.
set role "memoiz";

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

    -- emailing
    "unsubscribe_token" text DEFAULT '',

    -- payment
   "stripe_token" text DEFAULT '',

    -- time
    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "last_update" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ON "user" ("uid");
CREATE UNIQUE INDEX ON "user" ("email");

-- Subscription

CREATE TABLE "subscription" (
    "uid" text NOT NULL,
    "owner_uid" text NOT NULL,

    -- which token has been used to pay this subscription
    "stripe_customer_token" text NOT NULL,
    -- which token has been generated while paying this sub
    "stripe_charge_token" text NOT NULL,

    -- which plan has been chosen and when does it end
    "plan_id" text NOT NULL,
    "price" int NOT NULL,
    "end" timestamp with time zone NOT NULL,

    -- stripe response serialized in JSON. Forensic purpose.
    "stripe_response" text NOT NULL,

    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "last_update" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ON "subscription" ("uid");
CREATE INDEX ON "subscription" ("owner_uid");
ALTER TABLE "subscription" ADD CONSTRAINT "subscription_owner_uid" FOREIGN KEY ("owner_uid") REFERENCES "user" ("uid") MATCH FULL;

-- Memo

CREATE TABLE "memo" (
    "uid" text NOT NULL,
    "owner_uid" text NOT NULL,

    "text" text NOT NULL DEFAULT '',
    "position" int NOT NULL DEFAULT 0,
    "state" text NOT NULL DEFAULT 'MemoActive',

    -- rich information
    -- could not be set
    -- r stands for rich
    "r_category" int NOT NULL DEFAULT 0,
    "r_image" text DEFAULT '',
    "r_url" text DEFAULT  '',
    "r_title" text DEFAULT  '',

    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "last_update" timestamp with time zone NOT NULL DEFAULT now(),
    "last_email" timestamp with time zone,
    "archive_time" timestamp with time zone,
    "deletion_time" timestamp with time zone
);

CREATE UNIQUE INDEX ON "memo" ("uid");
ALTER TABLE "memo" ADD CONSTRAINT "memo_owner_uid" FOREIGN KEY ("owner_uid") REFERENCES "user" ("uid") MATCH FULL;

----------------------
-- Domains
-- Use by Bing to analyze the content of a memo
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
    "memo_uid" text NOT NULL,
    "memo_text" text NOT NULL DEFAULT '',
    "category" int[] NOT NULL,
    "domains" text[] NOT NULL,
    "weight" int NOT NULL DEFAULT 0,
    "creation_time" timestamp with time zone DEFAULT now()
);

ALTER TABLE "domain_result" ADD CONSTRAINT "domain_result_memo_uid" FOREIGN KEY ("memo_uid") REFERENCES "memo" ("uid") MATCH FULL;
CREATE UNIQUE INDEX ON "domain_result" ("uid");
CREATE INDEX ON "domain_result" ("memo_text");
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
    "memo_uid" text NOT NULL,
    "memo_text" text NOT NULL DEFAULT '',
    "types" text DEFAULT '',
    "description" text DEFAULT '',
    "category" int[] NOT NULL,
    "creation_time" timestamp with time zone DEFAULT now()
);

ALTER TABLE "kg_result" ADD CONSTRAINT "kg_result_memo_uid" FOREIGN KEY ("memo_uid") REFERENCES "memo" ("uid") MATCH FULL;
CREATE UNIQUE INDEX ON "kg_result" ("uid");
CREATE INDEX ON "kg_result" ("memo_text");
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
ALTER TABLE "emailing_sent" ADD CONSTRAINT "emailing_sent_owner_uid" FOREIGN KEY ("owner_uid") REFERENCES "user" ("uid") MATCH FULL;

CREATE TABLE "emailing_unsubscribe" (
    "owner_uid" text NOT NULL,
    "token" text NOT NULL,
    "reason" text DEFAULT '',
    "creation_time" timestamp with time zone DEFAULT NULL
);

CREATE UNIQUE INDEX ON "emailing_unsubscribe" ("owner_uid");
ALTER TABLE "emailing_unsubscribe" ADD CONSTRAINT "emailing_unsubscribe_owner_uid" FOREIGN KEY ("owner_uid") REFERENCES "user" ("uid") MATCH FULL;

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
insert into "memo" (uid,owner_uid,text,position) values ('abcdabcd-1234-1234-1234-abcdabcdabcd','12341234-1234-1234-1234-123412341234', 'Text of a memo', 0);
insert into "memo" (uid,owner_uid,text,position) values ('1bcdabcd-1234-1234-1234-abcdabcdabcd','12341234-1234-1234-1234-123412341234', 'Text of another memo', 1);
insert into "memo" (uid,owner_uid,text,position) values ('2bcdabcd-1234-1234-1234-abcdabcdabcd','12341234-1234-1234-1234-123412341234', 'Your bones don''t break, mine do. That''s clear. Your cells react to bacteria and viruses differently than mine. You don''t get sick, I do. That''s also clear. But for some reason, you and I react the exact same way to water. We swallow it too fast, we choke. We get some in our lungs, we drown.', 2);
