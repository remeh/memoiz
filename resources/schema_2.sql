-- Turn r_category to string in "memo"
-- ----------------------

ALTER TABLE "memo" ADD COLUMN "_r_category" text NOT NULL DEFAULT 'Uncategorized';

UPDATE "memo" SET _r_category = 'Artist' WHERE r_category = 1;
UPDATE "memo" SET _r_category = 'Actor' WHERE r_category = 2;
UPDATE "memo" SET _r_category = 'Book' WHERE r_category = 3;
UPDATE "memo" SET _r_category = 'News' WHERE r_category = 4;
UPDATE "memo" SET _r_category = 'Movie' WHERE r_category = 5;
UPDATE "memo" SET _r_category = 'Music' WHERE r_category = 6;
UPDATE "memo" SET _r_category = 'Person' WHERE r_category = 7;
UPDATE "memo" SET _r_category = 'Place' WHERE r_category = 8;
UPDATE "memo" SET _r_category = 'Serie' WHERE r_category = 9;
UPDATE "memo" SET _r_category = 'Video' WHERE r_category = 10;
UPDATE "memo" SET _r_category = 'VideoGame' WHERE r_category = 11;
UPDATE "memo" SET _r_category = 'Food' WHERE r_category = 12;

ALTER TABLE "memo" DROP COLUMN "r_category";
ALTER TABLE "memo" RENAME COLUMN "_r_category" TO "r_category";

-- Turn category to string in "kg_type"
-- ----------------------

ALTER TABLE "kg_type" ADD COLUMN "_category" text NOT NULL DEFAULT 'Uncategorized';

UPDATE "kg_type" SET _category = 'Artist' WHERE category = 1;
UPDATE "kg_type" SET _category = 'Actor' WHERE category = 2;
UPDATE "kg_type" SET _category = 'Book' WHERE category = 3;
UPDATE "kg_type" SET _category = 'News' WHERE category = 4;
UPDATE "kg_type" SET _category = 'Movie' WHERE category = 5;
UPDATE "kg_type" SET _category = 'Music' WHERE category = 6;
UPDATE "kg_type" SET _category = 'Person' WHERE category = 7;
UPDATE "kg_type" SET _category = 'Place' WHERE category = 8;
UPDATE "kg_type" SET _category = 'Serie' WHERE category = 9;
UPDATE "kg_type" SET _category = 'Video' WHERE category = 10;
UPDATE "kg_type" SET _category = 'VideoGame' WHERE category = 11;
UPDATE "kg_type" SET _category = 'Food' WHERE category = 12;

ALTER TABLE "kg_type" DROP COLUMN "category";
ALTER TABLE "kg_type" RENAME COLUMN "_category" TO "category";

-- domain_result
-- ----------------------

DELETE FROM "domain_result";
ALTER TABLE "domain_result" DROP COLUMN "category";
ALTER TABLE "domain_result" ADD COLUMN "category" text[] NOT NULL;
CREATE INDEX ON "domain_result" ("category");

-- domain
-- ----------------------

DROP TABLE "domain";

CREATE TABLE "domain" (
   "domain" text NOT NULL,
   "category" text NOT NULL DEFAULT 'Uncategorized',
   "weight" int NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "domain" ("domain", "category");

-- kg_result
-- ----------------------

DELETE FROM "kg_result";
ALTER TABLE "kg_result" DROP COLUMN "category";
ALTER TABLE "kg_result" ADD COLUMN "category" text[] NOT NULL;
CREATE INDEX ON "kg_result" ("category");
