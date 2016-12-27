-- Rottentomatoes
INSERT INTO "domain" VALUES ('rottentomatoes', 5, 100);
INSERT INTO "domain" VALUES ('rottentomatoes', 9, 100);

-- IMDb
INSERT INTO "domain" VALUES ('imdb', 5, 100);
INSERT INTO "domain" VALUES ('imdb', 9, 100);

-- Allocine
INSERT INTO "domain" VALUES ('allocine', 5, 100);
INSERT INTO "domain" VALUES ('allocine', 9, 100);

-- Wikipedia is unknown because it could speak about everything
INSERT INTO "domain" VALUES ('wikipedia', 0, 100);

-- Youtube could be either a video or a music
INSERT INTO "domain" VALUES ('youtube', 10, 100); -- video
INSERT INTO "domain" VALUES ('youtube', 6, 75); -- music
